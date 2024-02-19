package gen

import (
	"fmt"
	"net/http"

	"github.com/trwk76/gocode/golang"
	"github.com/trwk76/goweb/api/spec"
)

func Path(s *APISpec, p *PathSpec, name string, param bool, sec spec.SecurityRequirements, tags []string) *PathSpec {
	res := &PathSpec{
		name: name,
		parm: param,
		sec:  sec,
		tags: tags,
	}

	if p != nil {
		p.sub.add(res)
	} else {
		s.Paths.add(res)
	}

	return res
}

func GET(p *PathSpec, op *OpSpec) {
	setOp(&p.get, op, http.MethodGet, false)
}

func PUT(p *PathSpec, op *OpSpec) {
	setOp(&p.put, op, http.MethodPut, true)
}

func POST(p *PathSpec, op *OpSpec) {
	setOp(&p.post, op, http.MethodPost, true)
}

func DELETE(p *PathSpec, op *OpSpec) {
	setOp(&p.del, op, http.MethodDelete, false)
}

func OPTIONS(p *PathSpec, op *OpSpec) {
	setOp(&p.opts, op, http.MethodOptions, false)
}

func HEAD(p *PathSpec, op *OpSpec) {
	setOp(&p.head, op, http.MethodHead, false)
}

func PATCH(p *PathSpec, op *OpSpec) {
	setOp(&p.ptch, op, http.MethodPatch, true)
}

func TRACE(p *PathSpec, op *OpSpec) {
	setOp(&p.trac, op, http.MethodTrace, false)
}

type (
	PathSpecs []*PathSpec

	PathSpec struct {
		name string
		parm bool
		get  *OpSpec
		put  *OpSpec
		post *OpSpec
		del  *OpSpec
		opts *OpSpec
		head *OpSpec
		ptch *OpSpec
		trac *OpSpec
		sec  spec.SecurityRequirements
		tags []string
		sub  PathSpecs
	}
)

func (p *PathSpecs) hasParam() bool {
	if len(*p) < 1 {
		return false
	}

	return (*p)[0].parm
}

func (p *PathSpecs) hasName(name string) bool {
	for _, itm := range *p {
		if itm.name == name {
			return true
		}
	}

	return false
}

func (p *PathSpecs) add(item *PathSpec) {
	if item.parm {
		if len(*p) > 0 {
			panic(fmt.Errorf("parameter path must be only path"))
		}
	} else {
		if p.hasParam() {
			panic(fmt.Errorf("parameter path already defined"))
		}

		if p.hasName(item.name) {
			panic(fmt.Errorf("path '%s' already declared", item.name))
		}
	}

	*p = append(*p, item)
}

func (p PathSpecs) generate() []golang.Stmt {
	res := make([]golang.Stmt, 0)

	for _, c := range p {
		res = append(res, c.generate()...)
	}

	return res
}

func (p PathSpecs) spec() spec.NamedPaths {
	res := make(spec.NamedPaths)

	for _, itm := range p {
		itm.spec(res, "", nil, nil)
	}

	return res
}

func (p PathSpec) generate() []golang.Stmt {
	return make([]golang.Stmt, 0)
}

func (p PathSpec) spec(res spec.NamedPaths, path string, sec spec.SecurityRequirements, tags []string) {
	if p.sec != nil {
		sec = p.sec
	}

	if p.tags != nil {
		tags = p.tags
	}

	name := p.name
	if p.parm {
		name = "{" + name + "}"
	}
	if path != "" {
		path += "/" + name
	} else {
		path = name
	}

	if p.get != nil || p.put != nil || p.post != nil || p.del != nil || p.opts != nil || p.head != nil || p.ptch != nil || p.trac != nil {
		res[path] = spec.PathItem{
			GET:     p.get.spec(),
			PUT:     p.put.spec(),
			POST:    p.post.spec(),
			DELETE:  p.del.spec(),
			OPTIONS: p.opts.spec(),
			HEAD:    p.head.spec(),
			PATCH:   p.ptch.spec(),
			TRACE:   p.trac.spec(),
		}
	}

	for _, itm := range p.sub {
		itm.spec(res, path, sec, tags)
	}
}

func setOp(dst **OpSpec, op *OpSpec, meth string, hasBody bool) {
	if *dst != nil {
		panic(fmt.Errorf("path already defines an operation for %s", meth))
	}

	if !hasBody && op.body != nil {
		panic(fmt.Errorf("operation has a request body though method '%s' supports none", meth))
	}

	*dst = op
}
