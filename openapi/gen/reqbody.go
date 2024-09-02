package gen

import (
	"reflect"

	"github.com/trwk76/goweb"
	"github.com/trwk76/goweb/openapi/spec"
)

func ReqBodyFor[T any](a *API, setup SetupFunc[spec.RequestBody]) spec.RequestBodyOrRef {
	return ReqBodyOf(a, reflect.TypeFor[T](), setup)
}

func ReqBodyOf(a *API, t reflect.Type, setup SetupFunc[spec.RequestBody]) spec.RequestBodyOrRef {
	sch := a.sch.SchemaOrRefOf(t)

	item := spec.RequestBody{
		Content: spec.MediaTypes{
			goweb.ContentTypeJSON: spec.MediaType{
				Schema: &sch,
			},
		},
		Required: true,
	}

	if setup != nil {
		setup(a, &item)
	}

	return spec.RequestBodyOrRef{Item: item}
}

func (r *ReqBodies) Add(key string, item spec.RequestBodyOrRef) spec.RequestBodyOrRef {
	if item.Ref.Ref != "" {
		return item
	}

	key = uniqueName(*r, key)
	(*r)[key] = item

	return spec.RequestBodyOrRef{Ref: spec.ComponentsRef("requestBodies", key)}
}

func (r *ReqBodies) ReqBody(item spec.RequestBodyOrRef) (spec.RequestBody, bool) {
	if item.Ref.Ref != "" {
		key, ok := spec.ComponentsKey(item.Ref.Ref, "requestBodies")
		if !ok {
			return spec.RequestBody{}, false
		}

		res, ok := (*r)[key]
		return res.Item, ok
	}

	return item.Item, true
}

type (
	ReqBodies spec.NamedRequestBodyOrRefs
)

func (r ReqBodies) spec() spec.NamedRequestBodyOrRefs {
	return spec.NamedRequestBodyOrRefs(r)
}
