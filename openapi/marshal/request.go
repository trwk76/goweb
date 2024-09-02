package marshal

import (
	"encoding/json"
	"fmt"
	"io"
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

func ReflectRequest(t reflect.Type, param func(fld reflect.StructField, in spec.ParameterIn, name string), body func(fld reflect.StructField)) {
	if t.Kind() != reflect.Struct {
		panic(fmt.Errorf("type '%s' is not a structure", t.String()))
	}

	for _, fld := range reflect.VisibleFields(t) {
		tag := fld.Tag.Get("openapi")

		if tag != "-" {
			if fld.Anonymous && tag == "" {
				// Is a base structure
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
		func(fld reflect.StructField, in spec.ParameterIn, name string) {
			fval := dest.FieldByIndex(fld.Index)
			got := false
			opt := false
			raw := ""

			if fval.Kind() == reflect.Pointer {
				opt = true
				fval = reflect.Indirect(fval)
			}

			switch in {
			case spec.ParameterCookie:
				if cook, err := req.Cookie(name); cook != nil && err == nil {
					raw = cook.Value
					got = true
				}
			case spec.ParameterHeader:
				if _, ok := req.Header[name]; ok {
					raw = req.Header.Get(name)
					got = true
				}
			case spec.ParameterPath:
				raw = req.PathValue(name)
				got = (raw != "")
			case spec.ParameterQuery:
				if _, ok := req.URL.Query()[name]; ok {
					raw = req.URL.Query().Get(name)
					got = true
				}
			}

			if got {
				if err := UnmarshalText(fval.Addr().Interface(), []byte(raw)); err != nil {
					*errs = append(
						*errs,
						PathError{Path: Path{PathMember(in), PathMember(name)}, Err: err},
					)
				}
			} else {
				if !opt {
					*errs = append(
						*errs,
						PathError{Path: Path{PathMember(in), PathMember(name)}, Err: fmt.Errorf("parameter requires a value")},
					)
				}
			}
		},
		func(fld reflect.StructField) {
			fval := dest.FieldByIndex(fld.Index)
			opt := false

			if fval.Kind() == reflect.Pointer {
				opt = true
				fval = reflect.Indirect(fval)
			}

			raw, err := io.ReadAll(req.Body)
			if err != nil {
				*errs = append(
					*errs,
					PathError{Path: Path{PathMember("body")}, Err: fmt.Errorf("error reading body: %s", err.Error())},
				)

				return
			}

			if len(raw) < 1 {
				if !opt {
					*errs = append(
						*errs,
						PathError{Path: Path{PathMember("body")}, Err: fmt.Errorf("request requires a body")},
					)
				} else {
					if err := json.Unmarshal(raw, fval.Addr().Interface()); err != nil {
						*errs = append(
							*errs,
							NewPathErrors(PathMember("body"), err)...,
						)
					}
				}
			}
		},
	)
}
