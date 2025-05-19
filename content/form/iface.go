package form

import (
	"github.com/trwk76/goweb/loc"
)

type (
	// URLUnmarshalerCtx is an interface for defining custom URL marshaling logic for a structure using the given context.
	// If a structure pointer implements this interface, the Parse method will call it instead of parsing the structure.
	URLUnmarshalerCtx[C any] interface {
		UnmarshalURLCtx(ctx any, l loc.Loc, v string) error
	}

	// URLUnmarshaler is an interface for defining custom URL marshaling logic for a structure.
	// If a structure pointer implements this interface, the Parse method will call it instead of parsing the structure.
	URLUnmarshaler interface {
		UnmarshalURL(l loc.Loc, v string) error
	}
)
