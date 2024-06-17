package web

import (
	"io"
	"net/http"
	"time"
)

type (
	Handler      func(ctx *Context) (Response, error)
	ErrorHandler func(ctx *Context, err error) Response
)

func Content(name string, contentType string, etag string, modTime time.Time, content io.ReadSeeker) Handler {
	return func(ctx *Context) (Response, error) {
		res := NewRespBuffer()

		SetContentType(res.hdr, contentType)
		SetETag(res.hdr, etag)

		http.ServeContent(&res, ctx.req, name, modTime, content)

		return &res, nil
	}
}
