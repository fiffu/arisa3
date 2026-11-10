package stonks

import (
	"context"
	"testing"

	"github.com/fiffu/arisa3/app/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	dgo "github.com/bwmarrin/discordgo"
)

func newSetupEvent(ctrl *gomock.Controller, perms int64, channel, emoji string, got *types.ICommandResponse) types.ICommandEvent {
	args := types.NewMockIArgs(ctrl)
	args.EXPECT().String("channel").Return(channel, channel != "").AnyTimes()
	args.EXPECT().String("emoji").Return(emoji, emoji != "").AnyTimes()

	evt := types.NewMockICommandEvent(ctrl)
	evt.EXPECT().Interaction().Return(&dgo.InteractionCreate{Interaction: &dgo.Interaction{
		GuildID: guild,
		Member:  &dgo.Member{User: &dgo.User{ID: author}, Permissions: perms},
	}}).AnyTimes()
	evt.EXPECT().Args().Return(args).AnyTimes()
	evt.EXPECT().Respond(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, r types.ICommandResponse) error {
		*got = r
		return nil
	})
	return evt
}

func Test_setup(t *testing.T) {
	cases := []struct {
		name      string
		perms     int64
		channel   string
		emoji     string
		saved     bool
		ephemeral bool
	}{
		{"saves settings", dgo.PermissionManageServer, channel, "<:birbcoin:" + coin + ">", true, false},
		{"needs Manage Server", dgo.PermissionSendMessages, channel, "<:birbcoin:" + coin + ">", false, true},
		{"rejects a built-in emoji", dgo.PermissionManageServer, channel, "🪙", false, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := newFakeWalletRepo()
			cog := &Cog{wallets: &wallets{repo: repo, emojis: &guildEmojis{}}}
			var resp types.ICommandResponse

			require.NoError(t, cog.setup(context.Background(), newSetupEvent(ctrl, tc.perms, tc.channel, tc.emoji, &resp)))

			emojiID, cached := cog.wallets.emojis.Get(guild)
			assert.Equal(t, tc.saved, cached)
			if tc.saved {
				assert.Equal(t, coin, emojiID)
				assert.Equal(t, [2]string{coin, channel}, repo.settings[guild])
			} else {
				assert.Empty(t, repo.settings)
			}
			require.NotNil(t, resp)
			assert.Equal(t, tc.ephemeral, resp.Data().Data.Flags&dgo.MessageFlagsEphemeral != 0)
		})
	}
}
