// package types contains custom types and their factories.

package types

//go:generate mockgen -source=appcog.go -destination=./appcog_mock.go -package=types

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/fiffu/arisa3/app/database"
)

type IApp interface {
	Configs() map[string]interface{}
	Database() database.IDatabase
	BotSession() *discordgo.Session
	Shutdown(context.Context)
}

type CogConfig interface{}
type StructPointer interface{}

type ICog interface {
	Name() string
	ConfigPointer() StructPointer
	OnStartup(ctx context.Context, app IApp, config CogConfig) error
}

// IShutdown is optional for a cog. The app calls OnShutdown after it closes the gateway session
// and before it closes the database. OnShutdown must return when ctx ends.
type IShutdown interface {
	OnShutdown(ctx context.Context) error
}
