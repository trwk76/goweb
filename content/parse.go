package content

import (
	"net/http"

	"github.com/trwk76/goweb/loc"
)

type (
	ParseFuncCtx[C any] func(ctx C, l loc.Loc, r *http.Request, dest any) error
	ParseFunc           func(l loc.Loc, r *http.Request, dest any) error
)
