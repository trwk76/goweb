package web

import "net/http"

var (
	ErrNotFound       Error = Error{StatusCode: http.StatusNotFound}
	ErrMethNotAllowed Error = Error{StatusCode: http.StatusMethodNotAllowed}
	ErrInternalError  Error = Error{StatusCode: http.StatusInternalServerError}
)

type (
	Error struct {
		StatusCode int
	}
)

func (e Error) Error() string {
	return http.StatusText(e.StatusCode)
}

func (e Error) Status() int {
	return e.StatusCode
}

func (e Error) Response() Response {
	return NewDefaultResponse(e.StatusCode)
}
