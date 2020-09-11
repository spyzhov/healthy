package app

import (
	"github.com/spyzhov/safe"
)

// init all necessary resources
func (app *Application) init() (err error) {
	// region Service
	if err = app.setExecutor(); err != nil {
		return safe.Wrap(err, "cannot initialize Executor")
	}
	if err = app.setConfig(); err != nil {
		return safe.Wrap(err, "cannot initialize Config for Steps")
	}
	if err = app.setSteps(); err != nil {
		return safe.Wrap(err, "cannot initialize Steps")
	}
	// endregion
	// region Web
	if err = app.setTemplates(); err != nil {
		return safe.Wrap(err, "cannot initialize http/templates")
	}
	if err = app.setFavicon(); err != nil {
		return safe.Wrap(err, "cannot initialize favicon.ico")
	}
	if err = app.setHttpRoutes(); err != nil {
		return safe.Wrap(err, "cannot initialize http/rotes")
	}
	// endregion
	// region Cli
	// blank
	// endregion
	return nil
}
