package app

import (
	"context"

	"github.com/fiffu/arisa3/app/log"
	"github.com/fiffu/arisa3/lib/envconfig"
	validator "github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	BotSecret   string                 `mapstructure:"bot_secret" envvar:"BOT_SECRET" validator:"required"`
	DatabaseDSN string                 `mapstructure:"database_dsn" envvar:"DATABASE_URL"`
	EnableDebug bool                   `mapstructure:"enable_debug" envvar:"ENABLE_DEBUG"`
	Cogs        map[string]interface{} `mapstructure:"cogs"`
}

func Configure(path string) (*Config, error) {
	ctx := context.Background()

	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	cfg := &Config{}
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	// Merge config from env vars
	if replaced, err := envconfig.MergeEnvVars(cfg, ""); err != nil {
		return nil, err
	} else if len(replaced) > 0 {
		for envKey, fld := range replaced {
			log.Infof(ctx,
				"Replaced %v with environment var %s",
				fld.Name,
				envKey,
			)
		}
	}

	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
