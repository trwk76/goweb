package gen

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/trwk76/goweb/openapi/spec"
)

func (a *API) NamedPath(name string, f func(p *Path)) *Path {
	if len(name) < 1 || strings.ContainsAny(name, "/") {
		panic(fmt.Errorf("path name '%s' is invalid", name))
	}

	if _, fnd := a.paths.named[name]; fnd {
		panic(fmt.Errorf("a subpath named '%s' is already declared", name))
	} else if a.paths.param != nil {
		panic(fmt.Errorf("path cannot have a named subpath since it declares a parameter"))
	}

	res := &Path{
		paths: paths{
			named: make(map[string]*Path),
		},
		api:   a,
		par:   nil,
		path:  "/" + name,
		meths: make(map[string]spec.Operation),
	}

	a.paths.named[name] = res

	if f != nil {
		f(res)
	}

	return res
}

func (a *API) ParamPath(param spec.ParameterOrRef, f func(p *Path)) *Path {
	pspec := a.Param(param)

	if pspec.Name == "" {
		panic(fmt.Errorf("parameter name must not be empty"))
	} else if pspec.In != spec.ParameterPath {
		panic(fmt.Errorf("parameter 'in' is not 'path'"))
	}

	if len(a.paths.named) > 0 {
		panic(fmt.Errorf("path cannot have a parameter since it already declares named subpaths"))
	} else if a.paths.param != nil {
		panic(fmt.Errorf("path already declares a parameter"))
	}

	res := &Path{
		paths: paths{
			named: make(map[string]*Path),
		},
		api:   a,
		par:   nil,
		path:  "/:" + pspec.Name,
		pref:  param,
		meths: make(map[string]spec.Operation),
	}

	a.paths.param = res

	if f != nil {
		f(res)
	}

	return res
}

func (p *Path) NamedPath(name string, f func(p *Path)) *Path {
	if len(name) < 1 || strings.ContainsAny(name, "/") {
		panic(fmt.Errorf("path name '%s' is invalid", name))
	}

	if _, fnd := p.paths.named[name]; fnd {
		panic(fmt.Errorf("a subpath named '%s' is already declared", name))
	} else if p.paths.param != nil {
		panic(fmt.Errorf("path cannot have a named subpath since it declares a parameter"))
	}

	res := &Path{
		paths: paths{
			named: make(map[string]*Path),
		},
		api:      p.api,
		par:      p,
		path:     p.path + "/" + name,
		meths:    make(map[string]spec.Operation),
		OpPrefix: p.OpPrefix,
		Security: p.Security,
		Tags:     p.Tags,
	}

	p.paths.named[name] = res

	if f != nil {
		f(res)
	}

	return res
}

func (p *Path) ParamPath(param spec.ParameterOrRef, f func(p *Path)) *Path {
	pspec := p.api.Param(param)

	if pspec.Name == "" {
		panic(fmt.Errorf("parameter name must not be empty"))
	} else if pspec.In != spec.ParameterPath {
		panic(fmt.Errorf("parameter 'in' is not 'path'"))
	}

	if len(p.paths.named) > 0 {
		panic(fmt.Errorf("path cannot have a parameter since it already declares named subpaths"))
	} else if p.paths.param != nil {
		panic(fmt.Errorf("path already declares a parameter"))
	} else if p.par.hasParam(pspec.Name) {
		panic(fmt.Errorf("parameter '%s' already declared in parent path", pspec.Name))
	}

	res := &Path{
		paths: paths{
			named: make(map[string]*Path),
		},
		api:      p.api,
		par:      p,
		path:     p.par.path + "/:" + pspec.Name,
		pref:     param,
		meths:    make(map[string]spec.Operation),
		OpPrefix: p.OpPrefix,
		Security: p.Security,
		Tags:     p.Tags,
	}

	p.paths.param = res

	if f != nil {
		f(res)
	}

	return res
}

func (p *Path) GET(op spec.Operation) {
	p.setMeth(http.MethodGet, op, false)
}

func (p *Path) PUT(op spec.Operation) {
	p.setMeth(http.MethodPut, op, true)
}

func (p *Path) POST(op spec.Operation) {
	p.setMeth(http.MethodPost, op, true)
}

func (p *Path) DELETE(op spec.Operation) {
	p.setMeth(http.MethodDelete, op, false)
}

func (p *Path) OPTIONS(op spec.Operation) {
	p.setMeth(http.MethodOptions, op, false)
}

func (p *Path) HEAD(op spec.Operation) {
	p.setMeth(http.MethodHead, op, false)
}

func (p *Path) PATCH(op spec.Operation) {
	p.setMeth(http.MethodPatch, op, true)
}

func (p *Path) TRACE(op spec.Operation) {
	p.setMeth(http.MethodTrace, op, false)
}

type (
	Path struct {
		paths
		api   *API
		par   *Path
		path  string
		pref  spec.ParameterOrRef
		meths map[string]spec.Operation

		OpPrefix string
		Security spec.SecurityRequirements
		Tags     Tags
	}

	paths struct {
		named map[string]*Path
		param *Path
	}
)

func (p *Path) hasParam(name string) bool {
	if pspec := p.api.Param(p.pref); pspec.In == spec.ParameterPath && pspec.Name == name {
		return true
	}

	if p.par != nil {
		return p.par.hasParam(name)
	}

	return false
}

func (p *Path) params() spec.ParameterOrRefs {
	var res spec.ParameterOrRefs

	if p.par != nil {
		res = p.par.params()
	}

	if p.pref.Ref.Ref != "" || p.pref.Item.In == spec.ParameterPath {
		res = append(res, p.pref)
	}

	return res
}

func (p *Path) setMeth(meth string, op spec.Operation, acceptBody bool) {
	if !acceptBody && op.RequestBody != nil {
		panic(fmt.Errorf("http method %s does not allow a request body", meth))
	}

	if p.OpPrefix != "" {
		op.OperationID = p.OpPrefix + op.OperationID
	}

	if len(op.Security) < 1 {
		op.Security = p.Security
	}

	op.Parameters = append(p.params(), op.Parameters...)
	op.Tags = append(op.Tags, p.Tags.spec()...)
	p.meths[meth] = op
}
