package web

import "github.com/trwk76/goweb/api/spec"

type (
	SecurityRequirements []SecurityRequirement
	SecurityRequirement  map[string][]string

	Principals map[string]Principal

	Principal interface {
		HasRole(role string) bool
	}

	SecurityProviders map[string]SecurityProvider

	SecurityProvider interface {
		Authenticate(c *Context) (Principal, error)
		APISpec() spec.SecurityScheme
	}
)

func (p Principals) satisfyAny(reqs SecurityRequirements) bool {
	for _, req := range reqs {
		if p.satisfy(req) {
			return true
		}
	}

	return false
}

func (p Principals) satisfy(req SecurityRequirement) bool {
	for key, roles := range req {
		prin, ok := p[key]
		if !ok {
			return false
		}

		for _, role := range roles {
			if !prin.HasRole(role) {
				return false
			}
		}
	}

	return true
}

func (p SecurityProviders) APISpec() spec.NamedSecuritySchemes {
	res := make(spec.NamedSecuritySchemes)

	for key, prv := range p {
		res[key] = prv.APISpec()
	}

	return res
}

func (p SecurityProviders) authenticate(c *Context) (Principals, error) {
	res := make(Principals)

	for key, prv := range p {
		if prin, err := prv.Authenticate(c); err != nil {
			return nil, err
		} else if prin != nil {
			res[key] = prin
		}
	}

	return res, nil
}
