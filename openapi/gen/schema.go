package gen

import (
	"encoding"
	"encoding/json"
	"fmt"
	"go/token"
	"math"
	"reflect"

	"github.com/iancoleman/strcase"
	"github.com/trwk76/goweb/openapi/marshal"
	"github.com/trwk76/goweb/openapi/spec"
)

func SchemaTypePascalCase(t reflect.Type) string {
	if t.PkgPath() == "" || !token.IsIdentifier(t.Name()) {
		return ""
	}

	return strcase.ToCamel(t.Name())
}

func SchemaTypeCamelCase(t reflect.Type) string {
	if t.PkgPath() == "" || !token.IsIdentifier(t.Name()) {
		return ""
	}

	return strcase.ToLowerCamel(t.Name())
}

func SchemaTypeSnakeCase(t reflect.Type) string {
	if t.PkgPath() == "" || !token.IsIdentifier(t.Name()) {
		return ""
	}

	return strcase.ToSnake(t.Name())
}

func SchemaOrRefFor[T any](a *API) spec.SchemaOrRef {
	return a.sch.SchemaOrRefOf(reflect.TypeFor[T]())
}

func (m *Schemas) SchemaOrRefOf(t reflect.Type) spec.SchemaOrRef {
	var res spec.SchemaOrRef
	var ptr bool = false

	if t.Kind() == reflect.Pointer {
		ptr = true
		t = t.Elem()
	}

	if key, fnd := m.types[t]; fnd {
		res = spec.SchemaOrRef{Ref: spec.ComponentsRef("schemas", key)}
	} else {
		key := ""

		if m.keyf != nil {
			if key = m.keyf(t); key != "" {
				key = uniqueName(m.keys, key)
				m.types[t] = key
			}
		}

		sch := SchemaOf(t, m)

		if key != "" {
			m.keys[key] = sch
			res = spec.SchemaOrRef{Ref: spec.ComponentsRef("schemas", key)}
		} else {
			res = spec.SchemaOrRef{Item: sch}
		}
	}

	if ptr {
		// type is nullable => <res> = {oneOf: {<res>, {type: "null"}}}
		return spec.SchemaOrRef{Item: spec.Schema{OneOf: []spec.SchemaOrRef{res, {Item: spec.Schema{Type: spec.TypeNull}}}}}
	}

	return res
}

func (m *Schemas) Schema(item spec.SchemaOrRef) (spec.Schema, bool) {
	if item.Ref.Ref != "" {
		key, ok := spec.ComponentsKey(item.Ref.Ref, "schemas")
		if !ok {
			return spec.Schema{}, false
		}

		res, ok := m.keys[key]
		return res, ok
	}

	return item.Item, true
}

