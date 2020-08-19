package transport

import "net/http"

// RoundTripper is the main interface, implements the http.RoundTripper interface with additional WithNext function
type RoundTripper interface {
	http.RoundTripper
	WithNext(next RoundTripper) RoundTripper
}
