package colours

// types.go implements the interfaces defined by interfaces.go.

import (
	"context"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/fiffu/arisa3/app/instrumentation"
	"github.com/fiffu/arisa3/lib"
)

// Reason encodes reasons for colour changes (or lack thereof).
type Reason string

func (r Reason) String() string { return string(r) }

const (
	Mutate Reason = "mutate"
	Reroll Reason = "reroll"
	Freeze Reason = "freeze"
)

/* Sentinel values */

// Never indicates that user has state within the domain, but no records for the given Reason.
var Never = time.Time{} // the zero value

// ColourState models a participant's state in the Colour Roles domain.
type ColourState struct {
	UserID     string
	LastFrozen time.Time
	LastMutate time.Time
	LastReroll time.Time
}

// History models a participant's history of colours.
type History struct {
	records    []*ColoursLogRecord
	start, end time.Time
}

func isPartOfHistory(r Reason) bool {
	switch r {
	case Reroll, Mutate:
		return true
	default:
		return false
	}
}

// session implements IDomainSession
type session struct {
	sess         *discordgo.Session
	cacheMembers lib.ICache[IDomainMember, string]
	cacheRoles   lib.ICache[IDomainRole, string]
}

func NewDomainSession(sess *discordgo.Session) IDomainSession {
	return &session{
		sess,
		lib.NewCache[IDomainMember, string](1 * time.Hour),
		lib.NewCache[IDomainRole, string](1 * time.Hour),
	}
}

func (s *session) GuildMember(ctx context.Context, guildID, userID string) (IDomainMember, error) {
	ctx, span := instrumentation.SpanInContext(ctx, instrumentation.Vendor(s.sess.GuildMember))
	defer span.End()

	// Cache lookup
	if cached, ok := s.cacheMembers.Peek(userID); ok {
		return cached, nil
	}

	// Query API for guild member
	mem, err := s.sess.GuildMember(guildID, userID, discordgo.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	// Query for guild roles.
	allRoles, err := s.GuildRoles(ctx, mem.GuildID)
	if err != nil {
		return nil, err
	}

	// Merge on member roles.
	roles := make([]IDomainRole, 0)
	for _, roleID := range mem.Roles {
		for _, guildRole := range allRoles {
			if roleID == guildRole.ID() {
				roles = append(roles, guildRole)
			}
		}
	}

	d := NewDomainMember(mem, roles)
	s.cacheMembers.Put(d)
	return d, nil
}

func (s *session) GuildRole(ctx context.Context, guildID, roleID string) (IDomainRole, error) {
	if cached, ok := s.cacheRoles.Peek(roleID); ok {
		return cached, nil
	}
	roles, err := s.GuildRoles(ctx, guildID)
	if err != nil {
		return nil, err
	}
	for _, role := range roles {
		if role.ID() == roleID {
			return role, nil
		}
	}
	return nil, nil
}

func (s *session) GuildRoles(ctx context.Context, guildID string) ([]IDomainRole, error) {
	roles, err := s.guildRolesNative(ctx, guildID)
	if err != nil {
		return nil, err
	}
	out := make([]IDomainRole, 0)
	for _, role := range roles {
		d := NewDomainRole(role.ID, role.Name, role.Color)
		s.cacheRoles.Put(d)
		out = append(out, d)
	}
	return out, nil
}

// guildRolesNative returns array of discordgo.Role instead of IDomainRole
func (s *session) guildRolesNative(ctx context.Context, guildID string) ([]*discordgo.Role, error) {
	ctx, span := instrumentation.SpanInContext(ctx, instrumentation.Vendor(s.sess.GuildRoles))
	defer span.End()

	return s.sess.GuildRoles(guildID, discordgo.WithContext(ctx))
}

func (s *session) guildRoleNative(ctx context.Context, guildID, roleID string) (*discordgo.Role, error) {
	roles, err := s.guildRolesNative(ctx, guildID)
	if err != nil {
		return nil, err
	}
	for _, role := range roles {
		if role.ID == roleID {
			return role, nil
		}
	}
	return nil, nil
}

func (s *session) GuildRoleReorder(ctx context.Context, guildID string, roles []IDomainRole) error {
	ctx, span := instrumentation.SpanInContext(ctx, instrumentation.Vendor(s.sess.GuildRoleReorder))
	defer span.End()

	nativeRoles := make([]*discordgo.Role, 0)
	for _, role := range roles {
		nativeRole, err := s.guildRoleNative(ctx, guildID, role.ID())
		if err != nil {
			return err
		}
		nativeRoles = append(nativeRoles, nativeRole)
	}
	_, err := s.sess.GuildRoleReorder(guildID, nativeRoles, discordgo.WithContext(ctx))
	return err
}

func (s *session) GuildRoleCreate(ctx context.Context, guildID string, name string, colour int) (roleID string, err error) {
	ctx, span := instrumentation.SpanInContext(ctx, instrumentation.Vendor(s.sess.GuildRoleCreate))
	defer span.End()

	roleParams := &discordgo.RoleParams{Name: name, Color: &colour}
	role, err := s.sess.GuildRoleCreate(guildID, roleParams, discordgo.WithContext(ctx))
	if err != nil {
		return "", err
	}
	return role.ID, nil
}

func (s *session) GuildRoleEdit(ctx context.Context, guildID, roleID, name string, colour int) error {
	ctx, span := instrumentation.SpanInContext(ctx, instrumentation.Vendor(s.sess.GuildRoleEdit))
	defer span.End()

	roleParams := discordgo.RoleParams{Name: name, Color: &colour}
	_, err := s.sess.GuildRoleEdit(guildID, roleID, &roleParams, discordgo.WithContext(ctx))
	if err != nil {
		return err
	}
	s.cacheRoles.Delete(roleID)
	return nil
}

func (s *session) GuildMemberRoleAdd(ctx context.Context, guildID, userID, roleID string) error {
	ctx, span := instrumentation.SpanInContext(ctx, instrumentation.Vendor(s.sess.GuildMemberRoleAdd))
	defer span.End()

	err := s.sess.GuildMemberRoleAdd(guildID, userID, roleID, discordgo.WithContext(ctx))
	if err != nil {
		return err
	}
	s.cacheMembers.Delete(userID)
	return nil
}

// guild implements IDomainGuild
type guild struct {
	id string
}

func NewDomainGuild(id string) IDomainGuild {
	return &guild{id}
}

func (g *guild) ID() string { return g.id }

// member implements IDomainMember
type member struct {
	mem   *discordgo.Member
	roles []IDomainRole
}

func NewDomainMember(mem *discordgo.Member, roles []IDomainRole) IDomainMember {
	return &member{mem, roles}
}

func (m *member) Guild() IDomainGuild  { return NewDomainGuild(m.mem.GuildID) }
func (m *member) UserID() string       { return m.mem.User.ID }
func (m *member) Nick() string         { return m.mem.Nick }
func (m *member) Username() string     { return m.mem.User.Username + "#" + m.mem.User.Discriminator }
func (m *member) Roles() []IDomainRole { return m.roles }
func (m *member) CacheKey() string     { return m.UserID() }

// colourRole implements IDomainRole.
type colourRole struct {
	roleID string
	name   string
	colour *Colour
}

func NewDomainRole(id, name string, colour int) IDomainRole {
	col := (&Colour{}).FromDecimal(colour)
	return &colourRole{id, name, col}
}

func (r *colourRole) ID() string       { return r.roleID }
func (r *colourRole) Name() string     { return r.name }
func (r *colourRole) Colour() *Colour  { return r.colour }
func (r *colourRole) CacheKey() string { return r.ID() }
