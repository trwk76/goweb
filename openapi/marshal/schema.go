package marshal

import (
	"reflect"
	"regexp"

	"github.com/trwk76/goweb/openapi/spec"
)

func ArraySchemaFor[T ~[]I, I any](m SchemaTypeMapping, minLength uint32, maxLength uint32, unique bool) spec.Schema {
	return ArraySchema(m.SchemaOrRefOf(reflect.TypeFor[I]()), minLength, maxLength, unique)
}

func ArraySchema(items spec.SchemaOrRef, minLength uint32, maxLength uint32, unique bool) spec.Schema {
	return spec.Schema{
		Type:        spec.TypeArray,
		Items:       &items,
		MinItems:    minLength,
		MaxLength:   maxLength,
		UniqueItems: unique,
	}
}

func EnumSchema[T comparable](items ...T) spec.Schema {
	vals := make([]any, len(items))

	for idx, itm := range items {
		vals[idx] = itm
	}

	return spec.Schema{Enum: vals}
}

func IntSchema[T integer](format spec.Format, min T, max T, multipleOf T) spec.Schema {
	var mult any

	if multipleOf != T(0) {
		mult = multipleOf
	}

	return spec.Schema{Type: spec.TypeInteger, Format: format, Minimum: min, Maximum: max, MultipleOf: mult}
}

func UintSchema[T uinteger](format spec.Format, min T, max T, multipleOf T) spec.Schema {
	var mult any

	if multipleOf != T(0) {
		mult = multipleOf
	}

	return spec.Schema{Type: spec.TypeInteger, Format: format, Minimum: min, Maximum: max, MultipleOf: mult}
}

func FloatSchema[T float](format spec.Format, min T, max T, multipleOf T) spec.Schema {
	var mult any

	if multipleOf != T(0) {
		mult = multipleOf
	}

	return spec.Schema{Type: spec.TypeNumber, Format: format, Minimum: min, Maximum: max, MultipleOf: mult}
}

func StringSchema(format spec.Format, minLength uint32, maxLength uint32, regex *regexp.Regexp) spec.Schema {
	res := spec.Schema{
		Type:      spec.TypeString,
		Format:    format,
		MinLength: minLength,
		MaxLength: maxLength,
	}

	if regex != nil {
		res.Pattern = regex.String()
	}

	return res
}

func StructSchemaFor[T any](m SchemaTypeMapping) spec.Schema {
	return StructSchemaOf(reflect.TypeFor[T](), m)
}

func StructSchemaOf(t reflect.Type, m SchemaTypeMapping) spec.Schema {
	bases := make([]spec.SchemaOrRef, 0)
	res := spec.Schema{Type: spec.TypeObject, Properties: make(spec.NamedSchemaOrRefs)}

	reflectStruct(
		t,
		func(fld reflect.StructField) {
			bases = append(bases, m.SchemaOrRefOf(fld.Type))
		},
		func(fld reflect.StructField, name string, required bool) {
			if required {
				res.Required = append(res.Required, name)
			}

			res.Properties[name] = m.SchemaOrRefOf(fld.Type)
		},
	)

	if len(bases) > 0 {
		res = spec.Schema{AllOf: append(bases, spec.SchemaOrRef{Item: res})}
	}

	return res
}

type (
	SchemaTypeMapping interface {
		SchemaOrRefOf(t reflect.Type) spec.SchemaOrRef
	}
)
