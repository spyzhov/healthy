package handler

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/spyzhov/safe"
	"go.uber.org/zap"
)

// The LoggedHandlerFunc type is an adapter to allow the use of
// ordinary functions as HTTP handlers. If f is a function
// with the appropriate signature, LoggedHandlerFunc(f) is a
// Handler that calls f.
// Addition:
//   1. Adds access log for each response;
//   2. Recover panics and log them;
//   3. Close request.Body;
type LoggedHandlerFunc func(http.ResponseWriter, *http.Request)

// ServeHTTP calls f(w, r).
func (f LoggedHandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res := &response{w, http.StatusOK, 0}
	defer func(start time.Time) {
		if rec := recover(); rec != nil {
			zap.L().Error("recover panic", zap.Any("recover", rec), zap.ByteString("stack", debug.Stack()))
			if res.ContentLength() == 0 {
				res.WriteHeader(http.StatusInternalServerError)
			}
		}
		zap.L().Info(
			"request processed",
			zap.String("proto", r.Proto),
			zap.String("method", r.Method),
			zap.String("host", r.Host),
			zap.String("request", r.RequestURI),
			zap.Int("status", res.Status()),
			zap.Int("content_length", res.ContentLength()),
			zap.String("clientip", r.RemoteAddr),
			zap.String("agent", r.UserAgent()),
			zap.String("referrer", r.Referer()),
			zap.Duration("request_time", time.Since(start)),
		)
	}(time.Now())
	defer safe.Close(r.Body, "request body")
	f(res, r)
}
