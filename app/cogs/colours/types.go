package colours

// types.go implements the interfaces defined by interfaces.go.

import (
	"errors"
	"fmt"
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

// NoState indicates user has never interacted with the Colour Roles domain.
var NoState = &ColourState{}

// Never indicates that user has state within the domain, but no records for the given Reason.
var Never = time.Time{} // the zero value

// Colour encodes a Colour Role's RGB value.
type Colour struct {
	R, G, B float64
}

func (c *Colour) scale255() (r, g, b int) {
	delta := 1 / 512.0 // to make truncation round to nearest number instead of flooring
	r = int((c.R + delta) * 255)
	g = int((c.R + delta) * 255)
	b = int((c.R + delta) * 255)
	return
}

func (c *Colour) ToDecimal() int {
	r, g, b := c.scale255()
	return r<<16 + g<<8 + b
}

// ToHexcode returns the Colour in HTML-encoded hexcode.
func (c *Colour) ToHexcode() string {
	// Stolen from https://github.com/gerow/go-color/blob/master/color.go
	r, g, b := c.scale255()
	return fmt.Sprintf(
		"%02x%02x%02x",
		byte(r), byte(g), byte(b),
	)
}

func (c *Colour) FromDecimal(colour int) *Colour {
	r, g, b := lib.DecimalToRGB(colour)
	return &Colour{r, g, b}
}

// FromHSV returns a new instance of Colour, converting from HSV input to RGB colour space.
func (c *Colour) FromHSV(h, s, v float64) *Colour {
	r, g, b := lib.HSVtoRGB(h, s, v)
	return &Colour{r, g, b}
}

// Random returns a new instance of Colour with freshly-seeded values.
func (c *Colour) Random() *Colour {
	return c.FromHSV(
		lib.UniformRange(0, 1),       // any hue
		lib.UniformRange(0.55, 0.85), // less variation on saturation
		lib.UniformRange(0.50, 0.90), // more variation on lightness
	)
}

// Nudge returns a copy of the current Colour with very slightly adjusted values.
func (c *Colour) Nudge() *Colour {
	step := func() float64 {
		distance := lib.UniformRange(0.08, 0.15)
		if lib.ChooseBool() {
			distance *= -1
		}
		return distance
	}
	clamp := lib.Clamper(0, 1)
	return &Colour{
		clamp(c.R + step()),
		clamp(c.G + step()),
		clamp(c.B + step()),
	}
}

// ColourState models a participant's state in the Colour Roles domain.
type ColourState struct {
	UserID     string
	LastFrozen time.Time
	LastMutate time.Time
	LastReroll time.Time
}

// session implements IDomainSession
type session struct {
	sess         *discordgo.Session
	cacheMembers lib.ICache
	cacheRoles   lib.ICache
}

func NewDomainSession(sess *discordgo.Session, cacheTimeout time.Duration) IDomainSession {
	return &session{
		sess, lib.NewMemoryCache(), lib.NewMemoryCache(),
	}
}

func (s *session) GuildMember(guildID, userID string) (IDomainMember, error) {
	// Cache lookup
	if cached, ok := s.cacheMembers.Peek(userID); ok {
		if mem, ok := (cached.CacheData()).(*member); ok {
			return mem, nil
		} else {
			return nil, errors.New("error coercing cached member")
		}
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
		if role, ok := (cached.CacheData()).(*colourRole); ok {
			return role, nil
		} else {
			return nil, errors.New("error coercing cached colourRole")
		}
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

func (s *session) GuildRoleEdit(guildID, roleID, name string, color int) error {
	_, err := s.sess.GuildRoleEdit(
		guildID, roleID, name, color,
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

func (m *member) Guild() IDomainGuild          { return NewDomainGuild(m.mem.GuildID) }
func (m *member) UserID() string               { return m.mem.User.ID }
func (m *member) Nick() string                 { return m.mem.Nick }
func (m *member) Username() string             { return m.mem.User.Username + m.mem.User.Discriminator }
func (m *member) Roles() []IDomainRole         { return m.roles }
func (m *member) CacheKey() string             { return m.UserID() }
func (m *member) CacheData() interface{}       { return m }
func (m *member) CacheDuration() time.Duration { return time.Hour + 1 }

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

func (r *colourRole) ID() string                   { return r.roleID }
func (r *colourRole) Name() string                 { return r.name }
func (r *colourRole) Colour() *Colour              { return r.colour }
func (r *colourRole) CacheKey() string             { return r.ID() }
func (r *colourRole) CacheData() interface{}       { return r }
func (r *colourRole) CacheDuration() time.Duration { return time.Hour * 1 }
