package gen

import "github.com/trwk76/goweb/openapi/spec"

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