func SchemaOf(t reflect.Type, m marshal.SchemaTypeMapping) spec.Schema {
	pval := reflect.New(t)
	ival := reflect.Indirect(pval)

	if isch, ok := ival.Interface().(Schemater); ok {
		res := isch.Schema()
		reqMarsh := false

		switch res.Type {
		case spec.TypeNull:
		case spec.TypeBoolean:
			if t.Kind() != reflect.Bool {
				reqMarsh = true
			}
		case spec.TypeInteger:
			switch t.Kind() {
			case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8, reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
			default:
				reqMarsh = true
			}
		case spec.TypeNumber:
			switch t.Kind() {
			case reflect.Float32, reflect.Float64, reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8, reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
			default:
				reqMarsh = true
			}
		case spec.TypeString:
			if t.Kind() != reflect.String {
				_, tmar := ival.Interface().(encoding.TextMarshaler)
				_, tunm := pval.Interface().(encoding.TextUnmarshaler)

				if tmar || tunm {
					if !tunm {
						panic(fmt.Errorf("type '%s' implements encoding.TextMarshaler but (%s).UnmarshalText() is missing", ival.Type().String(), pval.Type().String()))
					}

					if !tmar {
						panic(fmt.Errorf("type '%s' implements encoding.TextUnmarshaler but (%s).MarshalText() is missing", pval.Type().String(), ival.Type().String()))
					}
				} else {
					reqMarsh = true
				}
			}
		case spec.TypeArray:
			switch t.Kind() {
			case reflect.Array, reflect.Slice:
			default:
				reqMarsh = true
			}
		case spec.TypeObject:
			switch t.Kind() {
			case reflect.Map, reflect.Struct:
			default:
				reqMarsh = true
			}
		case spec.TypeNone:
			if len(res.AllOf) < 1 || t.Kind() != reflect.Struct {
				reqMarsh = true
			}
		}

		if reqMarsh {
			_, tmar := ival.Interface().(json.Marshaler)
			_, tunm := pval.Interface().(json.Unmarshaler)

			if tmar || tunm {
				if !tunm {
					panic(fmt.Errorf("type '%s' implements json.Marshaler but (%s).UnmarshalJSON() is missing", ival.Type().String(), pval.Type().String()))
				}

				if !tmar {
					panic(fmt.Errorf("type '%s' implements json.Unmarshaler but (%s).MarshalJSON() is missing", pval.Type().String(), ival.Type().String()))
				}
			} else {
				panic(fmt.Errorf("type '%s' does not implement json.Marshaler / json.Unmarshaler", t.String()))
			}
		}

		return res
	}

	switch t.Kind() {
	case reflect.Array:
		itms := m.SchemaOrRefOf(t.Elem())
		l := uint32(t.Len())

		return spec.Schema{
			Type:     spec.TypeArray,
			Items:    &itms,
			MinItems: l,
			MaxItems: l,
		}
	case reflect.Bool:
		return spec.Schema{Type: spec.TypeBoolean}
	case reflect.Float32:
		return marshal.FloatSchema(spec.FormatFloat, -math.MaxFloat32, math.MaxFloat32, 0)
	case reflect.Float64:
		return marshal.FloatSchema(spec.FormatDouble, -math.MaxFloat64, math.MaxFloat64, 0)
	case reflect.Int, reflect.Int64:
		// assumption: we are no longer running on a 32bit processor (int = int64)
		return marshal.IntSchema(spec.FormatInt64, math.MinInt64, math.MaxInt64, 0)
	case reflect.Int16:
		return marshal.IntSchema(spec.FormatNone, math.MinInt16, math.MaxInt16, 0)
	case reflect.Int32:
		return marshal.IntSchema(spec.FormatInt32, math.MinInt32, math.MaxInt32, 0)
	case reflect.Int8:
		return marshal.IntSchema(spec.FormatNone, math.MinInt8, math.MaxInt8, 0)
	case reflect.Map:
		key := SchemaOf(t.Key(), m)

		if key.Type != spec.TypeString {
			panic(fmt.Errorf("type '%s': map's key type '%s' is not a string", t.String(), t.Key().String()))
		}

		items := m.SchemaOrRefOf(t.Elem())

		if key.Pattern != "" {
			return spec.Schema{
				Type:              spec.TypeObject,
				PatternProperties: spec.NamedSchemaOrRefs{key.Pattern: items},
			}
		}

		return spec.Schema{
			Type:                 spec.TypeObject,
			AdditionalProperties: &items,
		}
	case reflect.Slice:
		itms := m.SchemaOrRefOf(t.Elem())

		return spec.Schema{
			Type:  spec.TypeArray,
			Items: &itms,
		}
	case reflect.String:
		return spec.Schema{Type: spec.TypeString}
	case reflect.Struct:
		return marshal.StructSchemaOf(t, m)
	case reflect.Uint, reflect.Uint64:
		// assumption: we are no longer running on a 32bit processor (int = int64)
		return marshal.UintSchema(spec.FormatNone, 0, uint64(math.MaxUint64), 0)
	case reflect.Uint16:
		return marshal.UintSchema(spec.FormatNone, 0, uint16(math.MaxUint16), 0)
	case reflect.Uint32:
		return marshal.UintSchema(spec.FormatInt32, 0, uint32(math.MaxUint32), 0)
	case reflect.Uint8:
		return marshal.UintSchema(spec.FormatNone, 0, uint8(math.MaxUint8), 0)
	}

	panic(fmt.Errorf("type '%s': schema extraction by reflection is not supported", t.String()))
}

type (
	Schemater interface {
		Schema() spec.Schema
	}

	Schemas struct {
		keyf  SchemaTypeKeyFunc
		keys  map[string]spec.Schema
		types map[reflect.Type]string
	}

	SchemaTypeKeyFunc func(t reflect.Type) string
)

func newSchemas(keyFunc SchemaTypeKeyFunc) Schemas {
	return Schemas{
		keyf:  keyFunc,
		keys:  make(map[string]spec.Schema),
		types: make(map[reflect.Type]string),
	}
}

func (m *Schemas) spec() spec.NamedSchemas {
	return m.keys
}

func checkSimpleSchema(sch spec.Schema) {
	switch sch.Type {
	case spec.TypeArray, spec.TypeNone, spec.TypeNull, spec.TypeObject:
		panic(fmt.Errorf("only simple schemas can be used with parameters and headers"))
	}
}

var (
	_ marshal.SchemaTypeMapping = (*Schemas)(nil)
	_ SchemaTypeKeyFunc         = SchemaTypeCamelCase
)
