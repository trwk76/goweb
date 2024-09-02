package gen

import (
	"reflect"

	"github.com/trwk76/goweb/openapi/spec"
)

func ParamFor[T any](a *API, name string, in spec.ParameterIn, setup func(item *spec.Parameter)) spec.ParameterOrRef {
	return ParamOf(a, reflect.TypeFor[T](), name, in, setup)
}

func ParamOf(a *API, t reflect.Type, name string, in spec.ParameterIn, setup func(item *spec.Parameter)) spec.ParameterOrRef {
	sch := SchemaOf(t, &a.sch)

	checkSimpleSchema(sch)

	item := spec.Parameter{
		Name:   name,
		In:     in,
		Schema: &sch,
	}

	if setup != nil {
		setup(&item)
	}

	return spec.ParameterOrRef{Item: item}
}

func (p *Params) Add(key string, item spec.ParameterOrRef) spec.ParameterOrRef {
	if item.Ref.Ref != "" {
		return item
	}

	key = uniqueName(*p, key)
	(*p)[key] = item

	return spec.ParameterOrRef{Ref: spec.ComponentsRef("parameters", key)}
}

func (p *Params) Param(item spec.ParameterOrRef) (spec.Parameter, bool) {
	if item.Ref.Ref != "" {
		key, ok := spec.ComponentsKey(item.Ref.Ref, "parameters")
		if !ok {
			return spec.Parameter{}, false
		}

		res, ok := (*p)[key]
		return res.Item, ok
	}

	return item.Item, true
}

type (
	Params spec.NamedParameterOrRefs
)

func (p Params) spec() spec.NamedParameterOrRefs {
	return spec.NamedParameterOrRefs(p)
}
