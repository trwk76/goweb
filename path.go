package web

import (
	"fmt"
	"go/token"
	"strings"
)

func (p *path) Path() string {
	res := ""

	if p.par != nil {
		res = p.par.Path()
	}

	return res + "/" + p.name
}

func (p *path) Named(name string, errHandler ErrorHandler) Path {
	if strings.Contains(name, "/") {
		panic(fmt.Errorf("path name must not contain '/'"))
	} else if len(name) < 1 {
		panic(fmt.Errorf("path name must not be empty"))
	}

	if p.param != nil {
		panic(fmt.Errorf("path '%s' cannot accept a named subpath since it already declares a parameter as subpath", p.Path()))
	} else if _, fnd := p.named[name]; fnd {
		panic(fmt.Errorf("path '%s' already declares a subpath named '%s'", p.Path(), name))
	}

	res := &path{
		paths: paths{
			named: make(map[string]*path),
		},
		par:   p,
		name:  name,
		meths: make(map[string]Handler),
		errh:  errHandler,
	}

	p.named[name] = res
	return res
}

func (p *path) Param(name string, isPath bool, errHandler ErrorHandler) Path {
	if !token.IsIdentifier(name) {
		panic(fmt.Errorf("'%s' is not a valid identifier", name))
	} else if p.hasParam(name) {
		panic(fmt.Errorf("path '%s' already declares parameter '%s'", p.par.Path(), name))
	}

	if len(p.named) > 0 {
		panic(fmt.Errorf("path '%s' cannot accept a parameter subpath since it already has one or more named subpaths", p.Path()))
	} else if p.param != nil {
		panic(fmt.Errorf("path '%s' already defines a parameter subpath", p.Path()))
	}

	if isPath {
		name = "*" + name
	} else {
		name = ":" + name
	}

	res := &path{
		paths: paths{
			named: make(map[string]*path),
		},
		par:   p,
		name:  name,
		meths: make(map[string]Handler),
		errh:  errHandler,
	}

	p.param = res
	return res
}

func (p *path) Handle(meth string, hdl Handler) {
	p.meths[meth] = hdl
}

type (
	Path interface {
		Path() string
		Named(name string, errHandler ErrorHandler) Path
		Param(name string, path bool, errHandler ErrorHandler) Path
	}

	path struct {
		paths
		par   *path
		name  string
		meths map[string]Handler
		errh  ErrorHandler
	}

	paths struct {
		named map[string]*path
		param *path
	}
)

func (p *path) hasParam(name string) bool {
	if p.name[0] == ':' || p.name[0] == '*' {
		if p.name[1:] == name {
			return true
		}
	}

	if p.par != nil {
		return p.par.hasParam(name)
	}

	return false
}

func (p paths) child(ctx *Context, path []string) (*path, []string) {
	if par := p.param; par != nil {
		if par.name[0] == ':' || par.name[0] == '*' {
			ctx.params[par.name[1:]] = strings.Join(path, "/")
			path = nil
		} else {
			ctx.params[par.name[1:]] = path[0]
			path = path[1:]
		}

		return par, path
	}

	res, ok := p.named[path[0]]
	if !ok {
		return nil, path
	}

	return res, path[1:]
}

var _ Path = (*path)(nil)
