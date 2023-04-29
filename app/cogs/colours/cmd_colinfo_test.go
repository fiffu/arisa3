package colours

import (
	"testing"
	"time"

	"github.com/fiffu/arisa3/lib/functional"
	"github.com/stretchr/testify/assert"
)

func Test_partitionColours(t *testing.T) {
	unix := func(n int64) time.Time {
		return time.Unix(n, 0)
	}
	colours := func(hxs ...int) []*Colour {
		ret := make([]*Colour, len(hxs))
		for i, hx := range hxs {
			ret[i] = (&Colour{}).FromDecimal(hx)
		}
		return ret
	}

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
