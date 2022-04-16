package types

import (
	dgo "github.com/bwmarrin/discordgo"
)

type IEmbed interface {
	Data() *dgo.MessageEmbed
}
