package check

// WithCtx checks the value using the given context by invoking the CheckableCtx or Checkable interfaces if the value implements them, otherwise returns nil.
func WithCtx[C any](ctx C, value any) error {
	if cc, ok := value.(CheckableCtx[C]); ok {
		return cc.CheckCtx(ctx)
	} else if c, ok := value.(Checkable); ok {
		return c.Check()
	}

	return nil
}

// With checks the value using the given context by invoking the CheckableCtx or Checkable interfaces if the value implements them, otherwise returns nil.
func With(value any) error {
	return WithCtx[any](nil, value)
}
