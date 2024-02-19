package gen

import (
	"fmt"

	web "github.com/trwk76/goweb"
	"github.com/trwk76/goweb/api/spec"
)

func NamedResp[T any](s *APISpec, name string, desc string, examples Examples[T]) RespSpec {
	if s.Resps.hasName(name) {
		panic(fmt.Errorf("api already declares a response named '%s'", name))
	}

	itm := &namedResp{
		name: name,
		impl: resp[T](s, desc, examples),
	}
	s.Resps = append(s.Resps, itm)

	return respRef{ptr: itm}
}

func Resp[T any](s *APISpec, desc string, examples Examples[T]) RespSpec {
	return resp[T](s, desc, examples)
}

type (
	RespSpecs []*namedResp

	RespSpec interface {
		spec() spec.Response
		impl() respImpl
	}

	namedResp struct {
		name string
		impl respImpl
	}

	respImpl struct {
		desc string
		sch  SchemaSpec
		exs  spec.NamedExamples
	}

	respRef struct {
		ptr *namedResp
	}
)

func (r *RespSpecs) hasName(name string) bool {
	for _, itm := range *r {
		if itm.name == name {
			return true
		}
	}

	return false
}

func (p *RespSpecs) spec() spec.NamedResponses {
	res := make(spec.NamedResponses)

	for _, itm := range *p {
		res[itm.name] = itm.impl.spec()
	}

	return res
}

func (r respImpl) spec() spec.Response {
	return spec.Response{
		Description: r.desc,
		Content: spec.MediaTypes{
			web.JSONContentType: spec.MediaType{
				Schema:   r.sch.spec(),
				Examples: r.exs,
			},
		},
	}
}

func (r respImpl) impl() respImpl {
	return r
}

func (r respRef) spec() spec.Response {
	return spec.Response{Ref: "#/components/responses/" + r.ptr.name}
}

func (r respRef) impl() respImpl {
	return r.ptr.impl
}

func resp[T any](s *APISpec, desc string, examples Examples[T]) respImpl {
	return respImpl{
		desc: desc,
		sch:  Schema[T](s),
		exs:  examples.spec(),
	}
}

var (
	_ RespSpec = respImpl{}
	_ RespSpec = respRef{}
)
