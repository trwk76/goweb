package gen

import "github.com/trwk76/goweb/openapi/spec"

type (
	Operation struct {
		Summary     string
		Description string
		Deprecated  bool
		Params      spec.ParameterOrRefs
		RequestBody *spec.RequestBody
		Responses   spec.NamedResponseOrRefs
	}
)
