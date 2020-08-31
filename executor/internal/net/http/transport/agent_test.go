package transport

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAgentRoundTripper_RoundTrip(t *testing.T) {
	names := []string{"Agent", "User-Agent"}
	value := "service-name"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, name := range names {
			if r.Header.Get(name) != value {
				t.Error("Header corrupted")
			}
		}
	}))
	defer server.Close()
	var (
		c = &http.Client{
			Transport: NewAgent(value),
		}
		request, err = http.NewRequest("GET", server.URL, nil)
	)
	if err != nil {
		t.Errorf("Can't create new request: %s", err)
	} else if _, err = c.Do(request); err != nil {
		t.Errorf("Can't do the request: %s", err)
	}
}
