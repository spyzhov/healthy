package cmd

import (
	"github.com/caarlos0/env/v6"
	"github.com/spf13/pflag"
	"github.com/spyzhov/healthy/app"
)

type Config struct {
	app.Config
	CallVersion bool
	Quiet       bool
	Verbose     bool
	Groups      []string
	Steps       []string
}

func NewConfig() (cfg *Config, err error) {
	cfg = new(Config)
	if err = env.Parse(cfg); err != nil {
		return nil, err
	}
	return FlagsConfig(cfg)
}

func FlagsConfig(cfg *Config) (*Config, error) {
	pflag.StringVarP(&cfg.Level, "log-level", "l", cfg.Level, "log level for current run")
	pflag.StringVarP(&cfg.ConfigFile, "config-file", "f", cfg.ConfigFile, "config file to run")
	pflag.StringVarP(&cfg.ConfigYaml, "config-yaml", "y", cfg.ConfigYaml, "config yaml to run")

	pflag.StringSliceVarP(&cfg.Groups, "group", "g", cfg.Groups, "groups to run")
	pflag.StringSliceVarP(&cfg.Steps, "step", "s", cfg.Steps, "steps to run")
	pflag.BoolVarP(&cfg.Quiet, "quiet", "q", cfg.Quiet, "quiet output")
	pflag.BoolVarP(&cfg.Verbose, "verbose", "v", cfg.Verbose, "verbose output")
	pflag.BoolVar(&cfg.CallVersion, "version", cfg.CallVersion, "print current version")

	pflag.Parse()

	return cfg, nil
}
