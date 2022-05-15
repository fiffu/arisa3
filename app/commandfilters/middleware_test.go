package commandfilters

import (
	"testing"

	"github.com/fiffu/arisa3/app/types"
	"github.com/stretchr/testify/assert"
)

func True(types.ICommandEvent) bool {
	return true
}

func False(types.ICommandEvent) bool {
	return false
}

func Panic(types.ICommandEvent) bool {
	panic("oh no")
}

func Test(t *testing.T) {
	testCases := []struct {
		desc   string
		mw     *Middleware
		expect bool
	}{
		{
			desc:   "false AND true OR false == false",
			mw:     NewMiddleware(False).And(True).Or(False),
			expect: false,
		},
		{
			desc:   "false OR true AND false == false",
			mw:     NewMiddleware(False).Or(True).And(False),
			expect: false,
		},
		{
			desc:   "false OR false OR true",
			mw:     NewMiddleware(False).Or(False).Or(True),
			expect: true,
		},
		{
			desc:   "false AND panic should short-circuit",
			mw:     NewMiddleware(False).And(Panic),
			expect: false,
		},
		{
			desc:   "true OR panic should short-circuit",
			mw:     NewMiddleware(True).Or(Panic),
			expect: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			actual := tc.mw.Exec(nil)
			assert.Equal(t, tc.expect, actual)
		})
	}
}
