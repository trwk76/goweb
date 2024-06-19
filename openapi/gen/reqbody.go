package gen

import "github.com/trwk76/goweb/openapi/spec"

func RequestBodyFor[T any](a *API, desc string, required bool, mtypes MediaTypes, examples Examples[T]) spec.RequestBody {
	return spec.RequestBody{
		Description: desc,
		Content:     ContentFor(a, mtypes, examples),
		Required:    required,
	}
}

func (a *API) RequestBodyOrRef(key string, spc spec.RequestBody) spec.RequestBodyOrRef {
	if key != "" {
		key = uniqueName(a.spc.Components.RequestBodies, key)
		a.spc.Components.RequestBodies[key] = spec.RequestBodyOrRef{Item: spc}
		return spec.RequestBodyOrRef{Ref: spec.Ref("requestBodies", key)}
	}

	return spec.RequestBodyOrRef{Item: spc}
}

func (a *API) RequestBody(s spec.RequestBodyOrRef) spec.RequestBody {
	if s.Ref.Ref != "" {
		return a.spc.Components.RequestBodies[refKey(s.Ref.Ref)].Item
	}

	return s.Item
}
