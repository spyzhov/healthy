package web

import (
	"github.com/spyzhov/safe"
)

// init all necessary resources
func (app *Application) init() (err error) {
	if err = app.setTemplates(); err != nil {
		return safe.Wrap(err, "cannot initialize http/templates")
	}
	if err = app.setHttpRoutes(); err != nil {
		return safe.Wrap(err, "cannot initialize http/rotes")
	}
	return nil
}
