package models

import (
	"fmt"
	"reflect"
)

type (
	// FormatError is an error that indicates a value is not valid for a specific type.
	FormatError struct {
		Value string
		Type  reflect.Type
	}
)

// Error returns the error message for the FormatError.
func (e FormatError) Error() string {
	return fmt.Sprintf("'%s' is not a valid %s", e.Value, e.Type.String())
}

var (
	_ error = FormatError{}
)
