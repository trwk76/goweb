package openapi

import (
	"log"
	"net/http"

	"github.com/trwk76/goweb"
)

func GET[REQ any, RES any](p *PathBuilder, opts OperationOpts, handler OperationHandler[REQ, RES]) {
	setOperation[REQ, RES](p, http.MethodGet, false, opts, handler)
}

func PUT[REQ any, RES any](p *PathBuilder, opts OperationOpts, handler OperationHandler[REQ, RES]) {
	setOperation[REQ, RES](p, http.MethodPut, true, opts, handler)
}

func POST[REQ any, RES any](p *PathBuilder, opts OperationOpts, handler OperationHandler[REQ, RES]) {
	setOperation[REQ, RES](p, http.MethodPost, true, opts, handler)
}

func DELETE[REQ any, RES any](p *PathBuilder, opts OperationOpts, handler OperationHandler[REQ, RES]) {
	setOperation[REQ, RES](p, http.MethodDelete, true, opts, handler)
}

func OPTIONS[REQ any, RES any](p *PathBuilder, opts OperationOpts, handler OperationHandler[REQ, RES]) {
	setOperation[REQ, RES](p, http.MethodOptions, false, opts, handler)
}

func HEAD[REQ any, RES any](p *PathBuilder, opts OperationOpts, handler OperationHandler[REQ, RES]) {
	setOperation[REQ, RES](p, http.MethodHead, false, opts, handler)
}

func PATCH[REQ any, RES any](p *PathBuilder, opts OperationOpts, handler OperationHandler[REQ, RES]) {
	setOperation[REQ, RES](p, http.MethodPatch, true, opts, handler)
}

func TRACE[REQ any, RES any](p *PathBuilder, opts OperationOpts, handler OperationHandler[REQ, RES]) {
	setOperation[REQ, RES](p, http.MethodTrace, false, opts, handler)
}

type (
	OperationHandler[REQ any, RES any] func(ctx *goweb.Context, req REQ, res *RES)

	OperationOpts struct {
		ID          string
		Summary     string
		Description string
		Tags        []string
		Deprecated  bool
		MediaTypes  MediaTypes
		Security    SecurityReqs
		Parameters  ParamsOpts
		ReqBody     *ReqBodyOpt
		Responses   RespsOpts
	}

	ParamsOpts struct {
		Cookie map[string]ParamOpt
		Header map[string]ParamOpt
		Path   map[string]ParamOpt
		Query  map[string]ParamOpt
	}

	ParamOpt struct {
		Description     string
		Required        bool
		Deprecated      bool
		AllowEmptyValue bool
		Style           string
		Explode         bool
		AllowReserved   bool
		Examples        ExamplesOpt
	}

	ReqBodyOpt struct {
		Description string
		Required    bool
		Examples    ExamplesOpt
	}

	RespsOpts map[int]RespOpts

	RespOpts struct {
		Description string
		Headers     map[string]ParamOpt
		Examples    ExamplesOpt
	}

	ExamplesOpt map[string]ExampleOpt

	ExampleOpt struct {
		Summary     string
		Description string
		Value       any
	}

	OperationBuilder struct {
		opts OperationOpts
	}
)

func setOperation[REQ any, RES any](p *PathBuilder, meth string, allowBody bool, opts OperationOpts, handler OperationHandler[REQ, RES]) {
	if p.ops[meth] != nil {
		log.Fatalf("path '%s' already defines an operation for method '%s'", p.pth, meth)
	}

	opts.Tags = append(p.opts.Tags, opts.Tags...)

	if opts.MediaTypes == nil {
		opts.MediaTypes = p.opts.MediaTypes
	}

	if opts.Security == nil {
		opts.Security = p.opts.Security
	}

	res := &OperationBuilder{
		opts: opts,
	}

	p.ops[meth] = res
}
