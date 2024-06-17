package gen

import (
	"fmt"
	"go/token"
	"math"
	"reflect"

	"github.com/iancoleman/strcase"
	"github.com/trwk76/goweb/openapi/spec"
)

func SchemaFor[T any](a *API, mediaType MediaType) spec.Schema {
	return a.SchemaOf(reflect.TypeFor[T](), mediaType)
}

func SchemaOrRefFor[T any](a *API, mediaType MediaType) spec.SchemaOrRef {
	return a.SchemaOrRefOf(reflect.TypeFor[T](), mediaType)
}

func (a *API) SchemaOf(t reflect.Type, mediaType MediaType) spec.Schema {
	return a.Schema(a.SchemaOrRefOf(t, mediaType))
}

func (a *API) SchemaOrRefOf(t reflect.Type, mediaType MediaType) spec.SchemaOrRef {
	var res spec.Schema

	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	if se, ok := a.sch[t]; ok {
		if se.simple != "" {
			return spec.SchemaOrRef{Ref: spec.Ref("schemas", se.simple)}
		}

		if mediaType == nil {
			panic(fmt.Errorf("complex type requires a mediaType"))
		}

		if key, ok := se.mtypes[mediaType.ContentType()]; ok {
			return spec.SchemaOrRef{Ref: spec.Ref("schemas", key)}
		}
	}

	if sch, ok := reflect.Indirect(reflect.New(t)).Interface().(spec.Schemater); ok {
		res = sch.Schema()
	} else {
		switch t.Kind() {
		case reflect.Array:
			itm := a.SchemaOrRefOf(t.Elem(), mediaType)

			res.Type = spec.TypeArray
			res.Items = &itm
			res.MinItems = uint32(t.Len())
			res.MaxItems = uint32(t.Len())
		case reflect.Bool:
			res.Type = spec.TypeBoolean
		case reflect.Float32:
			res = numSchema(spec.TypeNumber, spec.FormatFloat, float32(-math.MaxFloat32), math.MaxFloat32)
		case reflect.Float64:
			res = numSchema(spec.TypeNumber, spec.FormatDouble, float64(-math.MaxFloat64), math.MaxFloat64)
		case reflect.Int16:
			res = numSchema(spec.TypeInteger, spec.FormatNone, int16(math.MinInt16), math.MaxInt16)
		case reflect.Int32:
			res = numSchema(spec.TypeInteger, spec.FormatInt32, int32(math.MinInt32), math.MaxInt32)
		case reflect.Int64:
			res = numSchema(spec.TypeInteger, spec.FormatInt64, int64(math.MinInt64), math.MaxInt64)
		case reflect.Int8:
			res = numSchema(spec.TypeInteger, spec.FormatNone, int8(math.MinInt8), math.MaxInt8)
		case reflect.Map:
			key := a.SchemaOf(t.Key(), mediaType)
			if key.Type != spec.TypeString {
				panic(fmt.Errorf("type '%s' is not a string though used as a map key", t.Key().String()))
			}

			itm := a.SchemaOrRefOf(t.Elem(), mediaType)

			res.Type = spec.TypeObject
			res.AdditionalProperties = &itm
		case reflect.Slice:
			itm := a.SchemaOrRefOf(t.Elem(), mediaType)

			res.Type = spec.TypeArray
			res.Items = &itm
		case reflect.String:
			res.Type = spec.TypeString
		case reflect.Struct:
			bases := make([]spec.SchemaOrRef, 0)
			props := make(spec.NamedSchemaOrRefs)
			req := make([]string, 0)

			res.Type = spec.TypeObject
			res.Required = req
			res.Properties = props

			for i := 0; i < t.NumField(); i++ {
				mediaType.ReflectField(a, t.Field(i), &bases, props, &req)
			}

			if len(bases) > 0 {
				res = spec.Schema{AllOf: append(bases, spec.SchemaOrRef{Item: res})}
			}
		case reflect.Uint16:
			res = numSchema(spec.TypeInteger, spec.FormatNone, uint16(0), math.MaxUint16)
		case reflect.Uint32:
			res = numSchema(spec.TypeInteger, spec.FormatNone, uint32(0), math.MaxUint32)
		case reflect.Uint64:
			res = numSchema(spec.TypeInteger, spec.FormatNone, uint64(0), math.MaxUint64)
		case reflect.Uint8:
			res = numSchema(spec.TypeInteger, spec.FormatNone, uint8(0), math.MaxUint8)
		}
	}

	if t.PkgPath() != "" && token.IsIdentifier(t.Name()) {
		key := uniqueName(a.spc.Components.Schemas, strcase.ToLowerCamel(t.Name()))

		se, ok := a.sch[t]

		switch res.Type {
		case spec.TypeNull, spec.TypeBoolean, spec.TypeInteger, spec.TypeNumber, spec.TypeString:
			se.simple = key
		default:
			if !ok {
				se.mtypes = make(map[string]string)
			}

			se.mtypes[mediaType.ContentType()] = key
		}

		a.sch[t] = se
		a.spc.Components.Schemas[key] = res

		return spec.SchemaOrRef{Ref: spec.Ref("schemas", key)}
	}

	return spec.SchemaOrRef{Item: res}
}

func (a *API) Schema(s spec.SchemaOrRef) spec.Schema {
	if s.Ref.Ref != "" {
		return a.spc.Components.Schemas[refKey(s.Ref.Ref)]
	}

	return s.Item
}

type (
	schemas map[reflect.Type]schemaEntry

	schemaEntry struct {
		simple string
		mtypes map[string]string
	}
)

func numSchema[T comparable](t spec.Type, f spec.Format, min T, max T) spec.Schema {
	return spec.Schema{
		Type:    t,
		Format:  f,
		Minimum: min,
		Maximum: max,
	}
}
