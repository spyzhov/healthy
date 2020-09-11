package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"

	. "github.com/spyzhov/healthy/step"
	"github.com/spyzhov/safe"
	"go.uber.org/zap"
)

type request struct {
	Group string `json:"group"`
	Name  string `json:"name"`
}

type response struct {
	Message string `json:"message"`
	Level   Status `json:"level"`
}

// Method httpValidate will be linked to /validate route and run any tests from frontend
func (app *Application) httpValidate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer func() {
		if rec := recover(); rec != nil {
			zap.L().Error("recover panic", zap.Any("recover", rec), zap.ByteString("stack", debug.Stack()))
			w.WriteHeader(http.StatusInternalServerError)
			write(w, &response{
				Message: fmt.Sprintf("panic: %v", rec),
				Level:   Error,
			})
		}
	}()
	var test request
	err := json.NewDecoder(r.Body).Decode(&test)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		write(w, &response{
			Message: fmt.Sprintf("error: %s", err.Error()),
			Level:   Error,
		})
		return
	}
	step := app.StepGroups.Get(test.Group).Get(test.Name)
	if safe.IsNil(step) {
		w.WriteHeader(http.StatusNotFound)
		write(w, &response{
			Message: fmt.Sprintf("Not found: %s / %s", test.Group, test.Name),
			Level:   Error,
		})
		return
	}

	res := step.Call()
	if res != nil {
		write(w, &response{
			Message: res.Message,
			Level:   res.Status,
		})
	} else {
		write(w, &response{
			Message: "Validation result is NIL",
			Level:   Warning,
		})
	}
}
