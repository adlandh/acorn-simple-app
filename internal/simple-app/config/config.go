// Package config contains application configuration
package config

import (
	"fmt"

	"github.com/caarlos0/env/v10"
)

type RedisConfig struct {
	URL    string `env:"URL,notEmpty"`
	Prefix string `env:"PREFIX" envDefault:"simple-app"`
}

type Config struct {
	Port  string      `env:"PORT" envDefault:"8080"`
	Redis RedisConfig `envPrefix:"REDIS_"`
}

func NewConfig() (*Config, error) {
	var cfg Config

	if err := env.ParseWithOptions(&cfg, env.Options{
		RequiredIfNoDef: false,
	}); err != nil {
		return nil, fmt.Errorf("error parsing config: %w", err)
	}

	return &cfg, nil
}
