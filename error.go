package web

import "net/http"

type (
	Error struct {
		Status int
		Nested error
	}

	ErrorHandler func(ctx *Context, err error) Response
)

func (e Error) Error() string {
	if e.Nested != nil {
		return e.Nested.Error()
	}

	return http.StatusText(e.Status)
}

func EnsureError(err error) Error {
	if e, ok := err.(Error); ok {
		return e
	}

	return Error{
		Status: http.StatusInternalServerError,
		Nested: err,
	}
}
