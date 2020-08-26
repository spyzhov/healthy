package app

import (
	"context"
	"sync"
	"time"

	"github.com/spyzhov/healthy/config"
	"github.com/spyzhov/healthy/executor"
	"github.com/spyzhov/healthy/step"
	"github.com/spyzhov/safe"
	"go.uber.org/zap"
)

type Application struct {
	Config    *Config
	Logger    *zap.Logger
	Info      *BuildInfo
	Ctx       context.Context
	CtxCancel context.CancelFunc
	WaitGroup *sync.WaitGroup

	error chan error
	once  sync.Once

	StepConfig *config.Config
	StepGroups *step.Groups
	Executor   *executor.Executor
}

func New(config *Config) (app *Application, err error) {
	app = &Application{
		error:      make(chan error, 1),
		Info:       NewBuildInfo(),
		WaitGroup:  new(sync.WaitGroup),
		StepGroups: step.NewGroups(),
	}

	app.Ctx, app.CtxCancel = context.WithCancel(context.Background())
	defer func() {
		if err != nil {
			app.Close()
		}
	}()
	if config == nil {
		app.Config, err = NewConfig()
		if err != nil {
			return nil, err
		}
	} else {
		app.Config = config
	}
	app.Logger, err = NewLogger(app.Config.Level)
	if err != nil {
		return nil, err
	}
	app.Logger.Debug("debug mode on")

	return app, app.init()
}

// Close all necessary resources
func (app *Application) Close() {
	zap.L().Debug("Application stops")
	if app == nil {
		return
	}

	defer close(app.error)
	defer safe.Close(app.Executor, "Executor")
	//defer safe.Close(app.Resource, "resource name")
}

// Start initialize all long-living processes
func (app *Application) Start() {
	select {
	case err := <-app.error:
		app.Logger.Panic("crashed", zap.Error(err))
	case <-app.Ctx.Done():
		app.Logger.Info("stops via context")
	case sig := <-WaitExit():
		app.Logger.Info("stop", zap.Stringer("signal", sig))
	}
}

// Stop waits for all resources be cosed
func (app *Application) Stop() {
	if app == nil {
		return
	}
	app.Logger.Info("stopping...")
	app.CtxCancel()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	go func() {
		defer cancel()
		app.WaitGroup.Wait()
	}()

	<-ctx.Done()

	if ctx.Err() != context.Canceled {
		app.Logger.Panic("service stopped with timeout")
	} else {
		app.Logger.Info("service stopped with success")
	}
}

// Error - register global error, but only once
func (app *Application) Error(err error) {
	if !safe.IsNil(err) {
		app.once.Do(func() {
			app.error <- err
		})
	}
}
