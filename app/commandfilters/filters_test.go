package commandfilters

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/fiffu/arisa3/app/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	desc     string
	isMember bool
	isAdmin  bool
	expect   bool
}

type middleware func(ev types.ICommandEvent) bool

func runTest(t *testing.T, mw middleware, tc testCase) {
	t.Helper()

	// Mocking instrumentation
	ctrl := gomock.NewController(t)

	mockEvt := types.NewMockICommandEvent(ctrl)
	mockInteraction := &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{},
	}

	// Inject mock return value
	if tc.isMember {
		mockInteraction.Interaction.Member = &discordgo.Member{}
	}
	if tc.isAdmin {
		adminFlag := int64(discordgo.PermissionAdministrator)
		mockInteraction.Interaction.Member.Permissions |= adminFlag
	}
	mockEvt.EXPECT().Interaction().
		AnyTimes().
		Return(mockInteraction)

	// Exec and assert
	actual := mw(mockEvt)
	assert.Equal(t, tc.expect, actual)
}

func Test_IsGuildAdmin(t *testing.T) {
	testCases := []testCase{
		{
			desc:     "non-guild member should return false",
			isMember: false,
			expect:   false,
		},
		{
			desc:     "non-admin guild member should return false",
			isMember: true,
			expect:   false,
		},
		{
			desc:     "admin guild member should return true",
			isMember: true,
			isAdmin:  true,
			expect:   true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			runTest(t, IsGuildAdmin, tc)
		})
	}
}

func Test_IsFromGuild(t *testing.T) {
	testCases := []testCase{
		{
			desc:     "non-guild member should return false",
			isMember: false,
			expect:   false,
		},
		{
			desc:     "guild member should return true",
			isMember: true,
			expect:   true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			runTest(t, IsFromGuild, tc)
		})
	}
}
