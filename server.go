package web

import (
	"fmt"
	"net/http"
	"strings"
)

func New(errHandler ErrorHandler) Server {
	if errHandler == nil {
		errHandler = func(ctx *Context, err error) response {
			e, ok := err.(Error)
			if !ok {
				e = ErrInternalError
			}

			return e
		}
	}

	return Server{
		path: path{
			paths: paths{
				named: make(map[string]*path),
			},
			name:  "",
			meths: make(map[string]Handler),
			errh:  errHandler,
		},
	}
}

type (
	Server struct {
		path
	}
)

func (s *Server) Run(srv *http.Server, certPath string, keyPath string) error {
	srv.Handler = s

	if certPath != "" && keyPath != "" {
		return srv.ListenAndServeTLS(certPath, keyPath)
	}

	return srv.ListenAndServe()
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		res response
		err error
	)

	ctx := newContext(s, r)
	nod := &s.path
	errh := nod.errh

	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("request panic: %v", r)
			}

			res = errh(&ctx, err)
		}

		res.WriteResponse(w)
	}()

	path := parsePath(r.URL.Path)
	for len(path) > 0 {
		if nod == nil {
			break
		}

		if nod.errh != nil {
			errh = nod.errh
		}

		switch path[0] {
		case ".":
			path = path[1:]
		case "..":
			path = path[1:]
			nod = nod.par
		default:
			nod, path = nod.child(&ctx, path)
		}
	}

	if nod == nil || len(nod.meths) < 1 {
		res = errh(&ctx, ErrNotFound)
		return
	}

	hdl, ok := nod.meths[r.Method]
	if !ok {
		res = errh(&ctx, ErrMethNotAllowed)
	}

	if res, err = hdl(&ctx); err != nil {
		res = errh(&ctx, err)
	}
}

func parsePath(path string) []string {
	for strings.Contains(path, "//") {
		path = strings.ReplaceAll(path, "//", "/")
	}

	path = strings.TrimPrefix(path, "/")
	return strings.Split(path, "/")
}

var _ http.Handler = (*Server)(nil)
