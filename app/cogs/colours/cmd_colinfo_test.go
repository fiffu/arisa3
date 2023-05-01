package colours

import (
	"strings"
	"testing"
	"time"

	"github.com/fiffu/arisa3/lib/functional"
	"github.com/stretchr/testify/assert"
)

func unix(n int64) time.Time {
	return time.Unix(n, 0)
}

func concat(s ...string) string {
	return strings.Join(s, "")
}

func colours(hxs ...int) []*Colour {
	ret := make([]*Colour, len(hxs))
	for i, hx := range hxs {
		ret[i] = (&Colour{}).FromDecimal(hx)
	}
	return ret
}

func Test_partitionColours(t *testing.T) {
	testCases := []struct {
		desc     string
		records  []*ColoursLogRecord
		expect   []*Colour
		interval int
	}{
		// All cases will use 50 seconds span, from unix(0) to unix(50)
		{
			desc:     "single log record",
			interval: 10,
			records: []*ColoursLogRecord{
				{TStamp: unix(20), ColourHex: "ff0000"},
			},
			expect: colours(0xff0000, 0xff0000, 0xff0000, 0xff0000, 0xff0000),
		},
		{
			desc:     "diff colour each step",
			interval: 10,
			records: []*ColoursLogRecord{
				{TStamp: unix(00), ColourHex: "ff0000"},
				{TStamp: unix(10), ColourHex: "ff000f"},
				{TStamp: unix(20), ColourHex: "ff00ff"},
				{TStamp: unix(30), ColourHex: "ff0fff"},
				{TStamp: unix(40), ColourHex: "ffffff"},
			},
			expect: colours(0xff0000, 0xff000f, 0xff00ff, 0xff0fff, 0xffffff),
		},
		{
			desc:     "2 colours across all steps",
			interval: 10,
			records: []*ColoursLogRecord{
				{TStamp: unix(00), ColourHex: "ffffff"},
				{TStamp: unix(20), ColourHex: "000000"},
			},
			expect: colours(0xffffff, 0xffffff, 0x000000, 0x000000, 0x000000),
		},
		{
			desc:     "returns the last of multiple colours within each span",
			interval: 10,
			records: []*ColoursLogRecord{
				{TStamp: unix(00), ColourHex: "110000"},
				{TStamp: unix(11), ColourHex: "11000f"},
				{TStamp: unix(12), ColourHex: "1100ff"},
				{TStamp: unix(20), ColourHex: "220000"},
				{TStamp: unix(21), ColourHex: "22000f"},
				{TStamp: unix(22), ColourHex: "2200ff"},
				{TStamp: unix(30), ColourHex: "330000"},
				{TStamp: unix(31), ColourHex: "33000f"},
				{TStamp: unix(32), ColourHex: "3300ff"},
			},
			expect: colours(0x110000, 0x1100ff, 0x2200ff, 0x3300ff, 0x3300ff),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			actual := partitionColours(
				&History{tc.records, unix(0), unix(50)},
				time.Duration(tc.interval)*time.Second,
			)
			assert.Equal(t,
				functional.Map(tc.expect, func(c *Colour) string { return c.ToHexcode() }),
				functional.Map(actual, func(c *Colour) string { return c.ToHexcode() }),
			)
		})
	}
}

func Test_formatColInfo(t *testing.T) {
	now := time.Unix(1682000000, 0)

	testCases := []struct {
		desc            string
		rerollCDEndTime time.Time
		lastMutateTime  time.Time
		lastFrozenTime  time.Time
		history         *History

		expectDesc string
	}{
		{
			desc:            "reroll CD elapsed",
			rerollCDEndTime: now.Add(1 * time.Hour),
			expectDesc: concat(
				"**Reroll cooldown:**\n1 hour\n\n",
				"**Last mutate:**\n_(Never)_",
			),
		},
		{
			desc:           "have mutated, 1 hour before",
			lastMutateTime: now.Add(-1 * time.Hour),
			expectDesc: concat(
				"**Reroll cooldown:**\n_(No cooldown, reroll available)_\n\n",
				"**Last mutate:**\n1 hour ago",
			),
		},
		{
			desc:           "have mutated, exactly now",
			lastMutateTime: now,
			expectDesc: concat(
				"**Reroll cooldown:**\n_(No cooldown, reroll available)_\n\n",
				"**Last mutate:**\nMoments ago",
			),
		},
		{
			desc:           "have mutated, but frozen",
			lastMutateTime: now.Add(-1 * time.Hour),
			lastFrozenTime: now.Add(-1 * time.Hour),
			expectDesc: concat(
				"**Reroll cooldown:**\n_(No cooldown, reroll available)_\n\n",
				"**Last mutate:**\n1 hour ago\n",
				"Frozen 1 hour ago",
			),
		},
		{
			desc:           "never mutated, but frozen",
			lastMutateTime: Never,
			lastFrozenTime: now.Add(-1 * time.Hour),
			expectDesc: concat(
				"**Reroll cooldown:**\n_(No cooldown, reroll available)_\n\n",
				"**Last mutate:**\n_(Never)_\n",
				"Frozen 1 hour ago",
			),
		},
		{
			desc:           "have mutated, have history",
			lastMutateTime: now.Add(-1 * time.Hour),
			history: &History{
				start: unix(00),
				end:   unix(40),
				records: []*ColoursLogRecord{
					{TStamp: unix(00), ColourHex: "ffffff"},
					{TStamp: unix(20), ColourHex: "000000"},
				},
			},
			expectDesc: concat(
				"**Reroll cooldown:**\n_(No cooldown, reroll available)_\n\n",
				"**Last mutate:**\n1 hour ago\n\n",
				"**Image history, newest → oldest:**",
			),
		},
	}
	cog := &Cog{cfg: &Config{MutateCooldownMins: 1}}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			reply, err := cog.formatColInfo(
				now,
				tc.rerollCDEndTime,
				tc.lastMutateTime,
				tc.lastFrozenTime,
				tc.history,
			)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectDesc, reply.desc)
		})
	}
}

func Test_horizontalPartitionImage(t *testing.T) {
	r, g, b := 0xff0000, 0x00ff00, 0x0000ff
	hpi := horizontalPartitionImage{
		partitions:      colours(r, g, b),
		partitionWidth:  3,
		partitionHeight: 2,
	}

	bitmap := [][]int{
		{r, r, r, g, g, g, b, b, b},
		{r, r, r, g, g, g, b, b, b},
	}
	for y, row := range bitmap {
		for x, pixel := range row {
			expect := (&Colour{}).FromDecimal(pixel)
			actual := hpi.At(x, y)

			xr, xg, xb, xa := expect.RGBA()
			ar, ag, ab, aa := actual.RGBA()
			assert.Equal(t, xr, ar)
			assert.Equal(t, xg, ag)
			assert.Equal(t, xb, ab)
			assert.Equal(t, xa, aa)
		}
	}
}

func Test_makeColHistoryImg(t *testing.T) {
	h := &History{
		start: unix(00),
		end:   unix(40),
		records: []*ColoursLogRecord{
			{TStamp: unix(00), ColourHex: "ffffff"},
			{TStamp: unix(20), ColourHex: "000000"},
		},
	}

	file, fileExt, fileContent, err := makeColHistoryImg(h, 20*time.Second)
	assert.NoError(t, err)
	assert.NotNil(t, file)
	assert.Equal(t, "png", fileExt)
	assert.Equal(t, "image/png", fileContent)
}
