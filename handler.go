package web

import (
	"io"
	"io/fs"
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

func Redirect(status int, url string) Handler {
	return func(ctx *Context) (Response, error) {
		return NewRedirectResponse(status, url), nil
	}
}

type (
	fsFile struct {
		file fs.File
	}
)

func (f fsFile) Seek(offset int64, whence int) (int64, error) {
	s, ok := f.file.(io.Seeker)
	if !ok {
		return 0, fs.ErrInvalid
	}

	return s.Seek(offset, whence)
}

func (f fsFile) Read(dest []byte) (int, error) {
	return f.file.Read(dest)
}
