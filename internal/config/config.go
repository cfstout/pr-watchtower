package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	GitHub struct {
		RefreshInterval time.Duration `mapstructure:"refresh_interval"`
		Queries         struct {
			NeedsReview string `mapstructure:"needs_review"`
			MyPRs       string `mapstructure:"my_prs"`
		} `mapstructure:"queries"`
	} `mapstructure:"github"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.pr-watchtower")

	// Defaults
	viper.SetDefault("github.refresh_interval", 2*time.Minute)
	viper.SetDefault("github.queries.needs_review", "review-requested:@me state:open")
	viper.SetDefault("github.queries.my_prs", "author:@me state:open")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		// Config file not found; ignore error if desired or warn
		// For now, we'll just use defaults if not found, but maybe we want to create one?
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return &cfg, nil
}
