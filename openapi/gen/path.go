package gen

import (
	"fmt"
	"strings"

	"github.com/trwk76/goweb/openapi/spec"
)

func (a *API) NamedPath(name string, setup SetupFunc[PathItem]) {
	setup(a, newNamedPathItem(a, nil, name))
}

func (a *API) ParamPath(param spec.ParameterOrRef, setup SetupFunc[PathItem]) {
	setup(a, newParamPathItem(a, nil, param))
}

func (p *PathItem) NamedPath(name string, setup SetupFunc[PathItem]) {
	setup(p.api, newNamedPathItem(p.api, p.par, name))
}

func (p *PathItem) ParamPath(param spec.ParameterOrRef, setup SetupFunc[PathItem]) {
	setup(p.api, newParamPathItem(p.api, p.par, param))
}

type (
	PathItem struct {
		api   *API
		par   *PathItem
		pth   string
		param spec.ParameterOrRef
	}
)

func (p *PathItem) ensureSpec() *spec.PathItem {
	res, ok := p.api.paths[p.pth]
	if !ok {
		res = &spec.PathItem{}
		p.api.paths[p.pth] = res
	}

	return res
}

func newNamedPathItem(api *API, par *PathItem, name string) *PathItem {
	pth := "/" + name

	if par != nil {
		pth = par.par.pth + pth
	}

	return &PathItem{
		api: api,
		par: par,
		pth: pth,
	}
}

func newParamPathItem(api *API, par *PathItem, param spec.ParameterOrRef) *PathItem {
	pspc, ok := api.Params().Param(param)
	if !ok {
		panic(fmt.Errorf("parameter cannot be resolved"))
	} else if pspc.In != spec.ParameterPath {
		panic(fmt.Errorf("not a path parameter"))
	}

	pth := fmt.Sprintf("/{%s}", pspc.Name)

	if par != nil {
		if strings.Contains(par.pth, pth) {
			panic(fmt.Errorf("path already contains a parameter named '%s'", pspc.Name))
		}

		pth = par.par.pth + pth
	}

	return &PathItem{
		api:   api,
		par:   par,
		pth:   pth,
		param: param,
	}
}

func (p *PathItem) params() spec.ParameterOrRefs {
	res := spec.ParameterOrRefs{}

	if p.par != nil {
		res = p.par.params()
	}

	if p.param.Ref.Ref != "" || p.param.Item.Name != "" {
		res = append(res, p.param)
	}

	return res
}
