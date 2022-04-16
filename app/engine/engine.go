package engine

import (
	"encoding/json"
	"errors"
	"fmt"

	dgo "github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

var (
	ErrCogParseConfig        = errors.New("unable to parse cog config")
	ErrUnexpectedConfigValue = errors.New("config type assert failed")
)

func ParseConfig(in interface{}, out interface{}) error {
	if out == nil {
		return nil
	}
	log.Warn().Msgf("ParseConfig-tgt: %#v, nil? %v", out, out == nil)
	log.Warn().Msgf("ParseConfig-in:  %#v", in)
	bytes, err := json.Marshal(in)
	log.Warn().Msgf("ParseConfig-byt:  %#v", string(bytes))
	if err != nil {
		return fmt.Errorf("%w: %s", ErrCogParseConfig, err)
	}
	err = json.Unmarshal(bytes, out)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrCogParseConfig, err)
	}
	log.Warn().Msgf("ParseConfig-out: %#v, nil? %v", out, out == nil)
	return nil
}

func UnexpectedConfigType(wanted interface{}, got interface{}) error {
	return fmt.Errorf("%w, wanted: %T, got: %#v", ErrUnexpectedConfigValue, wanted, got)
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
