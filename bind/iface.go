package bind

import (
	"net/http"
)

type (
	// BindableCtx is an interface for defining custom binding logic for a structure using the given context.
	// If a structure pointer implements this interface, the Parse method will call it instead of parsing the structure.
	BindableCtx[C any] interface {
		BindCtx(ctx C, r *http.Request) error
	}

	// Bindable is an interface for defining custom binding logic for a structure.
	// If a structure pointer implements this interface, the Parse method will call it instead of parsing the structure.
	Bindable interface {
		Bind(r *http.Request) error
	}
)
