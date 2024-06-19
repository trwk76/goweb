package gen

import (
	"fmt"
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
		meths: make(map[string]*Operation),
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
		meths: make(map[string]*Operation),
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
		meths:    make(map[string]*Operation),
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
		meths:    make(map[string]*Operation),
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

type (
	Path struct {
		paths
		api   *API
		par   *Path
		path  string
		pref  spec.ParameterOrRef
		meths map[string]*Operation

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
