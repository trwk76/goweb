package marshal

import (
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
	"strings"
)

func UnmarshalStruct[T any](dest *T, raw []byte, extra *map[string]json.RawMessage) error {
	var errs PathErrors

	props := make(map[string]json.RawMessage)

	if err := json.Unmarshal(raw, &props); err != nil {
		return fmt.Errorf("json object expected")
	}

	unmarshalStruct(reflect.Indirect(reflect.ValueOf(dest)), props, &errs)

	if extra != nil {
		*extra = props
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

func unmarshalStruct(dest reflect.Value, props map[string]json.RawMessage, errs *PathErrors) {
	reflectStruct(
		dest.Type(),
		func(fld reflect.StructField) {
			unmarshalStruct(dest.FieldByIndex(fld.Index), props, errs)
		},
		func(fld reflect.StructField, name string, required bool) {
			raw, fnd := props[name]

			if !fnd {
				if required {
					*errs = append(*errs, NewPathErrors(PathMember(name), fmt.Errorf("property requires a value"))...)
				}
			} else {
				if err := json.Unmarshal(raw, dest.FieldByIndex(fld.Index).Addr().Interface()); err != nil {
					*errs = append(*errs, NewPathErrors(PathMember(name), err)...)
				}

				delete(props, name)
			}
		},
	)
}

func reflectStruct(t reflect.Type, base func(fld reflect.StructField), property func(fld reflect.StructField, name string, required bool)) {
	if t.Kind() != reflect.Struct {
		panic(fmt.Errorf("type '%s' is not a structure", t.String()))
	}

	for i := 0; i < t.NumField(); i++ {
		fld := t.Field(i)
		tag := fld.Tag.Get("json")

		if tag != "-" {
			if fld.Anonymous && tag == "" {
				// Is a base structure
				base(fld)
			} else {
				tflds := strings.Split(tag, ",")
				name := tflds[0]
				opts := tflds[1:]

				if name == "" {
					name = fld.Name
				}

				if slices.Contains(opts, "string") {
					panic(fmt.Errorf("type '%s', field '%s': 'string' json marshaling is not supported", t.String(), fld.Name))
				}

				property(fld, name, !slices.Contains(opts, "omitempty"))
			}
		}
	}
}
