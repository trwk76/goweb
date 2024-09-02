package gen

import (
	"reflect"

	"github.com/trwk76/goweb/openapi/spec"
)

func HeaderFor[T any](a *API, setup SetupFunc[spec.Header]) spec.HeaderOrRef {
	return HeaderOf(a, reflect.TypeFor[T](), setup)
}

func HeaderOf(a *API, t reflect.Type, setup SetupFunc[spec.Header]) spec.HeaderOrRef {
	sch := SchemaOf(t, &a.sch)

	checkSimpleSchema(sch)

	item := spec.Header{
		Schema: &sch,
	}

	if setup != nil {
		setup(a, &item)
	}

	return spec.HeaderOrRef{Item: item}
}

func (h *Headers) Add(key string, item spec.HeaderOrRef) spec.HeaderOrRef {
	if item.Ref.Ref != "" {
		return item
	}

	key = uniqueName(*h, key)
	(*h)[key] = item

	return spec.HeaderOrRef{Ref: spec.ComponentsRef("headers", key)}
}

func (h *Headers) Header(item spec.HeaderOrRef) (spec.Header, bool) {
	if item.Ref.Ref != "" {
		key, ok := spec.ComponentsKey(item.Ref.Ref, "headers")
		if !ok {
			return spec.Header{}, false
		}

		res, ok := (*h)[key]
		return res.Item, ok
	}

	return item.Item, true
}

type (
	Headers spec.NamedHeaderOrRefs
)

func (h Headers) spec() spec.NamedHeaderOrRefs {
	return spec.NamedHeaderOrRefs(h)
}
