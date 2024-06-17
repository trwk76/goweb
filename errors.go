package web

import "net/http"

var (
	ErrNotFound       Error = Error{Status: http.StatusNotFound}
	ErrMethNotAllowed Error = Error{Status: http.StatusMethodNotAllowed}
	ErrInternalError  Error = Error{Status: http.StatusInternalServerError}
)

type (
	Error struct {
		Status int
	}
)

func (e Error) Error() string {
	return http.StatusText(e.Status)
}

func (e Error) WriteResponse(w http.ResponseWriter) {
	NewDefaultResponse(e.Status).WriteResponse(w)
}
