package colours

import (
	"encoding/hex"
	"fmt"
	"math"

	"github.com/fiffu/arisa3/lib"
)

// Colour encodes a Colour Role's RGB value.
type Colour struct {
	R, G, B float64
}

func (c *Colour) String() string {
	return c.ToHexcode()
}

func (c *Colour) scale255() (r, g, b int) {
	delta := 1 / 512.0 // to make truncation round to nearest number instead of flooring
	r = int((c.R + delta) * 255)
	g = int((c.G + delta) * 255)
	b = int((c.B + delta) * 255)
	return
}

func (c *Colour) scale0xFFFF() (r, g, b uint32) {
	delta := 0.0
	r = uint32((c.R + delta) * 0xFFFF)
	g = uint32((c.G + delta) * 0xFFFF)
	b = uint32((c.B + delta) * 0xFFFF)
	return
}

// RGBA implements interface color.Color of standard lib.
func (c *Colour) RGBA() (uint32, uint32, uint32, uint32) {
	r, g, b := c.scale0xFFFF()
	return r, g, b, 0xFFFF
}

func (c *Colour) ToDecimal() int {
	r, g, b := c.scale255()
	return (r << 16) + (g << 8) + b
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
	r := colour >> 16
	colour -= r * 65536

	g := colour / 256
	colour -= g * 256

	b := colour
	return &Colour{float64(r) / 255, float64(g) / 255, float64(b) / 255}
}

// FromHSV returns a new instance of Colour, converting from HSV input to RGB colour space.
func (c *Colour) FromHSV(h, s, v float64) *Colour {
	r, g, b := hsvToRGB(h, s, v)
	return &Colour{r, g, b}
}

// FromRGBHex returns a new instance of Colour, converting hex-encoded RGB.
func (c *Colour) FromRGBHex(rgbHex string) *Colour {
	byteArr, err := hex.DecodeString(rgbHex)
	// fmt.Println(byteArr)
	if err != nil || len(byteArr) != 3 {
		return &Colour{0, 0, 0}
	}
	return &Colour{
		R: float64(byteArr[0]) / 255.0,
		G: float64(byteArr[1]) / 255.0,
		B: float64(byteArr[2]) / 255.0,
	}
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
		if lib.CoinFlip() {
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

// hsvToRGB converts from HSV tuple to RGB colour space.
func hsvToRGB(h, s, v float64) (r, g, b float64) {
	// Adapted from Python stdlib
	// https://github.com/python/cpython/blob/3.10/Lib/colorsys.py

	if s == 0 {
		return v, v, v
	}
	i := int(math.Floor(h * 6))
	f := h*6 - float64(i)
	p := v * (1.0 - s)
	q := v * (1.0 - s*f)
	t := v * (1.0 - s*(1.0-f))

	switch i % 6 {
	case 0:
		return v, t, p
	case 1:
		return q, v, p
	case 2:
		return p, v, t
	case 3:
		return p, q, v
	case 4:
		return t, p, v
	case 5:
		return v, p, q
	default:
		// Cannot get here
		return 0, 0, 0
	}
}
