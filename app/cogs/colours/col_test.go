package colours

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func lineJoin(s ...string) string {
	return strings.Join(s, "")
}

func Test_formatColInfo(t *testing.T) {
	now := time.Unix(1682000000, 0)

	testCases := []struct {
		desc            string
		rerollCDEndTime time.Time
		lastMutateTime  time.Time
		lastFrozenTime  time.Time

		expectDesc string
	}{
		{
			desc:            "reroll CD elapsed",
			rerollCDEndTime: now.Add(1 * time.Hour),
			expectDesc: lineJoin(
				"**Reroll cooldown:**\n1 hour\n\n",
				"**Last mutate:**\n_(Never)_",
			),
		},
		{
			desc:           "have mutated, 1 hour before",
			lastMutateTime: now.Add(-1 * time.Hour),
			expectDesc: lineJoin(
				"**Reroll cooldown:**\n_(No cooldown, reroll available)_\n\n",
				"**Last mutate:**\n1 hour ago",
			),
		},
		{
			desc:           "have mutated, exactly now",
			lastMutateTime: now,
			expectDesc: lineJoin(
				"**Reroll cooldown:**\n_(No cooldown, reroll available)_\n\n",
				"**Last mutate:**\nMoments ago",
			),
		},
		{
			desc:           "have mutated, but frozen",
			lastMutateTime: now.Add(-1 * time.Hour),
			lastFrozenTime: now.Add(-1 * time.Hour),
			expectDesc: lineJoin(
				"**Reroll cooldown:**\n_(No cooldown, reroll available)_\n\n",
				"**Last mutate:**\n1 hour ago\n",
				"**Colour has been frozen for:**\n1 hour",
			),
		},
		{
			desc:           "never mutated, but frozen",
			lastMutateTime: Never,
			lastFrozenTime: now.Add(-1 * time.Hour),
			expectDesc: lineJoin(
				"**Reroll cooldown:**\n_(No cooldown, reroll available)_\n\n",
				"**Last mutate:**\n_(Never)_\n",
				"**Colour has been frozen for:**\n1 hour",
			),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			reply := (&Cog{}).formatColInfo(
				now,
				tc.rerollCDEndTime,
				tc.lastMutateTime,
				tc.lastFrozenTime,
			)
			assert.Equal(t, tc.expectDesc, reply)
		})
	}
}
