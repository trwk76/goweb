package web

import "net/http"

var (
	ErrNotFound       Error = Error{Status: http.StatusNotFound}
	ErrMethNotAllowed Error = Error{Status: http.StatusMethodNotAllowed}
)

type (
	Error struct {
		Status int
	}
)

func (e Error) Error() string {
	return http.StatusText(e.Status)
}

func (e Error) Write(w http.ResponseWriter) {
	NewDefaultResponse(e.Status).Write(w)
}
