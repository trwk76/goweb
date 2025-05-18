package loc

type (
	// OverflowError is an error that indicates a value is too large to fit into the target type.
	OverflowError struct{}
)

// Error returns the error message for the OverflowError.
func (OverflowError) Error() string {
	return "value is too big to fit the target type"
}

var (
	_ error = OverflowError{}
)
