package gen

import "github.com/trwk76/goweb/api/spec"

type (
	AuthSpecs map[string]AuthSpec

	AuthSpec struct {
		s spec.SecurityScheme
	}
)

func (a AuthSpecs) spec() spec.NamedSecuritySchemes {
	res := make(spec.NamedSecuritySchemes)

	for key, auth := range a {
		res[key] = auth.s
	}

	return res
}
