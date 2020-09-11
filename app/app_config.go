package app

import (
	"github.com/caarlos0/env/v6"
	"github.com/spf13/pflag"
)

type Config struct {
	LogLevel   string `env:"HEALTHY_LOG_LEVEL" envDefault:"warn"`
	LogFormat  string `env:"HEALTHY_LOG_FORMAT" envDefault:"json"`
	ConfigFile string `env:"HEALTHY_CONFIG_FILE"`
	ConfigYaml string `env:"HEALTHY_CONFIG_YAML"`

	Port           int `env:"HEALTHY_PORT" envDefault:"80"`
	ManagementPort int `env:"HEALTHY_MANAGEMENT_PORT" envDefault:"3280"`

	CallWeb     bool
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
	pflag.StringVarP(&cfg.LogLevel, "log-level", "l", cfg.LogLevel, "log level for current run")
	pflag.StringVar(&cfg.LogFormat, "log-format", cfg.LogFormat, "log format for current run")
	pflag.StringVarP(&cfg.ConfigFile, "config-file", "f", cfg.ConfigFile, "config file to run")
	pflag.StringVarP(&cfg.ConfigYaml, "config-yaml", "y", cfg.ConfigYaml, "config yaml to run")

	pflag.IntVarP(&cfg.Port, "port", "p", cfg.Port, "web: public port")
	pflag.IntVarP(&cfg.ManagementPort, "management-port", "m", cfg.ManagementPort, "web: management port")

	pflag.StringSliceVarP(&cfg.Groups, "group", "g", cfg.Groups, "cli: groups to run")
	pflag.StringSliceVarP(&cfg.Steps, "step", "s", cfg.Steps, "cli: steps to run")
	pflag.BoolVarP(&cfg.Quiet, "quiet", "q", cfg.Quiet, "cli: quiet output")
	pflag.BoolVarP(&cfg.Verbose, "verbose", "v", cfg.Verbose, "cli: verbose output")

	pflag.BoolVar(&cfg.CallVersion, "version", cfg.CallVersion, "print current version")
	pflag.BoolVarP(&cfg.CallWeb, "web", "w", cfg.CallWeb, "run as web")

	pflag.Parse()

	return cfg, nil
}
