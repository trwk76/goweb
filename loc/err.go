package loc

type (
	OverflowError struct{}
)

func (OverflowError) Error() string {
	return "value is too big to fit the target type"
}

var (
	_ error = OverflowError{}
)
