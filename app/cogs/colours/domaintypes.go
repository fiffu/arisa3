package colours

// domaintypes.go describes the core types and interfaces of the Colour Roles domain.
// Any domain logic should operate on these primitives.

import (
	"fmt"
	"time"

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

// ToHexcode returns the Colour in HTML-encoded hexcode.
func (c *Colour) ToHexcode() string {
	// Stolen from https://github.com/gerow/go-color/blob/master/color.go
	delta := 1 / 512.0 // to make truncation round to nearest number instead of flooring
	return fmt.Sprintf(
		"%02x%02x%02x",
		byte((c.R+delta)*255),
		byte((c.G+delta)*255),
		byte((c.B+delta)*255),
	)
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

// IColoursDomain describes the colour roles domain
type IColoursDomain interface {
	GetLastMutate(IDomainMember) (time.Time, error)
	GetLastReroll(IDomainMember) (time.Time, error)
	GetLastFrozen(IDomainMember) (time.Time, error)
	Mutate(IDomainMember) (*Colour, error)
	Reroll(IDomainMember) (*Colour, error)
	MakeRoleName(IDomainMember) string
}

// IDomainMember describes information that IColoursDomain derives from discordgo.User.
type IDomainMember interface {
	UserID() string
	UserName() string
	Nick() string
	ColourRole() IDomainRole
}

// IDomainRole describes information that IColoursDomain derives from discordgo.Role.
type IDomainRole interface {
	Name() string
	Colour() *Colour
}

// IDomainRepository describes methods that IColoursDomain uses to fetch/store data.
type IDomainRepository interface {
	FetchUserState(IDomainMember, Reason) (time.Time, error)
	UpdateMutate(IDomainMember, *Colour) error
	UpdateReroll(IDomainMember, *Colour) error
	UpdateFreeze(IDomainMember) error
	UpdateUnfreeze(IDomainMember) error
}
