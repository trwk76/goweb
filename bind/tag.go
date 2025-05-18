package bind

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
	Tag struct {
		Src  TagSrc
		Name string
		Opt  bool
		Del  TagDel
	}

	TagSrc string
	TagDel string
)

const (
	TagSrcBody   TagSrc = "body"
	TagSrcCookie TagSrc = "cookie"
	TagSrcHeader TagSrc = "header"
	TagSrcPath   TagSrc = "path"
	TagSrcQuery  TagSrc = "query"
)

const (
	TagDelNone  TagDel = ""
	TagDelComma TagDel = ","
	TagDelSpace TagDel = " "
	TagDelPipe  TagDel = "|"
)
