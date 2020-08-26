package web

import (
	"github.com/caarlos0/env/v6"
)

type Config struct {
	Port           int `env:"PORT" envDefault:"80"`
	ManagementPort int `env:"MANAGEMENT_PORT" envDefault:"3280"`
}

func NewConfig() (cfg *Config, err error) {
	cfg = new(Config)
	return cfg, env.Parse(cfg)
}
