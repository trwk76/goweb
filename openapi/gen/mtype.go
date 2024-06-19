package gen

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/trwk76/goweb/openapi/spec"
	"gopkg.in/yaml.v3"
)

func ContentFor[T any](a *API, mtypes MediaTypes, examples Examples[T]) spec.MediaTypes {
	t := reflect.TypeFor[T]()
	res := make(spec.MediaTypes)

	for key, mtyp := range mtypes.Items {
		sch := a.SchemaOrRefOf(t, key, mtyp)

		res[key] = spec.MediaType{
			Schema:   &sch,
			Examples: examples.spec(mtyp),
		}
	}

	return res
}

type (
	MediaTypes struct {
		Default string
		Items   map[string]MediaType
	}

	MediaType interface {
		ReflectField(a *API, key string, fld reflect.StructField, bases *[]spec.SchemaOrRef, props spec.NamedSchemaOrRefs, req *[]string)
		Example(value any) any
	}

	JSONMediaType struct{}
	YAMLMediaType struct{}
)

var (
	JSON JSONMediaType = JSONMediaType{}
	YAML YAMLMediaType = YAMLMediaType{}
)

func (t JSONMediaType) ReflectField(a *API, key string, fld reflect.StructField, bases *[]spec.SchemaOrRef, props spec.NamedSchemaOrRefs, req *[]string) {
	tag := fld.Tag.Get("json")

	if tag == "-" {
		return
	}

	flds := strings.Split(tag, ",")
	name := flds[0]
	flds = flds[1:]

	if name == "" && fld.Anonymous {
		// base
		*bases = append(*bases, a.SchemaOrRefOf(fld.Type, key, t))
		return
	}

	if name == "" {
		name = fld.Name
	}

	if slices.Contains(flds, "string") {
		props[name] = spec.SchemaOrRef{Item: spec.Schema{Type: spec.TypeString}}
	} else {
		props[name] = a.SchemaOrRefOf(fld.Type, key, t)
	}

	if !slices.Contains(flds[1:], "omitempty") {
		*req = append(*req, name)
	}
}

func (t JSONMediaType) Example(value any) any {
	return value
}

func (t YAMLMediaType) ReflectField(a *API, key string, fld reflect.StructField, bases *[]spec.SchemaOrRef, props spec.NamedSchemaOrRefs, req *[]string) {
	tag := fld.Tag.Get("yaml")

	if tag == "-" {
		return
	}

	flds := strings.Split(tag, ",")
	name := flds[0]
	flds = flds[1:]

	if slices.Contains(flds, "inline") {
		// base
		*bases = append(*bases, a.SchemaOrRefOf(fld.Type, key, t))
		return
	}

	if name == "" {
		name = fld.Name
	}

	props[name] = a.SchemaOrRefOf(fld.Type, key, t)

	if !slices.Contains(flds, "omitempty") {
		*req = append(*req, name)
	}
}

func (t YAMLMediaType) Example(value any) any {
	raw, err := yaml.Marshal(value)
	if err != nil {
		panic(fmt.Errorf("error marshaing example %v: %s", value, err.Error()))
	}

	return string(raw)
}

var (
	_ MediaType = JSON
	_ MediaType = YAML
)
