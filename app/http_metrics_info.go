package app

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

func (app *Application) httpInfo(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(app.Info); err != nil {
		app.Logger.Warn("error on write response", zap.Error(err))
	}
}
