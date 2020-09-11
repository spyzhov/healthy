package app

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

func (app *Application) httpHealthCheck(w http.ResponseWriter, _ *http.Request) {
	info, status := app.healthCheck()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(info); err != nil {
		app.Logger.Warn("error on write response", zap.Error(err))
	}
}
