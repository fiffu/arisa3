package types

import (
	"context"

	dgo "github.com/bwmarrin/discordgo"
)

type IApp interface {
	Configs() map[string]interface{}
	Debugf(context.Context, string, ...interface{})
	Infof(context.Context, string, ...interface{})
	Warnf(context.Context, string, ...interface{})
	Errorf(context.Context, error, string, ...interface{})
	ContextWithValue(ctx context.Context, key, value string) context.Context
}

type CogConfig interface{}

type ICog interface {
	Name() string
	OnStartup(ctx context.Context, sess *dgo.Session, config CogConfig) error
}
