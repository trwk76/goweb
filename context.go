package web

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (c *Context) Time() time.Time {
	return c.tm
}

func (c *Context) CorrID() string {
	return c.corrID
}

func (c *Context) Request() *http.Request {
	return c.req
}

type (
	Context struct {
		srv    *Server
		tm     time.Time
		corrID string
		req    *http.Request
		params map[string]string
	}
)

func newContext(srv *Server, req *http.Request) Context {
	return Context{
		srv:    srv,
		tm:     time.Now(),
		corrID: uuid.NewString(),
		req:    req,
		params: make(map[string]string),
	}
}

func init() {
	uuid.EnableRandPool()
}
