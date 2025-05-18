package bind

import (
	"net/http"

	"github.com/trwk76/goweb/loc"
)

type (
	BindableCtx[C any] interface {
		BindCtx(ctx C, r *http.Request) error
	}

	Bindable interface {
		Bind(r *http.Request) error
	}

	URLUnmarshalerCtx[C any] interface {
		UnmarshalURLCtx(ctx any, l loc.Loc, v string) error
	}

	URLUnmarshaler interface {
		UnmarshalURL(l loc.Loc, v string) error
	}
)
