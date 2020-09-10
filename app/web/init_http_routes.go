package web

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spyzhov/healthy/handler"
)

// setHttpRoutes declare all necessary HTTP methods
func (app *Application) setHttpRoutes() error {
	// region HTTP
	app.Http.Handle("/", handler.LoggedHandlerFunc(app.httpIndex))
	app.Http.Handle("/validate", handler.LoggedHandlerFunc(app.httpValidate))
	app.Http.Handle("/favicon.ico", handler.LoggedHandlerFunc(app.httpFavicon))
	//endregion
	// region Management
	app.Management.Handle("/health", handler.LoggedHandlerFunc(app.httpHealthCheck))
	app.Management.Handle("/info", handler.LoggedHandlerFunc(app.httpInfo))
	app.Management.Handle("/metrics", promhttp.Handler())
	//endregion
	return nil
}
