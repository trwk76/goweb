package json

import (
	"fmt"
	"reflect"
	"strings"
)

func MustParseTag(typ reflect.Type, fld reflect.StructField, val string) Tag {
	res, err := ParseTag(val)
	if err != nil {
		panic(fmt.Errorf("%s, %s: %s", typ.String(), fld.Name, err.Error()))
	}

	return res
}

func ParseTag(s string) (Tag, error) {
	var res Tag

	name, opts, ok := strings.Cut(s, ",")
	res.Name = name

	if ok {
		for _, opt := range strings.Split(opts, ",") {
			switch opt {
			case "omitempty":
				res.OmitEmpty = true
			case "omitzero":
				res.OmitZero = true
			case "string":
				res.String = true
			default:
				return res, fmt.Errorf("unsupported json tag option '%s'", opt)
			}
		}
	}

	return res, nil
}

type (
	Tag struct {
		Name      string
		OmitEmpty bool
		OmitZero  bool
		String    bool
	}
)
