package web

import (
	"fmt"
	"net/http"
	"strings"
)

type (
	Server struct {
		Path
		secp SecurityProviders
	}
)

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var res Response

	ctx := newContext(r)
	cur := &s.Path

	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("request panic: %v", r)
			}

			res = ctx.Error(err)
		}

		if res.Status == 0 {
			res = DefaultResponse(http.StatusNotImplemented)
		}

		res.write(w)
	}()

	pth := parsePath(r.URL.Path)

	for len(pth) > 0 {
		if cur.errh != nil {
			ctx.errh = cur.errh
		}

		switch pth[0] {
		case ".":
			pth = pth[1:]
		case "..":
			if cur.par != nil {
				cur = cur.par
			}
			pth = pth[1:]
		default:
			cur, pth = cur.child(&ctx, pth)
		}

		if cur == nil {
			break
		}
	}

	if cur == nil || len(cur.mth) < 1 {
		res = ctx.Error(Error{Status: http.StatusNotFound})
		return
	}

	mth, ok := cur.mth[strings.ToUpper(r.Method)]
	if !ok {
		res = ctx.Error(Error{Status: http.StatusMethodNotAllowed})
		return
	}

	if prins, err := s.secp.authenticate(&ctx); err != nil {
		res = ctx.Error(err)
		return
	} else {
		ctx.prins = prins
	}

	if mth.sec != nil && !ctx.prins.satisfyAny(mth.sec) {
		res = ctx.Error(Error{Status: http.StatusUnauthorized})
		return
	}

	res = mth.hdl(&ctx)
}
