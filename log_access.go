package web

type (
	AccessWriter interface {
		Begin(ctx *Context)
		End(ctx *Context, res Response, msecs int64)
	}
)

func LogAccessBegin(ctx *Context) {
	if accessWriter == nil {
		return
	}

	defer recover()
	accessWriter.Begin(ctx)
}

func LogAccessEnd(ctx *Context, res Response, msecs int64) {
	if accessWriter == nil {
		return
	}

	defer recover()
	accessWriter.End(ctx, res, msecs)
}

func SetAccessWriter(w AccessWriter) {
	accessWriter = w
}

var accessWriter AccessWriter
