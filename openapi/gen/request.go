package gen

import (
	"fmt"
	"reflect"

	"github.com/trwk76/goweb/openapi/marshal"
	"github.com/trwk76/goweb/openapi/spec"
)

func RequestFor[T any](a *API, op *spec.Operation, paramSetup SetupFunc[spec.Parameter], bodySetup SetupFunc[spec.RequestBody]) {
	RequestOf(a, reflect.TypeFor[T](), op, paramSetup, bodySetup)
}

func RequestOf(a *API, t reflect.Type, op *spec.Operation, paramSetup SetupFunc[spec.Parameter], bodySetup SetupFunc[spec.RequestBody]) {
	decl := make(map[string]string)

	marshal.ReflectRequest(
		t,
		func(fld reflect.StructField, in spec.ParameterIn, name string) {
			key := fmt.Sprintf("%s@%s", in, name)

			if _, ex := decl[key]; ex {
				panic(fmt.Errorf("type '%s', field '%s': redefines parameter '%s'", t.String(), fld.Name, key))
			}

			decl[key] = key

			fldt := fld.Type
			opt := false
			if fldt.Kind() == reflect.Pointer {
				fldt = fldt.Elem()
				opt = true
			}

			if _, fnd := FindParam(a, op.Parameters, in, name); !fnd {
				op.Parameters = append(op.Parameters, ParamOf(a, fldt, name, in, func(a *API, item *spec.Parameter) { item.Required = !opt }))
			}
		},
		func(fld reflect.StructField) {
			if _, ex := decl["body"]; ex {
				panic(fmt.Errorf("type '%s', field '%s': redefines request body", t.String(), fld.Name))
			}

			decl["body"] = "body"

			fldt := fld.Type
			opt := false
			if fldt.Kind() == reflect.Pointer {
				fldt = fldt.Elem()
				opt = true
			}

			if op.RequestBody == nil {
				body := ReqBodyOf(a, fldt, func(a *API, item *spec.RequestBody) { item.Required = !opt })
				op.RequestBody = &body
			} else {
				_, ok := a.ReqBodies().ReqBody(*op.RequestBody)
				if !ok {
					panic(fmt.Errorf("type '%s', field '%s': body reference invalid", t.String(), fld.Name))
				}
			}
		},
	)
}

func FindParam(a *API, coll spec.ParameterOrRefs, in spec.ParameterIn, name string) (spec.Parameter, bool) {
	for _, item := range coll {
		if s, ok := a.Params().Param(item); ok && s.In == in && s.Name == name {
			return s, true
		}
	}

	return spec.Parameter{}, false
}
