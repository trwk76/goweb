package openapi

import "github.com/trwk76/goweb/openapi/spec"

func New(path string, info spec.Info) Builder {
	return Builder{
		path: path,
		info: info,
	}
}

type (
	Builder struct {
		Auths Auths
		path  string
		info  spec.Info
	}
)
