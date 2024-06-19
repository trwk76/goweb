package gen

import "github.com/trwk76/goweb/openapi/spec"

func ParamFor[T any](a *API, name string, in spec.ParameterIn, f func(s *spec.Parameter), examples Examples[T]) spec.Parameter {
	sch := SchemaFor[T](a, "", nil)

	res := spec.Parameter{
		Name:     name,
		In:       in,
		Schema:   &sch,
		Examples: examples.spec(JSON),
	}

	f(&res)
	return res
}

func (a *API) ParamOrRef(key string, spc spec.Parameter) spec.ParameterOrRef {
	if key != "" {
		key = uniqueName(a.spc.Components.Parameters, key)
		a.spc.Components.Parameters[key] = spec.ParameterOrRef{Item: spc}
		return spec.ParameterOrRef{Ref: spec.Ref("parameters", key)}
	}

	return spec.ParameterOrRef{Item: spc}
}

func (a *API) Param(s spec.ParameterOrRef) spec.Parameter {
	if s.Ref.Ref != "" {
		return a.spc.Components.Parameters[refKey(s.Ref.Ref)].Item
	}

	return s.Item
}
