package cmd

import (
	"time"

	"go.uber.org/zap"
)

// Start of action
func (app *Application) RunAction() error {
	app.WaitGroup.Add(1)
	go func() {
		defer app.WaitGroup.Done()
		defer func(start time.Time) {
			app.Logger.Debug("Action: done", zap.Duration("duration", time.Since(start)))
		}(time.Now())
		zap.L().Debug("Action: start")
		app.Status = app.execute(app.Config.Groups, app.Config.Steps)
		app.CtxCancel()
	}()
	return nil
}
