package gen

import "github.com/trwk76/goweb/openapi/spec"

func ResponseFor[T any](a *API, desc string, hdrs spec.NamedHeaderOrRefs, mtypes MediaTypes, examples Examples[T]) spec.Response {
	return spec.Response{
		Description: desc,
		Headers:     hdrs,
		Content:     ContentFor(a, mtypes, examples),
	}
}

func (a *API) ResponseOrRef(key string, spc spec.Response) spec.ResponseOrRef {
	if key != "" {
		key = uniqueName(a.spc.Components.Responses, key)
		a.spc.Components.Responses[key] = spec.ResponseOrRef{Item: spc}
		return spec.ResponseOrRef{Ref: spec.Ref("responses", key)}
	}

	return spec.ResponseOrRef{Item: spc}
}

func (a *API) Response(s spec.ResponseOrRef) spec.Response {
	if s.Ref.Ref != "" {
		return a.spc.Components.Responses[refKey(s.Ref.Ref)].Item
	}

	return s.Item
}
