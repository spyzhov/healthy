package app

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

func (app *Application) httpPublic(w http.ResponseWriter, r *http.Request) {
	if file, ok := app.files[r.RequestURI]; ok {
		w.Header().Set("Content-Type", file.ContentType)
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(file.Content); err != nil {
			app.Logger.Warn("error on write response", zap.Error(err))
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		if _, err := fmt.Fprint(w, "Not found"); err != nil {
			app.Logger.Warn("error on write response", zap.Error(err))
		}
	}
}
