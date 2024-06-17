package log

import web "github.com/trwk76/goweb"

type (
	AccessWriter interface {
		Begin(ctx *web.Context)
		End(ctx *web.Context, res web.Response)
	}
)

func LogAccessBegin(ctx *web.Context) {
	if accessWriter == nil {
		return
	}

	defer recover()
	accessWriter.Begin(ctx)
}

func LogAccessEnd(ctx *web.Context, res web.Response) {
	if accessWriter == nil {
		return
	}

	defer recover()
	accessWriter.End(ctx, res)
}

func SetAccessWriter(w AccessWriter) {
	accessWriter = w
}

var accessWriter AccessWriter
