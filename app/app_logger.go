package app

import (
	"fmt"

	"github.com/spyzhov/safe"
	"go.uber.org/zap"
)

func NewLogger(level string) (logger *zap.Logger, err error) {
	cfg := zap.NewProductionConfig()
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
