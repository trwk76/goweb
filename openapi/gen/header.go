package gen

import "github.com/trwk76/goweb/openapi/spec"

func (a *API) HeaderOrRef(key string, spc spec.Header) spec.HeaderOrRef {
	if key != "" {
		key = uniqueName(a.spc.Components.Headers, key)
		a.spc.Components.Headers[key] = spec.HeaderOrRef{Item: spc}
		return spec.HeaderOrRef{Ref: spec.Ref("headers", key)}
	}

	return spec.HeaderOrRef{Item: spc}
}

func (a *API) Header(s spec.HeaderOrRef) spec.Header {
	if s.Ref.Ref != "" {
		return a.spc.Components.Headers[refKey(s.Ref.Ref)].Item
	}

	return s.Item
}
