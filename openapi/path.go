package openapi

import (
	"fmt"
	"log"
)

func (b *PathBuilder) Named(name string, opts PathOpts, action func(p *PathBuilder)) *PathBuilder {
	path := fmt.Sprintf("/%s/%s", b.pth, name)

	if _, found := b.subn[name]; found {
		log.Fatalf("path '%s' already declared", path)
	}

	npb := newPathBuilder(b.api, path, b.opts.merge(opts))
	res := &npb
	b.subn[name] = res

	if action != nil {
		action(res)
	}

	return res
}

func (b *PathBuilder) Param(name string, opts PathOpts, action func(p *PathBuilder)) *PathBuilder {
	path := fmt.Sprintf("/%s/:%s", b.pth, name)

	if b.subp != nil {
		log.Fatalf("path '%s' already declares a sub parameter", b.pth)
	}

	b.subp = &paramPath{
		name: name,
		path: newPathBuilder(b.api, path, b.opts.merge(opts)),
	}

	res := &b.subp.path

	if action != nil {
		action(res)
	}

	return res
}

type (
	PathOpts struct {
		Summary     string
		Description string
		Tags        []string
		MediaTypes  MediaTypes
		Security    SecurityReqs
	}

	PathBuilder struct {
		api  *Builder
		pth  string
		opts PathOpts
		ops  map[string]*OperationBuilder
		subn map[string]*PathBuilder
		subp *paramPath
	}

	paramPath struct {
		name string
		path PathBuilder
	}
)

func newPathBuilder(api *Builder, path string, opts PathOpts) PathBuilder {
	return PathBuilder{
		api:  api,
		pth:  path,
		opts: opts,
		ops:  make(map[string]*OperationBuilder),
		subn: make(map[string]*PathBuilder),
		subp: nil,
	}
}

func (o PathOpts) merge(r PathOpts) PathOpts {
	if r.MediaTypes == nil {
		r.MediaTypes = o.MediaTypes
	}

	if r.Security == nil {
		r.Security = o.Security
	}

	return PathOpts{
		Summary:     r.Summary,
		Description: r.Description,
		Tags:        append(o.Tags, r.Tags...),
		MediaTypes:  r.MediaTypes,
		Security:    r.Security,
	}
}
