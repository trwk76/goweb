package openapi

import "github.com/trwk76/goweb/openapi/spec"

func New(path string, info spec.Info, mtyp MediaTypes) Builder {
	return Builder{
		path: path,
		info: info,
		mtyp: mtyp,
	}
}

type (
	Builder struct {
		Auths Auths
		path  string
		mtyp  MediaTypes
		info  spec.Info
	}
)
