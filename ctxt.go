package web

import (
	"net/http"

	"github.com/google/uuid"
)

type (
	Context struct {
		id    string
		req   *http.Request
		prins Principals
		errh  ErrorHandler
		parm  map[string]string
	}
)

func (c *Context) ID() string {
	return c.id
}

func (c *Context) Request() *http.Request {
	return c.req
}

func (c *Context) PathParam(name string) string {
	return c.parm[name]
}

func (c *Context) Principals() Principals {
	return c.prins
}

func (c *Context) Error(err error) Response {
	if c.errh != nil {
		return c.errh(c, err)
	}

	return DefaultResponse(EnsureError(err).Status)
}

func newContext(req *http.Request) Context {
	return Context{
		id:    uuid.NewString(),
		req:   req,
		prins: make(Principals),
		parm:  make(map[string]string),
	}
}

func init() {
	uuid.EnableRandPool()
}
