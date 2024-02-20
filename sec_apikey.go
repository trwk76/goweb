package web

import (
	"fmt"

	"github.com/trwk76/goweb/api/spec"
)

func APIKeySecurity(s *Server, key string, desc string, name string, in spec.SecurityAPIKeyIn, check APIKeyCredCheck) {
	s.secp[key] = apiKeySec{
		desc: desc,
		name: name,
		in:   in,
		chck: check,
	}
}

type (
	APIKeyCredCheck func(c *Context, key string) (Principal, error)

	apiKeySec struct {
		desc string
		name string
		in   spec.SecurityAPIKeyIn
		chck APIKeyCredCheck
	}
)

func (p apiKeySec) APISpec() spec.SecurityScheme {
	return spec.SecurityScheme{
		Type:        spec.SecurityTypeAPIKey,
		Description: p.desc,
		Name:        p.name,
		In:          p.in,
	}
}

func (p apiKeySec) Authenticate(c *Context) (Principal, error) {
	var key string

	switch p.in {
	case spec.SecurityAPIKeyInCookie:
		if cook, err := c.req.Cookie(p.name); err == nil {
			key = cook.Value
		}
	case spec.SecurityAPIKeyInHeader:
		key = c.req.Header.Get(p.name)
	case spec.SecurityAPIKeyInQuery:
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
