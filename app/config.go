package app

import (
	validator "github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	BotSecret   string                 `mapstructure:"bot_secret" envvar:"BOT_SECRET" validator:"required"`
	DatabaseDSN string                 `mapstructure:"database_dsn" envvar:"DATABASE_URL"`
	Cogs        map[string]interface{} `mapstructure:"cogs"`
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

	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
