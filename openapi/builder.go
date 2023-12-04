package openapi

import "github.com/trwk76/goweb/openapi/spec"

func New(path string, info spec.Info, tags spec.Tags, mtypes MediaTypes, secs SecuritySchemes, sec SecurityReqs) *Builder {
	res := &Builder{
		pth:  path,
		inf:  info,
		tags: tags,
		secs: secs,
	}

	res.PathBuilder = newPathBuilder(res, "", PathOpts{
		MediaTypes: mtypes,
		Security:   sec,
	})

	return res
}

type (
	Builder struct {
		PathBuilder

		pth  string
		inf  spec.Info
		tags []spec.Tag
		secs SecuritySchemes
	}
)
