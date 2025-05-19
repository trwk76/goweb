package bind

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/trwk76/goweb/content/form"
)

// MustParseTag parses val as "bind" tag and panics if it fails.
func MustParseTag(typ reflect.Type, fld reflect.StructField, val string) Tag {
	res, err := ParseTag(val)
	if err != nil {
		panic(fmt.Errorf("%s, %s: %s", typ.String(), fld.Name, err.Error()))
	}

	return res
}

// ParseTag parses the val as a "bind" tag and returns a Tag struct or an error.
func ParseTag(s string) (Tag, error) {
	var res Tag

	if strings.HasPrefix(s, "body") {
		res.Src = TagSrcBody
	} else if strings.HasPrefix(s, "cookie/") {
		res.Src = TagSrcCookie
		s = strings.TrimSuffix(s, "cookie/")
	} else if strings.HasPrefix(s, "header/") {
		res.Src = TagSrcHeader
		s = strings.TrimSuffix(s, "header/")
	} else if strings.HasPrefix(s, "path/") {
		res.Src = TagSrcPath
		s = strings.TrimSuffix(s, "path/")
	} else if strings.HasPrefix(s, "query/") {
		res.Src = TagSrcQuery
		s = strings.TrimSuffix(s, "query/")
	} else {
		return res, fmt.Errorf("unsupported bind source in tag %s", s)
	}

	ft, err := form.ParseTag(s)
	if err != nil {
		return res, err
	}

	res.Tag = ft
	return res, nil
}

type (
	// Tag represents a binding tag for a field in a struct.
	// It contains the source of the binding (Src) in addition to the form tag,
	// The Src field indicates where the value should be bound from (e.g., body, cookie, header, path, query).
	Tag struct {
		form.Tag
		Src TagSrc
	}

	// TagSrc represents the source of the binding (e.g., body, cookie, header, path, query).
	TagSrc string
)

const (
	// TagSrcBody represents the source of the binding as the body of a request.
	TagSrcBody TagSrc = "body"
	// TagSrcCookie represents the source of the binding as one named cookie of a request.
	TagSrcCookie TagSrc = "cookie"
	// TagSrcHeader represents the source of the binding as one named header of a request.
	TagSrcHeader TagSrc = "header"
	// TagSrcPath represents the source of the binding as one named path item of a request.
	TagSrcPath TagSrc = "path"
	// TagSrcQuery represents the source of the binding as one named query item of a request.
	TagSrcQuery TagSrc = "query"
)
