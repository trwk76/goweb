package bind

import "net/http"

type (
	BindableCtx[C any] interface {
		BindCtx(ctx C, r *http.Request) error
	}

	Bindable interface {
		Bind(r *http.Request) error
	}
)
