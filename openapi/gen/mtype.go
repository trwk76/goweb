package gen

import (
	"fmt"
	"reflect"
	"slices"

	web "github.com/trwk76/goweb"
	"github.com/trwk76/goweb/openapi/spec"
)

type (
	MediaTypes struct {
		items []MediaType
	}

	MediaType interface {
		ContentType() string
		ReflectField(a *API, fld reflect.StructField, bases *[]spec.SchemaOrRef, props spec.NamedSchemaOrRefs, req *[]string)
	}

	JSONMediaType struct {
		ctype string
	}

	YAMLMediaType struct {
		ctype string
	}
)

var (
	JSON JSONMediaType = JSONMediaType{ctype: web.ContentTypeJSON}
	YAML YAMLMediaType = YAMLMediaType{ctype: web.ContentTypeYAML}
)

func NewMediaTypes(items ...MediaType) MediaTypes {
	for idx, itm := range items {
		if slices.IndexFunc(items, func(m MediaType) bool { return m.ContentType() == itm.ContentType() }) != idx {
			panic(fmt.Errorf("content type '%s' specified more than once", itm.ContentType()))
		}
	}

	return MediaTypes{items: items}
}

func NewJSONMediaType(ctype string) JSONMediaType {
	return JSONMediaType{ctype: ctype}
}

func NewYAMLMediaType(ctype string) YAMLMediaType {
	return YAMLMediaType{ctype: ctype}
}
