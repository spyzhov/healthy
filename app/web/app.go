package web

import (
	"html/template"
	"net/http"

	"github.com/spyzhov/healthy/app"
	"go.uber.org/zap"
)

type Application struct {
	*app.Application
	Config     *Config
	Http       *http.ServeMux
	Management *http.ServeMux
	templates  map[string]*template.Template
	favicon    []byte
}

func New() (web *Application, err error) {
	web = &Application{
		Http:       http.NewServeMux(),
		Management: http.NewServeMux(),
	}
	web.Application, err = app.New(nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			web.Close()
		}
	}()

	web.Config, err = NewConfig()
	if err != nil {
		return nil, err
	}

	return web, web.init()
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

	app.Application.Start()
}
