package web

import (
	"net/http"
	"time"
)

// Handle function for health-check
func (app *Application) healthCheck() (info map[string]string, status int) {
	status = http.StatusOK
	info = map[string]string{
		"service": "healthy",
		"time":    time.Now().String(),
	}
	return info, status
}
