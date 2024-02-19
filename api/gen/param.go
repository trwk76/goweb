package gen

import (
	"fmt"

	"github.com/trwk76/goweb/api/spec"
)

func NamedParam[T any](s *APISpec, key string, name string, in spec.ParameterIn, desc string, required bool, deprecated bool, examples Examples[T]) ParamSpec {
	if s.Params.hasName(key) {
		panic(fmt.Errorf("api already declares a parameter named '%s'", key))
	}

	itm := &namedParam{
		name: key,
		impl: param(s, name, in, desc, required, deprecated, examples),
	}
	s.Params = append(s.Params, itm)

	return paramRef{ptr: itm}
}

func Param[T any](s *APISpec, name string, in spec.ParameterIn, desc string, required bool, deprecated bool, examples Examples[T]) ParamSpec {
	return param(s, name, in, desc, required, deprecated, examples)
}

type (
	ParamSpecs []*namedParam

	namedParam struct {
		name string
		impl paramImpl
	}

	ParamSpec interface {
		spec() spec.Parameter
		impl() paramImpl
	}

	paramImpl struct {
		name string
		in   spec.ParameterIn
		desc string
		depr bool
		req  bool
		sch  SchemaSpec
		exs  spec.NamedExamples
	}

	paramRef struct {
		ptr *namedParam
	}
)

func (p *ParamSpecs) hasName(name string) bool {
	for _, itm := range *p {
		if itm.name == name {
			return true
		}
	}

	return false
}

func (p *ParamSpecs) spec() spec.NamedParameters {
	res := make(spec.NamedParameters)

	for _, itm := range *p {
		res[itm.name] = itm.impl.spec()
	}

	return res
}

func (p paramImpl) spec() spec.Parameter {
	return spec.Parameter{
		Name:        p.name,
		In:          p.in,
		Description: p.desc,
		Required:    p.req,
		Deprecated:  p.depr,
		Schema:      p.sch.spec(),
		Examples:    p.exs,
	}
}

func (p paramImpl) impl() paramImpl {
	return p
}

func (r paramRef) spec() spec.Parameter {
	return spec.Parameter{Ref: "#/components/parameters/" + r.ptr.name}
}

func (r paramRef) impl() paramImpl {
	return r.ptr.impl
}

func param[T any](s *APISpec, name string, in spec.ParameterIn, desc string, required bool, deprecated bool, examples Examples[T]) paramImpl {
	sch := Schema[T](s)

	return paramImpl{
		name: name,
		in:   in,
		desc: desc,
		req:  required,
		depr: deprecated,
		sch:  sch,
		exs:  examples.spec(),
	}
}

var (
	_ ParamSpec = paramImpl{}
	_ ParamSpec = paramRef{}
)
