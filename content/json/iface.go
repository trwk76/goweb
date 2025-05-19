package json

import "github.com/trwk76/goweb/loc"

type (
	// JSONUnmarshalerCtx is an interface for defining custom JSON marshaling logic for a value using the given context.
	// If a value pointer implements this interface, the Unmarshal method will call it instead of using default behavior.
	JSONUnmarshalerCtx[C any] interface {
		UnmarshalJSONCtx(ctx any, l loc.Loc, v []byte) error
	}

	// JSONUnmarshaler is an interface for defining custom JSON marshaling logic for a value without a context.
	// If a value pointer implements this interface, the Unmarshal method will call it instead of using default behavior.
	JSONUnmarshaler interface {
		UnmarshalJSON(l loc.Loc, v []byte) error
	}
)
