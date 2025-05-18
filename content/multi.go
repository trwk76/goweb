package content

import (
	"net/http"

	web "github.com/trwk76/goweb"
	"github.com/trwk76/goweb/loc"
)

// MultiParseFuncCtx returns a ParseFunc able to handle different content types using the provided map.
func MultiParseFuncCtx[C any](ctypes map[string]ParseFuncCtx[C]) ParseFuncCtx[C] {
	return func(ctx C, l loc.Loc, r *http.Request, dest any) error {
		typ, err := ParseType(r.Header.Get(web.HeaderContentType))
		if err != nil {
			return err
		}

		for t, f := range ctypes {
			if typ.Type == t {
				return f(ctx, l, r, dest)
			}
		}

		typs := make([]string, 0, len(ctypes))

		for itm := range ctypes {
			typs = append(typs, itm)
		}

		return TypeError{Act: typ.Type, Exp: typs}
	}
}
