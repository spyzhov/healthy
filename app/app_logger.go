package app

import (
	"fmt"
	"strings"

	"github.com/spyzhov/safe"
	"go.uber.org/zap"
)

func NewLogger(level string, format string) (logger *zap.Logger, err error) {
	var cfg zap.Config
	switch strings.ToLower(format) {
	case "json":
		cfg = zap.NewProductionConfig()
	case "text":
		cfg = zap.NewDevelopmentConfig()
	default:
		return nil, fmt.Errorf("log-format should be one of the next: json, text")
	}
	cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)

	atom := zap.NewAtomicLevel()
	err = atom.UnmarshalText([]byte(level))
	if err != nil {
		return nil, err
	}

	cfg.Level = atom

	logger, err = cfg.Build()
	if err != nil {
		return nil, err
	}
	zap.ReplaceGlobals(logger)

	safe.Printf = func(format string, v ...interface{}) {
		logger.Warn(fmt.Sprintf(format, v...))
	}

	return logger, err
}
