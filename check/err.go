package check

import (
	"fmt"
)

type (
	FieldErrors map[string]error

	RequiredError struct{}

	ValueCountError struct {
		Expected int
		Actual   int
	}
)

func NewFieldErrors() FieldErrors {
	return make(FieldErrors)
}

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

func (e FieldErrors) Err() error {
	if len(e) < 1 {
		return nil
	}

	return e
}

func (FieldErrors) Error() string {
	return "one or more fields failed to validate"
}

func (RequiredError) Error() string {
	return "requires a value"
}

func (e ValueCountError) Error() string {
	return fmt.Sprintf("expected %d values, got %d", e.Expected, e.Actual)
}

var (
	_ error = RequiredError{}
	_ error = ValueCountError{}
)
