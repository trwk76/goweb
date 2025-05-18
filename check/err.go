package check

import (
	"fmt"
)

type (
	// FieldErrors is a map of field paths to errors.
	// It is used to collect validation errors for multiple fields in a structure.
	// The keys are the field paths (e.g., "field1/field2") and the values are the corresponding errors.
	// The map is used to store errors for each field, allowing for easy retrieval and display of validation errors.
	// The map is also used to implement the error interface, allowing it to be returned as an error.
	FieldErrors map[string]error

	// RequiredError is an error that indicates a required field is missing a value.
	RequiredError struct{}

	// ValueCountError is an error that indicates the number of values provided does not match the expected count.
	ValueCountError struct {
		Expected int
		Actual   int
	}
)

// NewFieldErrors creates a new instance of FieldErrors.
func NewFieldErrors() FieldErrors {
	return make(FieldErrors)
}

// Add adds an error into the receiver FieldErrors.
// If the error is nil, it removes the field from the map.
// If the error is a FieldErrors, it merges the errors into the receiver map using name as prefixes.
// If the error is a single error, it adds it to the map with the given name.
func (e FieldErrors) Add(name string, err error) {
	if err == nil {
		delete(e, name)
	} else if fe, ok := err.(FieldErrors); ok {
		name += "/"

		for key, val := range fe {
			e[name+key] = val
		}
	} else {
		e[name] = err
	}
}

// Err returns the receiver as an error if any at least one field has an error.
func (e FieldErrors) Err() error {
	if len(e) < 1 {
		return nil
	}

	return e
}

// Error returns the error message for the receiver FieldErrors.
func (FieldErrors) Error() string {
	return "one or more fields failed to validate"
}

// Error returns the error message for the RequiredError.
func (RequiredError) Error() string {
	return "requires a value"
}

// Error returns the error message for the ValueCountError.
func (e ValueCountError) Error() string {
	return fmt.Sprintf("expected %d values, got %d", e.Expected, e.Actual)
}

var (
	_ error = RequiredError{}
	_ error = ValueCountError{}
)
