package http

import (
	"net/http"
	"time"

	"github.com/spyzhov/healthy/executor/internal/net/http/transport"
)

func GetClient(timeout time.Duration, version string) *http.Client {
	return &http.Client{
		Timeout:   timeout,
		Transport: transport.NewAgent("healthy/" + version),
	}
}
