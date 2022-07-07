package rng

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/fiffu/arisa3/app/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func newTestPokiesRequest(ctrl *gomock.Controller, size int, ok bool) types.ICommandEvent {
	req := types.NewMockICommandEvent(ctrl)
	args := types.NewMockIArgs(ctrl)
	args.EXPECT().Int(gomock.Any()).Return(size, ok)
	req.EXPECT().Args().Return(args)
	return req
}

func Test_parseGrid(t *testing.T) {
	testCases := []struct {
		desc                   string
		size                   int
		ok                     bool
		expectRows, expectCols int
		expectTooBig           bool
	}{
		{
			desc:       "Default, no arg passed",
			ok:         false,
			expectRows: pokiesDefaultRows,
			expectCols: pokiesDefaultCols,
		},
		{
			desc:       "Zero size should yield default",
			size:       0,
			ok:         true,
			expectRows: pokiesDefaultRows,
			expectCols: pokiesDefaultCols,
		},
		{
			desc:       "Negative size should yield default",
			size:       -1,
			ok:         true,
			expectRows: pokiesDefaultRows,
			expectCols: pokiesDefaultCols,
		},
		{
			desc:       "Arg is passed",
			size:       4,
			ok:         true,
			expectRows: 4,
			expectCols: 4,
		},
		{
			desc:       "8 should be accepted",
			size:       8,
			ok:         true,
			expectRows: 8,
			expectCols: 8,
		},
		{
			desc:         "Arg passed is too big",
			size:         9,
			ok:           true,
			expectTooBig: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			req := newTestPokiesRequest(ctrl, tc.size, tc.ok)
			rows, cols, ok := parseGrid(req)
			if tc.expectTooBig {
				assert.False(t, ok)
			} else {
				assert.True(t, ok)
				assert.Equal(t, tc.expectRows, rows)
				assert.Equal(t, tc.expectCols, cols)
			}
		})
	}
}

func Test_buildGrid_tooBig(t *testing.T) {
	e := &discordgo.Emoji{Name: "aaaaaaaaaa", ID: "12345678"}
	_, ok := buildGrid(100, 100, []*discordgo.Emoji{e})
	assert.False(t, ok)
}
