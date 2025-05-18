package bind

import (
	"fmt"
	"reflect"
	"strings"
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

	loc, opts, ok := strings.Cut(s, ",")
	if loc == "body" {
		res.Src = TagSrcBody
		res.Name = loc
	} else {
		src, name, ok := strings.Cut(loc, "/")
		if !ok {
			return res, fmt.Errorf("invalid bind location %s", loc)
		}

		switch src {
		case "cookie":
			res.Src = TagSrcCookie
			res.Name = name
		case "header":
			res.Src = TagSrcHeader
			res.Name = name
		case "path":
			res.Src = TagSrcPath
			res.Name = name
		case "query":
			res.Src = TagSrcQuery
			res.Name = name
		default:
			return res, fmt.Errorf("unsupported bind source %s", src)
		}
	}

	if ok {
		for _, opt := range strings.Split(opts, ",") {
			switch opt {
			case "opt":
				res.Opt = true
			case "comma":
				res.Del = TagDelComma
			case "space":
				res.Del = TagDelSpace
			case "pipe":
				res.Del = TagDelPipe
			default:
				return res, fmt.Errorf("unsupported option %s", opt)
			}
		}
	}

	return res, nil
}

type (
	// Tag represents a binding tag for a field in a struct.
	// It contains the source of the binding (Src), the name of the field (Name),
	// whether the field is optional (Opt), and the delimiter for splitting values (Del).
	// The Src field indicates where the value should be bound from (e.g., body, cookie, header, path, query).
	// The Name field specifies the name of the field in the source.
	// The Opt field indicates if the field is optional (true) or required (false).
	// The Del field specifies the delimiter used to split multiple values in the source.
	Tag struct {
		Src  TagSrc
		Name string
		Opt  bool
		Del  TagDel
	}

	// TagSrc represents the source of the binding (e.g., body, cookie, header, path, query).
	TagSrc string

	// TagDel represents the delimiter used to split multiple values in the source (e.g., comma, space, pipe).
	TagDel string
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

const (
	// TagDelNone represents no delimiter for splitting values.
	TagDelNone TagDel = ""
	// TagDelComma represents a comma (,) delimiter to be used for splitting multiple values.
	TagDelComma TagDel = ","
	// TagDelSpace represents a space ( ) delimiter to be used for splitting multiple values.
	TagDelSpace TagDel = " "
	// TagDelPipe represents a comma (|) delimiter to be used for splitting multiple values.
	TagDelPipe TagDel = "|"
)
