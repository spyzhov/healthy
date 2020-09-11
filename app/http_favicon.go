package app

import (
	"net/http"

	"go.uber.org/zap"
)

func (app *Application) httpFavicon(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("content-type", "image/x-icon")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(app.favicon); err != nil {
		app.Logger.Warn("error on write response", zap.Error(err))
	}
}
