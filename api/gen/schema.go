package gen

import (
	"fmt"
	"math"
	"reflect"
	"strings"

	"github.com/trwk76/goweb/api/spec"
)

func Schema[T any](s *APISpec) SchemaSpec {
	var d T

	return s.Schemas.typ(reflect.TypeOf(d))
}

type (
	Schemater interface {
		Schema() spec.Schema
	}

	SchemaSpecs []*namedSchema

	namedSchema struct {
		name string
		impl schemaImpl
	}

	SchemaSpec interface {
		spec() spec.Schema
		impl() schemaImpl
	}

	schemaImpl struct {
		typ reflect.Type
		spc spec.Schema
	}

	schemaRef struct {
		ptr *namedSchema
	}
)

func (s *SchemaSpecs) typ(t reflect.Type) SchemaSpec {
	var res SchemaSpec
	var dst *schemaImpl

	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	for _, itm := range *s {
		if itm.impl.typ == t {
			return schemaRef{ptr: itm}
		}
	}

	if t.Name() != "" && t.PkgPath() != "" {
		named := &namedSchema{
			name: s.newName(t.Name()),
		}

		dst = &named.impl
		*s = append(*s, named)
		res = schemaRef{ptr: named}
	} else {
		dst = &schemaImpl{}
	}

	dst.typ = t

	val := reflect.Indirect(reflect.New(t))
	if sch, ok := val.Interface().(Schemater); ok {
		dst.spc = sch.Schema()
	} else {
		switch t.Kind() {
		case reflect.Array:
			itm := s.typ(t.Elem()).spec()
			dst.spc.Type = spec.TypeArray
			dst.spc.Items = &itm
			dst.spc.MinItems = uint32(t.Len())
			dst.spc.MaxItems = uint32(t.Len())
		case reflect.Bool:
			dst.spc.Type = spec.TypeBoolean
		case reflect.Float32:
			dst.spc = numSchema(spec.TypeNumber, spec.FormatFloat32, -math.MaxFloat32, math.MaxFloat32)
		case reflect.Float64:
			dst.spc = numSchema(spec.TypeNumber, spec.FormatFloat64, -math.MaxFloat64, math.MaxFloat64)
		case reflect.Int16:
			dst.spc = numSchema(spec.TypeInteger, spec.FormatInt16, math.MinInt16, math.MaxInt16)
		case reflect.Int32:
			dst.spc = numSchema(spec.TypeInteger, spec.FormatInt32, math.MinInt32, math.MaxInt32)
		case reflect.Int64:
			dst.spc = numSchema(spec.TypeInteger, spec.FormatInt64, math.MinInt64, math.MaxInt64)
		case reflect.Int8:
			dst.spc = numSchema(spec.TypeInteger, spec.FormatInt8, math.MinInt8, math.MaxInt8)
		case reflect.Map:
			key := s.typ(t.Key()).impl()
			val := s.typ(t.Elem())

			if key.spc.Type != spec.TypeString {
				panic(fmt.Errorf("map type '%s' must have a string key type", t.String()))
			}

			dst.spc.Type = spec.TypeObject
			if key.spc.Pattern != "" {
				dst.spc.AdditionalProperties = &spec.AdditionalProperties{Bool: false}
				dst.spc.PatternProperties = spec.NamedSchemas{key.spc.Pattern: val.spec()}
			} else {
				itm := val.spec()
				dst.spc.AdditionalProperties = &spec.AdditionalProperties{Schema: &itm}
			}
		case reflect.Slice:
			itm := s.typ(t.Elem()).spec()
			dst.spc.Type = spec.TypeArray
			dst.spc.Items = &itm
		case reflect.String:
			dst.spc.Type = spec.TypeString
		case reflect.Struct:
			dst.spc.Type = spec.TypeObject
			dst.spc.AdditionalProperties = &spec.AdditionalProperties{Bool: false}
			dst.spc.Properties = make(spec.NamedSchemas)

			for _, fld := range reflect.VisibleFields(t) {
				tag := fld.Tag.Get("json")

				if tag != "-" {
					opt := strings.HasSuffix(tag, ",omitempty")
					tag = strings.TrimSuffix(tag, ",omitempty")

					if tag == "" {
						tag = fld.Name
					}

					if !opt {
						dst.spc.Required = append(dst.spc.Required, tag)
					}

					dst.spc.Properties[tag] = s.typ(fld.Type).spec()
				}
			}
		case reflect.Uint16:
			dst.spc = numSchema(spec.TypeInteger, spec.FormatUint16, 0, math.MaxUint16)
		case reflect.Uint32:
			dst.spc = numSchema(spec.TypeInteger, spec.FormatUint32, 0, math.MaxUint32)
		case reflect.Uint64:
			dst.spc = numSchema(spec.TypeInteger, spec.FormatUint64, 0, math.MaxUint64)
		case reflect.Uint8:
			dst.spc = numSchema(spec.TypeInteger, spec.FormatUint8, 0, math.MaxUint8)
		default:
			panic(fmt.Errorf("type '%s' cannot be used as a JSON schema", t.String()))
		}
	}

	if res == nil {
		res = *dst
	}

	return res
}

func (s *SchemaSpecs) newName(name string) string {
	if !s.hasName(name) {
		return name
	}

	base := name
	i := 1
	name = fmt.Sprintf("%s%d", base, i)

	for s.hasName(name) {
		i++
		name = fmt.Sprintf("%s%d", base, i)
	}

	return name
}

func (s *SchemaSpecs) hasName(name string) bool {
	for _, itm := range *s {
		if itm.name == name {
			return true
		}
	}

	return false
}

func (s SchemaSpecs) spec() spec.NamedSchemas {
	res := make(spec.NamedSchemas)

	for _, itm := range s {
		res[itm.name] = itm.impl.spc
	}

	return res
}

func (s schemaImpl) spec() spec.Schema {
	return s.spc
}

func (s schemaImpl) impl() schemaImpl {
	return s
}

func (s schemaImpl) isSimple() bool {
	switch s.spc.Type {
	case spec.TypeBoolean, spec.TypeInteger, spec.TypeNumber, spec.TypeString:
		return true
	}

	return false
}

func (r schemaRef) spec() spec.Schema {
	return spec.Schema{Ref: "#/components/schemas/" + r.ptr.name}
}

func (r schemaRef) impl() schemaImpl {
	return r.ptr.impl
}

func numSchema(t spec.Type, f spec.Format, min float64, max float64) spec.Schema {
	return spec.Schema{
		Type:    t,
		Format:  f,
		Minimum: &min,
		Maximum: &max,
	}
}

var (
	_ SchemaSpec = schemaImpl{}
	_ SchemaSpec = schemaRef{}
)
