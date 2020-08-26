package app

import "github.com/spyzhov/healthy/executor"

func (app *Application) setExecutor() error {
	app.Logger.Debug("initialize Executor")
	app.Executor = executor.NewExecutor(app.Ctx, app.Info.Version)
	return nil
}
