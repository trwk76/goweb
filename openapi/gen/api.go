package gen

import "github.com/trwk76/goweb/openapi/spec"

func New(path string, info spec.Info) *API {
	res := &API{
		spc: spec.OpenAPI{
			OpenAPI: spec.Version,
			Info:    info,
			Servers: []spec.Server{{URL: path, Description: "Current server"}},
			Paths:   make(spec.Paths),
			Components: &spec.Components{
				Schemas:         make(spec.NamedSchemas),
				Responses:       make(spec.NamedResponseOrRefs),
				Parameters:      make(spec.NamedParameterOrRefs),
				Examples:        make(spec.NamedExampleOrRefs),
				RequestBodies:   make(spec.NamedRequestBodyOrRefs),
				Headers:         make(spec.NamedHeaderOrRefs),
				SecuritySchemes: make(spec.NamedSecuritySchemeOrRefs),
			},
		},
		sch: make(schemas),
	}

	return res
}

type (
	API struct {
		spc spec.OpenAPI
		sch schemas
	}
)
