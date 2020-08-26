package cmd

import (
	"os"

	"github.com/spyzhov/healthy/app"
	"go.uber.org/zap"
)

type Application struct {
	*app.Application
	Config *Config
	Status int
}

func New() (cmd *Application, err error) {
	cmd = new(Application)
	cmd.Config, err = NewConfig()
	if err != nil {
		return nil, err
	}
	if cmd.Config.CallVersion {
		cmd.printVersion()
	}

	cmd.Application, err = app.New(&cmd.Config.Config)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			cmd.Close()
		}
	}()

	return cmd, cmd.init()
}

// Start initialize all long-living processes
func (app *Application) Start() {
	defer app.Stop()

	// Run Action
	if err := app.RunAction(); err != nil {
		app.Logger.Panic("Action start error", zap.Error(err))
	}

	app.Application.Start()
}

// Close all necessary resources
func (app *Application) Close() {
	defer os.Exit(app.Status)
	app.Application.Close()
}
