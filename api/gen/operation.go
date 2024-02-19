package gen

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/trwk76/gocode/golang"
	web "github.com/trwk76/goweb"
	"github.com/trwk76/goweb/api"
	"github.com/trwk76/goweb/api/spec"
)

func Op[REQ any, RES any, HDL OpHandler[REQ, RES]](s *APISpec, handler HDL, info OpInfo) *OpSpec {
	funcName := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()

	rest, ress, _ := buildResHandler[RES](s, info)
	reqt, reqs := buildReqHandler[REQ](s, info)

	lbd := golang.FuncExpr{
		Params: golang.Params{
			{
				Name: "ctx",
				Type: golang.PtrType{
					Target: golang.NamedType{Name: "Context", Pkg: webPkg},
				},
			},
		},
		Return: golang.Return{
			{Type: golang.NamedType{Name: "Response", Pkg: webPkg}},
		},
		Body: golang.BlockStmt{
			golang.VarDecls{
				{
					Name: "req",
					Type: s.file.Type(reqt),
				},
				{
					Name: "res",
					Type: s.file.Type(rest),
				},
				{
					Name: "perr",
					Type: s.file.NamedType("ParseError", apiPkg, nil),
				},
			},
		},
	}

	lbd.Body = append(lbd.Body, reqs...)

	lbd.Body = append(lbd.Body, golang.ExprStmt{
		Expr: golang.CallExpr{
			Func: golang.SymbolExpr{Name: funcName},
			Args: golang.Args{
				golang.SymbolExpr{Name: "ctx"},
				golang.SymbolExpr{Name: "req"},
				golang.UnaryExpr{Op: golang.AddrOf, Expr: golang.SymbolExpr{Name: "res"}},
			},
		},
	})

	lbd.Body = append(lbd.Body, ress)

	return &OpSpec{
		id:   info.ID,
		summ: info.Summary,
		desc: info.Description,
		depr: info.Deprecated,
		parm: info.Params,
		body: info.ReqBody,
		resp: info.Resps,
		sec:  info.Security,
		tags: info.Tags,
		fnc:  lbd,
	}
}

type (
	OpInfo struct {
		ID          string
		Summary     string
		Description string
		Deprecated  bool
		Params      ParamsSpec
		ReqBody     ReqBodySpec
		Resps       RespsSpec
		Security    spec.SecurityRequirements
		Tags        []string
	}

	OpHandler[REQ any, RES any] func(ctx *web.Context, req REQ, res *RES)

	OpSpec struct {
		id   string
		summ string
		desc string
		depr bool
		parm ParamsSpec
		body ReqBodySpec
		resp RespsSpec
		sec  spec.SecurityRequirements
		tags []string
		fnc  golang.FuncExpr
	}

	ParamsSpec []ParamSpec
	RespsSpec  map[int]RespSpec
)

func (s *OpSpec) spec() *spec.Operation {
	if s == nil {
		return nil
	}

	var body *spec.RequestBody

	if s.body != nil {
		tmp := s.body.spec()
		body = &tmp
	}

	return &spec.Operation{
		OperationID: s.id,
		Summary:     s.summ,
		Description: s.desc,
		Deprecated:  s.depr,
		Parameters:  s.parm.spec(),
		RequestBody: body,
		Responses:   s.resp.spec(),
		Security:    s.sec,
		Tags:        s.tags,
	}
}

func (p ParamsSpec) spec() spec.Parameters {
	res := make(spec.Parameters, len(p))

	for idx, itm := range p {
		res[idx] = itm.spec()
	}

	return res
}

func (p ParamsSpec) find(name string, in spec.ParameterIn) ParamSpec {
	for _, itm := range p {
		impl := itm.impl()
		if impl.name == name && impl.in == in {
			return itm
		}
	}

	return nil
}

func (p RespsSpec) spec() spec.Responses {
	res := make(spec.Responses)

	for sta, itm := range p {
		res[fmt.Sprintf("%03d", sta)] = itm.spec()
	}

	return res
}

func (r RespsSpec) find(status int) RespSpec {
	if res, ok := r[status]; ok {
		return res
	}

	if res, ok := r[0]; ok {
		return res
	}

	panic(fmt.Errorf("no response specified for status code '%03d'", status))
}

