package web

import (
	"net/http"

	"github.com/google/uuid"
)

type (
	Context struct {
		id   string
		req  *http.Request
		parm map[string]string
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

func newContext(req *http.Request) Context {
	return Context{
		id:   uuid.NewString(),
		req:  req,
		parm: make(map[string]string),
	}
}

func init() {
	uuid.EnableRandPool()
}
