package stonks

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/fiffu/arisa3/app/types"

	dgo "github.com/bwmarrin/discordgo"
)

// customEmoji matches the text Discord sends for a custom emoji, such as <:birbcoin:123> or <a:birb:123>.
var customEmoji = regexp.MustCompile(`^<a?:\w+:(\d+)>$`)

func parseEmojiID(s string) (string, bool) {
	m := customEmoji.FindStringSubmatch(strings.TrimSpace(s))
	if m == nil {
		return "", false
	}
	return m[1], true
}

func (c *Cog) setupCommand() *types.Command {
	return types.NewCommand("stonks-setup").ForChat().
		Desc("Set the market channel and the birbcoin emoji for this server").
		MemberPermissions(dgo.PermissionManageServer).
		Options(
			types.NewOption("channel").Channel().
				ChannelType([]dgo.ChannelType{dgo.ChannelTypeGuildText}).
				Desc("Channel for market updates").
				Required(),
			types.NewOption("emoji").String().
				Desc("Custom emoji. A react with it gives the message author one birbcoin").
				Required(),
		).
		Handler(c.setup)
}

func (c *Cog) setup(ctx context.Context, req types.ICommandEvent) error {
	i := req.Interaction()
	if i.GuildID == "" || i.Member == nil {
		return req.Respond(ctx, types.NewResponse().Ephemeral().Content("Run this command in a server."))
	}
	// Server admins can grant the command to roles without Manage Server, so check here too.
	if i.Member.Permissions&dgo.PermissionManageServer == 0 {
		return req.Respond(ctx, types.NewResponse().Ephemeral().Content("You need the Manage Server permission."))
	}

	// For a channel option, Discord sends the channel ID as a string.
	channelID, _ := req.Args().String("channel")
	emoji, _ := req.Args().String("emoji")
	emojiID, ok := parseEmojiID(emoji)
	if channelID == "" || !ok {
		return req.Respond(ctx, types.NewResponse().Ephemeral().
			Content("Pick a text channel and a custom emoji. Built-in emojis such as 🪙 do not work."))
	}

	if err := c.wallets.repo.SaveGuildSettings(ctx, i.GuildID, emojiID, channelID); err != nil {
		return fmt.Errorf("save guild settings: %w", err)
	}
	c.wallets.emojis.Set(i.GuildID, emojiID)

	embed := types.NewEmbed().
		Title("Stonks setup saved").
		Field("Market channel", fmt.Sprintf("<#%s>", channelID), true).
		Field("Coin emoji", strings.TrimSpace(emoji), true)
	return req.Respond(ctx, types.NewResponse().Embeds(embed))
}
