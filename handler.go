package web

import (
	"io"
	"net/http"
	"time"
)

type (
	Handler      func(ctx *Context) (response, error)
	ErrorHandler func(ctx *Context, err error) response
)

func Content(name string, contentType string, etag string, modTime time.Time, content io.ReadSeeker) Handler {
	return func(ctx *Context) (response, error) {
		res := NewRespBuffer()

		SetHeaderContentType(res.hdr, contentType)
		SetHeaderETag(res.hdr, etag)

		http.ServeContent(&res, ctx.req, name, modTime, content)

		return &res, nil
	}
}
