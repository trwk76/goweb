package web

func BasicSecurity(s *Server, key string, check BasicCredCheck) {
	s.secp[key] = httpBasicSec{
		chck: check,
	}
}

type (
	BasicCredCheck func(c *Context, userName string, password string) (Principal, error)

	httpBasicSec struct {
		chck BasicCredCheck
	}
)

func (p httpBasicSec) Authenticate(c *Context) (Principal, error) {
	login, pwd, ok := c.req.BasicAuth()
	if !ok {
		return nil, nil
	}

	return p.chck(c, login, pwd)
}

var _ SecurityProvider = httpBasicSec{}
