package models

import (
	"encoding/json"
	"reflect"

	"github.com/trwk76/goweb/loc"
)

// MarshalJSONBool enables to implement the json.Marshaler interface for bool in a single line.
func MarshalJSONBool(v bool) ([]byte, error) {
	return json.Marshal(v)
}

// UnmarshalJSONBool enables to implement the json.Unmarshaler interface for bool in a single line.
func UnmarshalJSONBool(v *bool, data []byte) error {
	if err := json.Unmarshal(data, v); err != nil {
		return FormatError{Value: string(data), Type: reflect.TypeFor[bool]()}
	}

	return nil
}

// MarshalTextBool enables to implement the loc.TextMarshaler and encoding.TextMarshaler interfaces for bool in a single line.
func MarshalTextBool(v bool, l loc.Loc) ([]byte, error) {
	return []byte(l.FormatBool(v)), nil
}

// UnmarshalTextBool enables to implement the loc.TextUnmarshaler and encoding.TextUnmarshaler interfaces for bool in a single line.
func UnmarshalTextBool(v *bool, l loc.Loc, data []byte) error {
	tmp, err := l.ParseBool(string(data))
	if err != nil {
		return FormatError{Value: string(data), Type: reflect.TypeFor[bool]()}
	}

	*v = tmp
	return nil
}
