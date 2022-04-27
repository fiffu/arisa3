package colours

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Decimal(t *testing.T) {
	// ffff00
	type testCase struct {
		r, g, b    int     // input [0, 255]
		rx, gx, bx float64 // expect [0, 1]
	}
	for _, tc := range []testCase{
		{
			0, 0, 0,
			0, 0, 0,
		},
		{
			255, 0, 0,
			1, 0, 0,
		},
		{
			255, 127, 0,
			1, 0.5, 0,
		},
	} {
		// 0.1% accuracy
		epsilon255 := 255.0 / 1000
		epsilon1 := 1.0 / 1000

		dec := (tc.r << 16) + (tc.g << 8) + tc.b
		col := (&Colour{}).FromDecimal(dec)
		assert.LessOrEqual(t, math.Abs(tc.rx-col.R), epsilon255)
		assert.LessOrEqual(t, math.Abs(tc.gx-col.G), epsilon255)
		assert.LessOrEqual(t, math.Abs(tc.bx-col.B), epsilon255)

		dec2 := col.ToDecimal()
		assert.LessOrEqual(t, math.Abs(float64(dec2-dec)), epsilon1)
	}
}

func Test_Random(t *testing.T) {
	for i := 0; i < 1000; i++ {
		c := (&Colour{}).Random()
		// assert r, g, b components are in range [0,1]
		lo, hi := 0.0, 1.0
		for _, f := range []float64{c.R, c.G, c.B} {
			assert.LessOrEqual(t, lo, f)
			assert.LessOrEqual(t, f, hi)
		}
	}
}
