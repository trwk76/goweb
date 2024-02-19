package gen

import (
	"fmt"

	web "github.com/trwk76/goweb"
	"github.com/trwk76/goweb/api/spec"
)

func NamedReqBody[T any](s *APISpec, name string, desc string, required bool, examples Examples[T]) ReqBodySpec {
	if s.ReqBodies.hasName(name) {
		panic(fmt.Errorf("api already declares a request body named '%s'", name))
	}

	itm := &namedReqBody{
		name: name,
		impl: reqBody(s, desc, required, examples),
	}
	s.ReqBodies = append(s.ReqBodies, itm)

	return reqBodyRef{ptr: itm}
}

func ReqBody[T any](s *APISpec, desc string, req bool, examples Examples[T]) ReqBodySpec {
	return reqBody[T](s, desc, req, examples)
}

type (
	ReqBodySpecs []*namedReqBody

	ReqBodySpec interface {
		spec() spec.RequestBody
		impl() reqBodyImpl
	}

	namedReqBody struct {
		name string
		impl reqBodyImpl
	}

	reqBodyImpl struct {
		desc string
		req  bool
		sch  SchemaSpec
		exs  spec.NamedExamples
	}

	reqBodyRef struct {
		ptr *namedReqBody
	}
)

func (p *ReqBodySpecs) hasName(name string) bool {
	for _, itm := range *p {
		if itm.name == name {
			return true
		}
	}

	return false
}

func (p *ReqBodySpecs) spec() spec.NamedRequestBodies {
	res := make(spec.NamedRequestBodies)

	for _, itm := range *p {
		res[itm.name] = itm.impl.spec()
	}

	return res
}

func (r reqBodyImpl) spec() spec.RequestBody {
	return spec.RequestBody{
		Description: r.desc,
		Required:    r.req,
		Content: spec.MediaTypes{
			web.JSONContentType: spec.MediaType{
				Schema:   r.sch.spec(),
				Examples: r.exs,
			},
		},
	}
}

func (r reqBodyImpl) impl() reqBodyImpl {
	return r
}

func (r reqBodyRef) spec() spec.RequestBody {
	return spec.RequestBody{Ref: "#/components/requestBodies/" + r.ptr.name}
}

func (r reqBodyRef) impl() reqBodyImpl {
	return r.ptr.impl
}

func reqBody[T any](s *APISpec, desc string, required bool, examples Examples[T]) reqBodyImpl {
	sch := Schema[T](s)

	return reqBodyImpl{
		desc: desc,
		req:  required,
		sch:  sch,
		exs:  examples.spec(),
	}
}

var (
	_ ReqBodySpec = reqBodyImpl{}
	_ ReqBodySpec = reqBodyRef{}
)
