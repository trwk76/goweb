package web

import (
	"fmt"
	"strings"
)

type (
	Path struct {
		par  *Path
		nam  string
		typ  PathType
		errh ErrorHandler
		mth  meths
		sub  paths
	}

	PathType string

	Handler func(ctx *Context) Response
)

const (
	PathNamed   PathType = "named"
	PathParam   PathType = "param"
	PathSubPath PathType = "subpath"
)

func (p *Path) Name() string {
	switch p.typ {
	case PathParam:
		return ":" + p.nam
	case PathSubPath:
		return "*" + p.nam
	}

	return p.nam
}

func (p *Path) BasePath() string {
	if p.par == nil {
		return "/"
	} else if p.par.par == nil {
		return "/" + p.Name()
	}

	return p.par.BasePath() + "/" + p.Name()
}

func (p *Path) Path(path string) *Path {
	for _, itm := range parsePath(path) {
		switch itm {
		case ".":
		case "..":
			if p.par != nil {
				p = p.par
			}
		default:
			p = p.Child(itm)
		}
	}

	return p
}

func (p *Path) ErrorHandler(errh ErrorHandler) {
	p.errh = errh
}

func (p *Path) Handle(mth string, sec SecurityRequirements, handler Handler) {
	if p.mth == nil {
		p.mth = make(meths)
	}

	mth = strings.ToUpper(mth)

	if _, found := p.mth[mth]; found {
		panic(fmt.Errorf("path '%s' already has a handler for method '%s'", p.BasePath(), mth))
	}

	p.mth[mth] = meth{
		sec: sec,
		hdl: handler,
	}
}

func (p *Path) Child(name string) *Path {
	var param *Path
	var named []*Path

	if p.typ == PathSubPath {
		panic(fmt.Errorf("path parameter cannot bear subpaths"))
	}

	typ := PathNamed
	if strings.HasPrefix(name, ":") {
		typ = PathParam
		name = name[1:]
	} else if strings.HasPrefix(name, "*") {
		typ = PathSubPath
		name = name[1:]
	}

	for _, itm := range p.sub {
		if itm.typ == PathNamed {
			if typ == PathNamed {
				if name == itm.nam {
					return itm
				}
			}

			named = append(named, itm)
		} else {
			if typ != PathNamed {
				if typ == itm.typ && name == itm.nam {
					return itm
				}
			}

			param = itm
		}
	}

	if typ != PathNamed && param != nil {
		panic(fmt.Errorf("path '%s' already has a parameter subpath", p.BasePath()))
	}

	res := &Path{
		par: p,
		nam: name,
		typ: typ,
	}

	if typ != PathNamed {
		param = res
	} else {
		named = append(named, res)
	}

	p.sub = named
	if param != nil {
		p.sub = append(p.sub, param)
	}

	return res
}

func (p *Path) child(c *Context, path []string) (*Path, []string) {
	for _, sub := range p.sub {
		switch sub.typ {
		case PathNamed:
			if sub.nam == path[0] {
				return sub, path[1:]
			}
		case PathParam:
			c.parm[sub.nam] = path[0]
			return sub, path[1:]
		case PathSubPath:
			c.parm[sub.nam] = strings.Join(path, "/")
			return sub, nil
		}
	}

	return nil, path
}

type (
	paths []*Path

	meths map[string]meth

	meth struct {
		sec SecurityRequirements
		hdl Handler
	}
)

func parsePath(path string) []string {
	for strings.Contains(path, "//") {
		path = strings.ReplaceAll(path, "//", "/")
	}

	path = strings.TrimPrefix(strings.TrimSuffix(path, "/"), "/")
	if path == "" {
		return nil
	}

	return strings.Split(path, "/")
}
