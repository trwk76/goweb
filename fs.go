package web

import (
	"fmt"
	"io/fs"
	"mime"
	"net/http"
	"path/filepath"
	"time"
)

func FS(p Path, f fs.FS, errHandler ErrorHandler) {
	p = p.ParamPath("path", true, errHandler)
	GET_HEAD(p, func(ctx *Context) (Response, error) {
		path := ctx.Param("path")
		return serveFile(ctx, f, path)
	})
}

func File(f fs.FS, path string, etag string) Handler {
	return func(ctx *Context) (Response, error) {
		res, err := serveFile(ctx, f, path)
		SetETag(res.hdr, etag)
		return res, err
	}
}

func serveFile(ctx *Context, f fs.FS, path string) (*ResponseBuffer, error) {
	res := NewRespBuffer()
	modTime := time.Time{}

	file, err := f.Open(path)
	if err != nil {
		return &res, fmt.Errorf("error opening file '%s': %s", path, err.Error())
	}

	defer file.Close()

	stat, err := file.Stat()
	if err == nil {
		modTime = stat.ModTime()
	}

	SetContentType(res.hdr, mime.TypeByExtension(filepath.Ext(path)))
	http.ServeContent(&res, ctx.req, filepath.Base(path), modTime, fsFile{file: file})
	return &res, nil
}
