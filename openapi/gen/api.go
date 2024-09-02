package gen

import (
	"os"
	"path/filepath"

	"github.com/trwk76/goweb/openapi/spec"
)

func New(path string, info spec.Info, sec spec.SecurityRequirements, typeKeyFunc SchemaTypeKeyFunc) API {
	return API{
		path:   path,
		info:   info,
		sec:    sec,
		paths:  make(spec.Paths),
		sch:    newSchemas(typeKeyFunc),
		resps:  make(Resps),
		rbds:   make(ReqBodies),
		params: make(Params),
		hdrs:   make(Headers),
		secs:   make(Securities),
		tags:   make(Tags, 0),
	}
}

func (a *API) Schemas() *Schemas {
	return &a.sch
}

func (a *API) Resps() *Resps {
	return &a.resps
}

func (a *API) ReqBodies() *ReqBodies {
	return &a.rbds
}

func (a *API) Params() *Params {
	return &a.params
}

func (a *API) Headers() *Headers {
	return &a.hdrs
}

func (a *API) Securities() *Securities {
	return &a.secs
}

func (a *API) Tags() *Tags {
	return &a.tags
}

func (a *API) Spec() spec.OpenAPI {
	return spec.OpenAPI{
		OpenAPI: spec.Version,
		Info:    a.info,
		Servers: []spec.Server{{URL: a.path, Description: "Current server."}},
		Components: &spec.Components{
			Schemas:         a.sch.spec(),
			Responses:       a.resps.spec(),
			RequestBodies:   a.rbds.spec(),
			Parameters:      a.params.spec(),
			Headers:         a.hdrs.spec(),
			SecuritySchemes: a.secs.spec(),
		},
		Security: a.sec,
		Tags:     a.tags,
	}
}

func WriteSpecTo(dir string, s spec.OpenAPI) {
	if err := os.MkdirAll(dir, os.FileMode(0777)); err != nil {
		panic(err)
	}

	if err := os.WriteFile(filepath.Join(dir, "openapi.json"), s.JSON(), os.FileMode(0644)); err != nil {
		panic(err)
	}

	if err := os.WriteFile(filepath.Join(dir, "openapi.yaml"), s.YAML(), os.FileMode(0644)); err != nil {
		panic(err)
	}
}

type (
	API struct {
		path   string
		info   spec.Info
		sec    spec.SecurityRequirements
		paths  spec.Paths
		sch    Schemas
		resps  Resps
		rbds   ReqBodies
		params Params
		hdrs   Headers
		secs   Securities
		tags   Tags
	}

	SetupFunc[T any] func(a *API, item *T)
)
