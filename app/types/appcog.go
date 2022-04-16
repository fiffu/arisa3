package types

import (
	"context"

	dgo "github.com/bwmarrin/discordgo"
)

type IApp interface {
	Configs() map[string]interface{}
}

type CogConfig interface{}
type StructPointer interface{}

type ICog interface {
	Name() string
	ConfigPointer() StructPointer
	OnStartup(ctx context.Context, sess *dgo.Session, config CogConfig) error
}
