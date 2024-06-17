package log

import web "github.com/trwk76/goweb"

type (
	AccessWriter interface {
		Begin(ctx *web.Context)
		End(ctx *web.Context, res web.Response)
	}
)
