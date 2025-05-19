package models

import (
	"encoding/json"
	"reflect"

	"github.com/trwk76/goweb/loc"
)

// MarshalJSONInt enables to implement the json.Marshaler interface for int in a single line.
func MarshalJSONInt[I integer](v I) ([]byte, error) {
	return json.Marshal(v)
}

// UnmarshalJSONInt enables to implement the json.Unmarshaler interface for int in a single line.
func UnmarshalJSONInt[I integer](v *I, data []byte) error {
	if err := json.Unmarshal(data, v); err != nil {
		return FormatError{Value: string(data), Type: reflect.TypeFor[I]()}
	}

	return nil
}

// MarshalTextInt enables to implement the loc.TextMarshaler and encoding.TextMarshaler interfaces for int in a single line.
func MarshalTextInt[I integer](v I, l loc.Loc) ([]byte, error) {
	return []byte(l.FormatInt(int64(v))), nil
}

// UnmarshalTextInt enables to implement the loc.TextUnmarshaler and encoding.TextUnmarshaler interfaces for int in a single line.
func UnmarshalTextInt[I integer](v *I, l loc.Loc, data []byte) error {
	tmp, err := l.ParseInt(string(data))
	if err != nil {
		return FormatError{Value: string(data), Type: reflect.TypeFor[int8]()}
	}

	*v = I(tmp)
	return nil
}

type (
	integer interface {
		~int | ~int8 | ~int16 | ~int32 | ~int64
	}
)
