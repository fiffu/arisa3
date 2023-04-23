package colours

// types.go implements the interfaces defined by interfaces.go.

import (
	"time"

	"github.com/bwmarrin/discordgo"
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

func (s *session) GuildMember(guildID, userID string) (IDomainMember, error) {
	// Cache lookup
	if cached, ok := s.cacheMembers.Peek(userID); ok {
		return cached, nil
	}

	// Query API for guild member
	mem, err := s.sess.GuildMember(guildID, userID)
	if err != nil {
		return nil, err
	}

	// Query for guild roles.
	allRoles, err := s.GuildRoles(mem.GuildID)
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

func (s *session) GuildRole(guildID, roleID string) (IDomainRole, error) {
	if cached, ok := s.cacheRoles.Peek(roleID); ok {
		return cached, nil
	}
	roles, err := s.GuildRoles(guildID)
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

func (s *session) GuildRoles(guildID string) ([]IDomainRole, error) {
	roles, err := s.guildRolesNative(guildID)
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
func (s *session) guildRolesNative(guildID string) ([]*discordgo.Role, error) {
	return s.sess.GuildRoles(guildID)
}

func (s *session) guildRoleNative(guildID, roleID string) (*discordgo.Role, error) {
	roles, err := s.guildRolesNative(guildID)
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

func (s *session) GuildRoleReorder(guildID string, roles []IDomainRole) error {
	nativeRoles := make([]*discordgo.Role, 0)
	for _, role := range roles {
		nativeRole, err := s.guildRoleNative(guildID, role.ID())
		if err != nil {
			return err
		}
		nativeRoles = append(nativeRoles, nativeRole)
	}
	_, err := s.sess.GuildRoleReorder(guildID, nativeRoles)
	return err
}

func (s *session) GuildRoleCreate(guildID string) (roleID string, err error) {
	role, err := s.sess.GuildRoleCreate(guildID)
	if err != nil {
		return "", err
	}
	return role.ID, nil
}

func (s *session) GuildRoleEdit(guildID, roleID, name string, colour int) error {
	_, err := s.sess.GuildRoleEdit(
		guildID, roleID, name, colour,
		false, 0, false,
	)
	if err != nil {
		return err
	}
	s.cacheRoles.Delete(roleID)
	return nil
}

func (s *session) GuildMemberRoleAdd(guildID, userID, roleID string) error {
	err := s.sess.GuildMemberRoleAdd(guildID, userID, roleID)
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