func buildReqHandler[REQ any](s *APISpec, info OpInfo) (reflect.Type, []golang.Stmt) {
	var d REQ

	t := reflect.TypeOf(d)
	reqs := make([]golang.Stmt, 0)
	prms := make(map[string]bool)
	doc := make(map[string]bool)

	for _, p := range info.Params {
		impl := p.impl()
		doc[impl.name+"@"+string(impl.in)] = true
	}

	if info.ReqBody != nil {
		doc["body"] = true
	}

	for _, fld := range reflect.VisibleFields(t) {
		tag := fld.Tag.Get("api")

		if tag != "" {
			if _, found := prms[tag]; found {
				panic(fmt.Errorf("'%s' declared more than once", tag))
			}

			prms[tag] = true
			name := ""
			in := ""

			parseFunc := "Parse"
			t := fld.Type
			if t.Kind() == reflect.Pointer {
				parseFunc = "ParsePtr"
				t = t.Elem()
			}

			if tag == "body" {
				delete(doc, "body")
				in = "InBody"

				if info.ReqBody == nil {
					panic(fmt.Errorf("request body is not documented"))
				}

				if t != info.ReqBody.impl().sch.impl().typ {
					panic(fmt.Errorf("type mismatch"))
				}
			} else {
				delete(doc, tag)

				namv, inv, ok := strings.Cut(tag, "@")
				if !ok {
					panic(fmt.Errorf("invalid api request tag '%s'", tag))
				}

				d := info.Params.find(namv, spec.ParameterIn(inv))
				if d == nil {
					panic(fmt.Errorf("parameter '%s' is not documented", tag))
				}

				if t != d.impl().sch.impl().typ {
					panic(fmt.Errorf("type mismatch"))
				}

				name = namv
				switch inv {
				case string(spec.ParameterInCookie):
					in = "InCookie"
				case string(spec.ParameterInHeader):
					in = "InHeader"
				case string(spec.ParameterInPath):
					in = "InPath"
				case string(spec.ParameterInQuery):
					in = "InQuery"
				}
			}

			reqs = append(reqs, golang.ExprStmt{
				Expr: golang.CallExpr{
					Func: s.file.Symbol(parseFunc, apiPkg),
					Args: golang.Args{
						golang.UnaryExpr{
							Op: golang.AddrOf,
							Expr: golang.MemberExpr{
								Expr: golang.SymbolExpr{Name: "req"},
								Name: fld.Name,
							},
						},
						golang.SymbolExpr{Name: "ctx"},
						golang.StringExpr(name),
						s.file.Symbol(in, apiPkg),
						golang.UnaryExpr{
							Op:   golang.AddrOf,
							Expr: golang.SymbolExpr{Name: "perr"},
						},
					},
				},
			})
		}
	}

	return t, reqs
}

func buildResHandler[RES any](s *APISpec, info OpInfo) (reflect.Type, golang.Stmt, map[int]string) {
	var d RES
	var res golang.ElseStmt

	t := reflect.TypeOf(d)
	sf := make(map[int]string)

	for _, fld := range reflect.VisibleFields(t) {
		tag := fld.Tag.Get("api")

		if tag != "" {
			var status int

			if _, err := fmt.Sscanf(tag, "%03d", &status); err != nil && (status < 100 || status > 599) {
				panic(fmt.Errorf("response status tag '%s' is invalid", tag))
			}

			if fld.Type.Kind() != reflect.Pointer {
				panic(fmt.Errorf("response field type must be pointer"))
			}

			resp := info.Resps.find(status)
			if resp == nil {
				panic(fmt.Errorf("response '%03d' is not documented", status))
			}

			rimpl := resp.impl()
			rtyp := rimpl.sch.impl().typ

			if fld.Type.Elem() != rtyp {
				panic(fmt.Errorf("response '%03d' types mismatch: %s when %s expected", status, fld.Type.Elem().String(), rtyp.String()))
			}

			sf[status] = fld.Name

			res = golang.IfStmt{
				Cond: golang.BinaryExpr{
					Op: golang.NotEqual,
					Lhs: golang.MemberExpr{
						Expr: golang.SymbolExpr{Name: "res"},
						Name: fld.Name,
					},
					Rhs: golang.Nil,
				},
				Then: golang.BlockStmt{
					golang.ReturnStmt{
						Value: golang.CallExpr{
							Func: golang.SymbolExpr{
								Name: "JSONResponse",
								Pkg:  "web",
							},
						},
					},
				},
				Else: res,
			}
		}
	}

	return t, res, sf
}

var (
	dw     web.Response
	da     api.In
	webPkg string = reflect.TypeOf(dw).PkgPath()
	apiPkg string = reflect.TypeOf(da).PkgPath()
)
