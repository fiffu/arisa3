package instrumentation

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_discordAPIMatcher(t *testing.T) {
	testCases := []struct {
		path   string
		expect string
	}{
		// Matches full path
		{
			path:   "discord.com/api/v9/applications/964085462748774401/commands",
			expect: "discord.com/api/.+/applications/.+/commands",
		},
		{
			path:   "discord.com/api/v9/interactions/1176038219758968842/aW50ZXJhY3Rpb246MTE3NjAzODIxOTc1ODk2ODg0Mjp5TVdNQmJrZUNmVnFERjludzlKSG1VRndicDN0T1ZQeVFFWFI1WW5vN2U5Y2RKTGZVamZnbnFKQWRxS2dabzN4QU9qcjlWOTFIc2RIUGUwbUhrQWtjRWI1YjJiMHdGRzNVSjlHc0I0Q0JaTjc1dHNETnh4TU90ZnpsdGlGOVQ5Sg/callback",
			expect: "discord.com/api/.+/interactions/.+/.+/callback",
		},
		{
			path:   "discord.com/api/v9/guilds/294027089760223232/roles/620080566242508810",
			expect: "discord.com/api/.+/guilds/.+/roles/.+",
		},
		{
			path:   "discord.com/api/v9/guilds/294027089760223232/members/176045537215119360",
			expect: "discord.com/api/.+/guilds/.+/members/.+",
		},
		{
			path:   "discord.com/api/v9/guilds/294027089760223232/roles",
			expect: "discord.com/api/.+/guilds/.+/roles",
		},
		// Matches on base path
		{
			path:   "discord.com/api/v9/interactions/some-random-path",
			expect: "discord.com/api/v9/interactions/some-random-path",
		},
		{
			path:   "discord.com/api/v9/guilds/some-random-path",
			expect: "discord.com/api/v9/guilds/some-random-path",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.expect, func(t *testing.T) {
			actual := discordAPIMatcher.Match(context.Background(), tc.path)
			assert.Equal(t, tc.expect, actual)
		})
	}
}
