package app

import (
	"context"
	"html/template"
	"net/http"
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

	Http       *http.ServeMux
	Management *http.ServeMux
	StepConfig *config.Config
	StepGroups *step.Groups
	Executor   *executor.Executor

	templates map[string]*template.Template
}

func New() (app *Application, err error) {
	app = &Application{
		error:      make(chan error, 1),
		Http:       http.NewServeMux(),
		Management: http.NewServeMux(),
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

	app.Config, err = NewConfig()
	if err != nil {
		return nil, err
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
	app.Logger.Debug("Application stops")

	defer close(app.error)
	defer safe.Close(app.Executor, "Executor")
	//defer safe.Close(app.Resource, "resource name")
}

// Start initialize all long-living processes
func (app *Application) Start() {
	defer app.Stop()

	// Run HTTP handler
	if err := app.RunHttp(app.Http, app.Config.Port, "HTTP Server"); err != nil {
		app.Logger.Panic("HTTP Server start error", zap.Error(err))
	}

	// Run HTTP Management handler
	if err := app.RunHttp(app.Management, app.Config.ManagementPort, "HTTP Management Server"); err != nil {
		app.Logger.Panic("HTTP Management Server start error", zap.Error(err))
	}

	select {
	case err := <-app.error:
		app.Logger.Panic("service crashed", zap.Error(err))
	case <-app.Ctx.Done():
		app.Logger.Info("service stops via context")
	case sig := <-WaitExit():
		app.Logger.Info("service stop", zap.Stringer("signal", sig))
	}
}

// Stop waits for all resources be cosed
func (app *Application) Stop() {
	app.Logger.Info("service stopping...")
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
