package app

import (
	"github.com/spf13/viper"
)

type Config struct {
	BotSecret string                 `mapstructure:"botSecret"`
	Cogs      map[string]interface{} `mapstructure:"cogs"`

	GuildID        string // flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	RemoveCommands bool   // flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
}

func Configure(path string) (*Config, error) {
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	cfg := &Config{}
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
