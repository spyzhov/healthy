package handler

import "net/http"

// A ResponseWriter interface is used by an HTTP handler to
// construct an HTTP response.
//
// A ResponseWriter may not be used after the Handler.ServeHTTP method
// has returned.
type ResponseWriter interface {
	http.ResponseWriter
	// Status returns status code of response
	Status() int
	// Status returns status code of response
	ContentLength() int
}

type response struct {
	http.ResponseWriter
	status int
	length int
}

func (r *response) WriteHeader(statusCode int) {
	r.status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *response) Write(data []byte) (int, error) {
	n, err := r.ResponseWriter.Write(data)
	r.length += n
	return n, err
}

func (r *response) Status() int {
	return r.status
}

func (r *response) ContentLength() int {
	return r.length
}
