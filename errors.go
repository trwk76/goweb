package web

import "net/http"

var (
	ErrNotFound       HTTPError = HTTPError{Status: http.StatusNotFound}
	ErrMethNotAllowed HTTPError = HTTPError{Status: http.StatusMethodNotAllowed}
	ErrInternalError  HTTPError = HTTPError{Status: http.StatusInternalServerError}
)

type (
	HTTPError struct {
		Status int
	}
)

func (e HTTPError) Error() string {
	return http.StatusText(e.Status)
}

func (e HTTPError) Response() Response {
	return NewDefaultResponse(e.Status)
}
