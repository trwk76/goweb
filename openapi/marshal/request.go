package marshal

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/trwk76/goweb/openapi/spec"
)

func UnmarshalRequest[T any](dest *T, req *http.Request) error {
	var errs PathErrors

	unmarshalRequest(reflect.Indirect(reflect.ValueOf(dest)), req, &errs)

	if len(errs) > 0 {
		return errs
	}

	return nil
}

func ReflectRequest(t reflect.Type, base func(fld reflect.StructField), param func(fld reflect.StructField, in spec.ParameterIn, name string), body func(fld reflect.StructField)) {
	if t.Kind() != reflect.Struct {
		panic(fmt.Errorf("type '%s' is not a structure", t.String()))
	}

	for i := 0; i < t.NumField(); i++ {
		fld := t.Field(i)
		tag := fld.Tag.Get("openapi")

		if tag != "-" {
			if fld.Anonymous && tag == "" {
				// Is a base structure
				base(fld)
			} else {
				if tag == "body" {
					body(fld)
				} else {
					in, name, ok := strings.Cut(tag, "@")
					if !ok {
						panic(fmt.Errorf("type '%s', field '%s': openapi tag is invalid", t.String(), fld.Name))
					}

					switch spec.ParameterIn(in) {
					case spec.ParameterCookie, spec.ParameterHeader, spec.ParameterPath, spec.ParameterQuery:
					default:
						panic(fmt.Errorf("type '%s', field '%s': openapi tag is invalid: invalid paramter location '%s'", t.String(), fld.Name, in))
					}

					param(fld, spec.ParameterIn(in), name)
				}
			}
		}
	}
}

func unmarshalRequest(dest reflect.Value, req *http.Request, errs *PathErrors) {
	ReflectRequest(
		dest.Type(),
		func(fld reflect.StructField) {
			unmarshalRequest(dest.FieldByIndex(fld.Index), req, errs)
		},
		func(fld reflect.StructField, in spec.ParameterIn, name string) {

		},
		func(fld reflect.StructField) {

		},
	)
}
