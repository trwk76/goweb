package web

import "strings"

type (
	path struct {
		paths
		par   *path
		meths map[string]Handler
		errh  ErrorHandler
	}

	paramPath struct {
		path
		paramName string
		paramPath bool
	}

	paths struct {
		named map[string]*path
		param *paramPath
	}
)

func (p paths) child(ctx *Context, path []string) (*path, []string) {
	if par := p.param; par != nil {
		if par.paramPath {
			ctx.params[par.paramName] = strings.Join(path, "/")
			path = nil
		} else {
			ctx.params[par.paramName] = path[0]
			path = path[1:]
		}

		return &par.param.path, path
	}

	res, ok := p.named[path[0]]
	if !ok {
		return nil, path
	}

	return res, path[1:]
}
