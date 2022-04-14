package cogs

import (
	"context"
	"encoding/json"

	"github.com/bwmarrin/discordgo"
)

// Util is a method set available to all cogs
type util struct{ cog ICog }

func (c *util) ParseConfig(in interface{}, out interface{}) error {
	bytes, err := json.Marshal(in)
	if err != nil {
		return err // errors.New("unable to parse cog config")
	}
	err = json.Unmarshal(bytes, out)
	if err != nil {
		return err //
	}
	return nil
}

// Create context based on the event type
func (u *util) Context(a IApp, event interface{}) context.Context {
	ctx := a.ContextWithValue(context.Background(), "event", u.ParseEvent(event))
	ctx = a.ContextWithValue(ctx, "cog", u.cog.Name())
	return ctx
}

// ParseEvent returns the event name based on event handler structs provided by discordgo.
// Courtesy of regex match on discordgo@v0.24.0/eventhandlers.go
func (u *util) ParseEvent(event interface{}) string {
	switch event.(type) {
	case *discordgo.ChannelCreate:
		return "channelCreate"
	case *discordgo.ChannelDelete:
		return "channelDelete"
	case *discordgo.ChannelPinsUpdate:
		return "channelPinsUpdate"
	case *discordgo.ChannelUpdate:
		return "channelUpdate"
	case *discordgo.Connect:
		return "connect"
	case *discordgo.Disconnect:
		return "disconnect"
	case *discordgo.Event:
		return "event"
	case *discordgo.GuildBanAdd:
		return "guildBanAdd"
	case *discordgo.GuildBanRemove:
		return "guildBanRemove"
	case *discordgo.GuildCreate:
		return "guildCreate"
	case *discordgo.GuildDelete:
		return "guildDelete"
	case *discordgo.GuildEmojisUpdate:
		return "guildEmojisUpdate"
	case *discordgo.GuildIntegrationsUpdate:
		return "guildIntegrationsUpdate"
	case *discordgo.GuildScheduledEventCreate:
		return "guildScheduledEventCreate"
	case *discordgo.GuildScheduledEventUpdate:
		return "guildScheduledEventUpdate"
	case *discordgo.GuildScheduledEventDelete:
		return "guildScheduledEventDelete"
	case *discordgo.GuildMemberAdd:
		return "guildMemberAdd"
	case *discordgo.GuildMemberRemove:
		return "guildMemberRemove"
	case *discordgo.GuildMemberUpdate:
		return "guildMemberUpdate"
	case *discordgo.GuildMembersChunk:
		return "guildMembersChunk"
	case *discordgo.GuildRoleCreate:
		return "guildRoleCreate"
	case *discordgo.GuildRoleDelete:
		return "guildRoleDelete"
	case *discordgo.GuildRoleUpdate:
		return "guildRoleUpdate"
	case *discordgo.GuildUpdate:
		return "guildUpdate"
	case *discordgo.InteractionCreate:
		return "interactionCreate"
	case *discordgo.InviteCreate:
		return "inviteCreate"
	case *discordgo.InviteDelete:
		return "inviteDelete"
	case *discordgo.MessageAck:
		return "messageAck"
	case *discordgo.MessageCreate:
		return "messageCreate"
	case *discordgo.MessageDelete:
		return "messageDelete"
	case *discordgo.MessageDeleteBulk:
		return "messageDeleteBulk"
	case *discordgo.MessageReactionAdd:
		return "messageReactionAdd"
	case *discordgo.MessageReactionRemove:
		return "messageReactionRemove"
	case *discordgo.MessageReactionRemoveAll:
		return "messageReactionRemoveAll"
	case *discordgo.MessageUpdate:
		return "messageUpdate"
	case *discordgo.PresenceUpdate:
		return "presenceUpdate"
	case *discordgo.PresencesReplace:
		return "presencesReplace"
	case *discordgo.RateLimit:
		return "rateLimit"
	case *discordgo.Ready:
		return "ready"
	case *discordgo.RelationshipAdd:
		return "relationshipAdd"
	case *discordgo.RelationshipRemove:
		return "relationshipRemove"
	case *discordgo.Resumed:
		return "resumed"
	case *discordgo.ThreadCreate:
		return "threadCreate"
	case *discordgo.ThreadDelete:
		return "threadDelete"
	case *discordgo.ThreadListSync:
		return "threadListSync"
	case *discordgo.ThreadMemberUpdate:
		return "threadMemberUpdate"
	case *discordgo.ThreadMembersUpdate:
		return "threadMembersUpdate"
	case *discordgo.ThreadUpdate:
		return "threadUpdate"
	case *discordgo.TypingStart:
		return "typingStart"
	case *discordgo.UserGuildSettingsUpdate:
		return "userGuildSettingsUpdate"
	case *discordgo.UserNoteUpdate:
		return "userNoteUpdate"
	case *discordgo.UserSettingsUpdate:
		return "userSettingsUpdate"
	case *discordgo.UserUpdate:
		return "userUpdate"
	case *discordgo.VoiceServerUpdate:
		return "voiceServerUpdate"
	case *discordgo.VoiceStateUpdate:
		return "voiceStateUpdate"
	case *discordgo.WebhooksUpdate:
		return "webhooksUpdate"
	}
	return "unknown"
}
