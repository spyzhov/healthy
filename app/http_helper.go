package app

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

// function write writes any interface as json http/response
func write(w http.ResponseWriter, v interface{}) {
	if err := json.NewEncoder(w).Encode(v); err != nil {
		zap.L().Warn("error on write response", zap.Error(err))
	}
}
