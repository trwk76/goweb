package check

type (
	CheckableCtx[C any] interface {
		CheckCtx(ctx C) error
	}

	Checkable interface {
		Check() error
	}
)
