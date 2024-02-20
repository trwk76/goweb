package web

import (
	"fmt"
)

func APIKeySecurity(s *Server, key string, name string, in APIKeyIn, check APIKeyCredCheck) {
	s.secp[key] = apiKeySec{
		name: name,
		in:   in,
		chck: check,
	}
}

type (
	APIKeyCredCheck func(c *Context, key string) (Principal, error)

	APIKeyIn string

	apiKeySec struct {
		name string
		in   APIKeyIn
		chck APIKeyCredCheck
	}
)

const (
	APIKeyInCookie APIKeyIn = "cookie"
	APIKeyInHeader APIKeyIn = "header"
	APIKeyInQuery  APIKeyIn = "query"
)

func (p apiKeySec) Authenticate(c *Context) (Principal, error) {
	var key string

	switch p.in {
	case APIKeyInCookie:
		if cook, err := c.req.Cookie(p.name); err == nil {
			key = cook.Value
		}
	case APIKeyInHeader:
		key = c.req.Header.Get(p.name)
	case APIKeyInQuery:
		key = c.req.URL.Query().Get(p.name)
	default:
		panic(fmt.Errorf("invalid api key location '%s'", p.in))
	}

	if key == "" {
		return nil, nil
	}

	return p.chck(c, key)
}

var _ SecurityProvider = apiKeySec{}
