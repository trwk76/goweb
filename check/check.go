package check

func CheckCtx[C any](ctx C, value any) error {
	if cc, ok := value.(CheckableCtx[C]); ok {
		return cc.CheckCtx(ctx)
	} else if c, ok := value.(Checkable); ok {
		return c.Check()
	}

	return nil
}

func Check(value any) error {
	return CheckCtx[any](nil, value)
}
