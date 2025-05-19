package form

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

	name, opts, ok := strings.Cut(s, ",")
	res.Name = name

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
	// Tag represents a form tag for a field in a struct.
	// It contains the name of the field (Name),
	// whether the field is optional (Opt), and the delimiter for splitting values (Del).
	// The Name field specifies the name of the field in the source.
	// The Opt field indicates if the field is optional (true) or required (false).
	// The Del field specifies the delimiter used to split multiple values in the source.
	Tag struct {
		Name string
		Opt  bool
		Del  TagDel
	}

	// TagDel represents the delimiter used to split multiple values in the source (e.g., comma, space, pipe).
	TagDel string
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
