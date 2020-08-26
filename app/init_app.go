package app

import (
	"github.com/spyzhov/safe"
)

// init all necessary resources
func (app *Application) init() (err error) {
	if err = app.setExecutor(); err != nil {
		return safe.Wrap(err, "cannot initialize Executor")
	}
	if err = app.setConfig(); err != nil {
		return safe.Wrap(err, "cannot initialize Config for Steps")
	}
	if err = app.setSteps(); err != nil {
		return safe.Wrap(err, "cannot initialize Steps")
	}
	return nil
}
