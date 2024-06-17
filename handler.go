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
		res := newRespBuffer()

		if contentType != "" {
			res.SetContentType(contentType)
		}

		if etag != "" {
			res.SetETag(etag)
		}

		http.ServeContent(&res, ctx.req, name, modTime, content)

		return res.response(), nil
	}
}
