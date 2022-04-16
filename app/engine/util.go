package engine

import (
	"arisa3/app/types"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	dgo "github.com/bwmarrin/discordgo"
)

var (
	ErrCogParseConfig = errors.New("unable to parse cog config")
)

const (
	ctxFieldCog     = "cog"
	ctxFieldCommand = "command"
	ctxFieldUser    = "user"
	ctxFieldEvent   = "event"
)

type IDefaultStartup interface {
	ConfigType() interface{}
	Configure(ctx context.Context, cfg interface{}) error
	ReadyCallback(s *dgo.Session, r *dgo.Ready)
}

func ParseConfig(in interface{}, out interface{}) error {
	bytes, err := json.Marshal(in)
	if err != nil {
		return ErrCogParseConfig
	}
	err = json.Unmarshal(bytes, out)
	if err != nil {
		return ErrCogParseConfig
	}
	return nil
}

// DefaultStartup parses config and pushes it to cog, and sets up a handler for discordgo.Ready event.
func DefaultStartup(ctx context.Context, sess *dgo.Session, rawConfig types.CogConfig, cog IDefaultStartup) error {
	cfg := cog.ConfigType()
	if err := ParseConfig(rawConfig, &cfg); err != nil {
		return fmt.Errorf("unable to parse cog config: %w", err)
	}
	cog.Configure(ctx, cfg)

	sess.AddHandler(func(s *dgo.Session, r *dgo.Ready) {
		cog.ReadyCallback(s, r)
	})
	return nil
}

// ParseEvent returns the event name based on event handler structs provided by discordgo.
// Courtesy of regex match on discordgo@v0.24.0/eventhandlers.go
func ParseEvent(event interface{}) string {
	switch event.(type) {
	case *dgo.ChannelCreate:
		return "channelCreate"
	case *dgo.ChannelDelete:
		return "channelDelete"
	case *dgo.ChannelPinsUpdate:
		return "channelPinsUpdate"
	case *dgo.ChannelUpdate:
		return "channelUpdate"
	case *dgo.Connect:
		return "connect"
	case *dgo.Disconnect:
		return "disconnect"
	case *dgo.Event:
		return "event"
	case *dgo.GuildBanAdd:
		return "guildBanAdd"
	case *dgo.GuildBanRemove:
		return "guildBanRemove"
	case *dgo.GuildCreate:
		return "guildCreate"
	case *dgo.GuildDelete:
		return "guildDelete"
	case *dgo.GuildEmojisUpdate:
		return "guildEmojisUpdate"
	case *dgo.GuildIntegrationsUpdate:
		return "guildIntegrationsUpdate"
	case *dgo.GuildScheduledEventCreate:
		return "guildScheduledEventCreate"
	case *dgo.GuildScheduledEventUpdate:
		return "guildScheduledEventUpdate"
	case *dgo.GuildScheduledEventDelete:
		return "guildScheduledEventDelete"
	case *dgo.GuildMemberAdd:
		return "guildMemberAdd"
	case *dgo.GuildMemberRemove:
		return "guildMemberRemove"
	case *dgo.GuildMemberUpdate:
		return "guildMemberUpdate"
	case *dgo.GuildMembersChunk:
		return "guildMembersChunk"
	case *dgo.GuildRoleCreate:
		return "guildRoleCreate"
	case *dgo.GuildRoleDelete:
		return "guildRoleDelete"
	case *dgo.GuildRoleUpdate:
		return "guildRoleUpdate"
	case *dgo.GuildUpdate:
		return "guildUpdate"
	case *dgo.InteractionCreate:
		return "interactionCreate"
	case *dgo.InviteCreate:
		return "inviteCreate"
	case *dgo.InviteDelete:
		return "inviteDelete"
	case *dgo.MessageAck:
		return "messageAck"
	case *dgo.MessageCreate:
		return "messageCreate"
	case *dgo.MessageDelete:
		return "messageDelete"
	case *dgo.MessageDeleteBulk:
		return "messageDeleteBulk"
	case *dgo.MessageReactionAdd:
		return "messageReactionAdd"
	case *dgo.MessageReactionRemove:
		return "messageReactionRemove"
	case *dgo.MessageReactionRemoveAll:
		return "messageReactionRemoveAll"
	case *dgo.MessageUpdate:
		return "messageUpdate"
	case *dgo.PresenceUpdate:
		return "presenceUpdate"
	case *dgo.PresencesReplace:
		return "presencesReplace"
	case *dgo.RateLimit:
		return "rateLimit"
	case *dgo.Ready:
		return "ready"
	case *dgo.RelationshipAdd:
		return "relationshipAdd"
	case *dgo.RelationshipRemove:
		return "relationshipRemove"
	case *dgo.Resumed:
		return "resumed"
	case *dgo.ThreadCreate:
		return "threadCreate"
	case *dgo.ThreadDelete:
		return "threadDelete"
	case *dgo.ThreadListSync:
		return "threadListSync"
	case *dgo.ThreadMemberUpdate:
		return "threadMemberUpdate"
	case *dgo.ThreadMembersUpdate:
		return "threadMembersUpdate"
	case *dgo.ThreadUpdate:
		return "threadUpdate"
	case *dgo.TypingStart:
		return "typingStart"
	case *dgo.UserGuildSettingsUpdate:
		return "userGuildSettingsUpdate"
	case *dgo.UserNoteUpdate:
		return "userNoteUpdate"
	case *dgo.UserSettingsUpdate:
		return "userSettingsUpdate"
	case *dgo.UserUpdate:
		return "userUpdate"
	case *dgo.VoiceServerUpdate:
		return "voiceServerUpdate"
	case *dgo.VoiceStateUpdate:
		return "voiceStateUpdate"
	case *dgo.WebhooksUpdate:
		return "webhooksUpdate"
	}
	return "unknown"
}
