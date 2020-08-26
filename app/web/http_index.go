package web

import (
	"fmt"
	"net/http"

	"github.com/spyzhov/healthy/config"
	"github.com/spyzhov/healthy/step"
	"go.uber.org/zap"
)

type IndexContext struct {
	Name   string
	Front  config.Frontend
	Groups *step.Groups
}

func (app *Application) httpIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet || r.RequestURI != "/" {
		w.WriteHeader(http.StatusNotFound)
		if _, err := fmt.Fprint(w, "Not found"); err != nil {
			app.Logger.Warn("error on write response", zap.Error(err))
		}
		return
	}
	err := app.templates["index"].ExecuteTemplate(w, "index", &IndexContext{
		Name:   app.StepConfig.Name,
		Front:  app.StepConfig.Frontend,
		Groups: app.StepGroups,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := fmt.Fprintf(w, "error: %s", err.Error()); err != nil {
			app.Logger.Warn("error on write response", zap.Error(err))
		}
	}
}
