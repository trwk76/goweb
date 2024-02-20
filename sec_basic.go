package web

import "github.com/trwk76/goweb/api/spec"

func BasicSecurity(s *Server, key string, desc string, check BasicCredCheck) {
	s.secp[key] = httpBasicSec{
		desc: desc,
		chck: check,
	}
}

type (
	BasicCredCheck func(c *Context, userName string, password string) (Principal, error)

	httpBasicSec struct {
		desc string
		chck BasicCredCheck
	}
)

func (p httpBasicSec) APISpec() spec.SecurityScheme {
	return spec.SecurityScheme{
		Type:        spec.SecurityTypeHTTP,
		Description: p.desc,
		Scheme:      spec.SecurityHTTPSchemeBasic,
	}
}

func (p httpBasicSec) Authenticate(c *Context) (Principal, error) {
	login, pwd, ok := c.req.BasicAuth()
	if !ok {
		return nil, nil
	}

	return p.chck(c, login, pwd)
}

var _ SecurityProvider = httpBasicSec{}
