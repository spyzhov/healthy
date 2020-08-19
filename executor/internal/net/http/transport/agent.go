package transport

import (
	"net/http"
)

// AgentRoundTripper implements the http.RoundTripper interface.
type AgentRoundTripper struct {
	Next    http.RoundTripper
	headers map[string]string
}

// NewAgent is a constructor for AgentRoundTripper, that provides default transport.
func NewAgent(name string) *AgentRoundTripper {
	return &AgentRoundTripper{
		Next: http.DefaultTransport,
		headers: map[string]string{
			"Agent":      name,
			"User-Agent": name,
		},
	}
}

// WithNext can set next RoundTripper object
func (s *AgentRoundTripper) WithNext(next RoundTripper) RoundTripper {
	s.Next = next
	return s
}

// RoundTrip executes a single HTTP transaction, returning
// a Response for the provided Request.
func (s *AgentRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	for name, value := range s.headers {
		if request.Header.Get(name) == "" {
			request.Header.Set(name, value)
		}
	}
	return s.Next.RoundTrip(request)
}
