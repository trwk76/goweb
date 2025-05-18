package check

type (
	// CheckableCtx is an interface that defines a method for checking a value using a context.
	CheckableCtx[C any] interface {
		CheckCtx(ctx C) error
	}

	// Checkable is an interface that defines a method for checking a value without a context.
	Checkable interface {
		Check() error
	}
)
