package web

import (
	"fmt"
	"net/http"
	"strings"
)

type (
	Server struct {
		root path
	}
)

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var res response

	ctx := newContext(s, r)
	nod := &s.root
	errh := nod.errh

	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("request panic: %v", r)
			}

			res = errh(&ctx, err)
		}

		res.Write(w)
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

	if nod == nil {
		res = errh(&ctx, ErrNotFound)
		return
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
