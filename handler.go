package web

type (
	Handler      func(ctx *Context) (response, error)
	ErrorHandler func(ctx *Context, err error) response
)
