// package types contains custom types and their factories.

package types

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/fiffu/arisa3/app/database"
)

type IApp interface {
	Configs() map[string]interface{}
	Database() database.IDatabase
	BotSession() *discordgo.Session
	Shutdown()
}

type CogConfig interface{}
type StructPointer interface{}

type ICog interface {
	Name() string
	ConfigPointer() StructPointer
	OnStartup(ctx context.Context, app IApp, config CogConfig) error
}
