package gen

import (
	"reflect"

	"github.com/trwk76/goweb"
	"github.com/trwk76/goweb/openapi/spec"
)

func RespFor[T any](a *API, desc string, setup SetupFunc[spec.Response]) spec.ResponseOrRef {
	return RespOf(a, reflect.TypeFor[T](), desc, setup)
}

func RespOf(a *API, t reflect.Type, desc string, setup SetupFunc[spec.Response]) spec.ResponseOrRef {
	sch := a.sch.SchemaOrRefOf(t)

	item := spec.Response{
		Description: desc,
		Content: spec.MediaTypes{
			goweb.ContentTypeJSON: &spec.MediaType{Schema: &sch},
		},
	}

	if setup != nil {
		setup(a, &item)
	}

	return spec.ResponseOrRef{Item: item}
}

func (r *Resps) Add(key string, item spec.ResponseOrRef) spec.ResponseOrRef {
	if item.Ref.Ref != "" {
		return item
	}

	key = uniqueName(*r, key)
	(*r)[key] = item

	return spec.ResponseOrRef{Ref: spec.ComponentsRef("responses", key)}
}

func (r *Resps) Resp(item spec.ResponseOrRef) (spec.Response, bool) {
	if item.Ref.Ref != "" {
		key, ok := spec.ComponentsKey(item.Ref.Ref, "responses")
		if !ok {
			return spec.Response{}, false
		}

		res, ok := (*r)[key]
		return res.Item, ok
	}

	return item.Item, true
}

type (
	Resps spec.NamedResponseOrRefs
)

func (r Resps) spec() spec.NamedResponseOrRefs {
	return spec.NamedResponseOrRefs(r)
}
