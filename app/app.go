package app

import (
	"context"
	"html/template"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/spyzhov/healthy/config"
	"github.com/spyzhov/healthy/executor"
	"github.com/spyzhov/healthy/step"
	"github.com/spyzhov/safe"
	"go.uber.org/zap"
)

type Application struct {
	// region System
	Config    *Config
	Logger    *zap.Logger
	Info      *BuildInfo
	Ctx       context.Context
	CtxCancel context.CancelFunc
	WaitGroup *sync.WaitGroup

	error chan error
	once  sync.Once
	// endregion
	// region Service
	StepConfig *config.Config
	StepGroups *step.Groups
	Executor   *executor.Executor
	// endregion
	// region Web
	Http       *http.ServeMux
	Management *http.ServeMux
	templates  map[string]*template.Template
	favicon    []byte
	// endregion
	// region Cli
	Status int
	// endregion
}

func New() (app *Application, err error) {
	app = &Application{
		// region System
		error:      make(chan error, 1),
		Info:       NewBuildInfo(),
		WaitGroup:  new(sync.WaitGroup),
		StepGroups: step.NewGroups(),
		// endregion
		// region Web
		Http:       http.NewServeMux(),
		Management: http.NewServeMux(),
		// endregion
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
	if app.Config.CallVersion {
		app.printVersion()
	}
	app.Logger, err = NewLogger(app.Config.LogLevel, app.Config.LogFormat)
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

	defer os.Exit(app.Status)
	defer close(app.error)
	defer safe.Close(app.Executor, "Executor")
	//defer safe.Close(app.Resource, "resource name")
}

// Start initialize all long-living processes
func (app *Application) Start() {
	defer app.Stop()

	if app.Config.CallWeb {
		// Run HTTP handler
		if err := app.RunHttp(app.Http, app.Config.Port, "HTTP Server"); err != nil {
			app.Logger.Panic("HTTP Server start error", zap.Error(err))
		}

		// Run HTTP Management handler
		if err := app.RunHttp(app.Management, app.Config.ManagementPort, "HTTP Management Server"); err != nil {
			app.Logger.Panic("HTTP Management Server start error", zap.Error(err))
		}
	} else {
		// Run Action
		if err := app.RunAction(); err != nil {
			app.Logger.Panic("Action start error", zap.Error(err))
		}
	}

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
