package app

import (
	"github.com/caarlos0/env/v6"
)

type Config struct {
	Level          string `env:"LOG_LEVEL" envDefault:"info"`
	Port           int    `env:"PORT" envDefault:"80"`
	ManagementPort int    `env:"MANAGEMENT_PORT" envDefault:"3280"`
	ConfigFile     string `env:"CONFIG_FILE"`
	ConfigYaml     string `env:"CONFIG_YAML"`
}

func NewConfig() (cfg *Config, err error) {
	cfg = new(Config)
	return cfg, env.Parse(cfg)
}
