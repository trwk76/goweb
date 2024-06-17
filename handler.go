package web

import (
	"fmt"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"path/filepath"
	"time"
)

type (
	Handler      func(ctx *Context) (Response, error)
	ErrorHandler func(ctx *Context, err error) Response
)

func File(f fs.FS, path string, etag string) Handler {
	ctype := mime.TypeByExtension(filepath.Ext(path))

	return func(ctx *Context) (Response, error) {
		modTime := time.Time{}
		file, err := f.Open(path)
		if err != nil {
			panic(fmt.Errorf("error opening file '%s': %s", path, err.Error()))
		}

		defer file.Close()

		stat, err := file.Stat()
		if err == nil {
			modTime = stat.ModTime()
		}

		res := NewRespBuffer()

		SetContentType(res.hdr, ctype)
		SetETag(res.hdr, etag)
		http.ServeContent(&res, ctx.req, filepath.Base(path), modTime, fsFile{file: file})

		return &res, nil
	}
}

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
