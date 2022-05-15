package commandfilters

import (
	"github.com/bwmarrin/discordgo"
	"github.com/fiffu/arisa3/app/types"
)

func getMember(ev types.ICommandEvent) *discordgo.Member {
	return ev.Interaction().Member
}

func IsFromGuild(ev types.ICommandEvent) bool {
	return getMember(ev) != nil
}

func IsGuildAdmin(ev types.ICommandEvent) bool {
	if !IsFromGuild(ev) {
		return false
	}
	adminPerms := getMember(ev).Permissions & discordgo.PermissionAdministrator
	return adminPerms > 0
}
