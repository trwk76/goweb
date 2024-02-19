package api

import (
	"encoding/json"
	"fmt"

	web "github.com/trwk76/goweb"
)

func Parse[T any](dst *T, ctx *web.Context, name string, in In, req bool, err *ParseError) bool {
	switch in {
	case InBody:
		return ParseBody(dst, ctx, req, err)
	case InCookie:
		return ParseCookie(dst, ctx, name, req, err)
	case InHeader:
		return ParseHeader(dst, ctx, name, req, err)
	case InPath:
		return ParsePath(dst, ctx, name, err)
	case InQuery:
		return ParseQuery(dst, ctx, name, req, err)
	}

	panic(fmt.Errorf("parse location '%s' is not supported", in))
}

func ParsePtr[T any](dst **T, ctx *web.Context, name string, in In, req bool, err *ParseError) bool {
	var val T

	if !Parse(&val, ctx, name, in, req, err) {
		return false
	}

	*dst = &val
	return true
}

func ParseBody[T any](dst *T, ctx *web.Context, req bool, err *ParseError) bool {
	if perr := json.NewDecoder(ctx.Request().Body).Decode(dst); perr != nil {
		err.Add(string(InBody), perr.Error())
		return false
	}

	return true
}

func ParseCookie[T any](dst *T, ctx *web.Context, name string, req bool, err *ParseError) bool {
	cook, cerr := ctx.Request().Cookie(name)
	if cerr != nil {
		if req {
			err.AddRequired(toPath(InCookie, name))
		}

		return false
	}

	if perr := unmarshalText(dst, []byte(cook.Value)); perr != nil {
		err.Add(toPath(InCookie, name), perr.Error())
		return false
	}

	return true
}

func ParseHeader[T any](dst *T, ctx *web.Context, name string, req bool, err *ParseError) bool {
	text := ctx.Request().Header.Get(name)
	if text == "" {
		if req {
			err.AddRequired(toPath(InHeader, name))
		}

		return false
	}

	if perr := unmarshalText(dst, []byte(text)); perr != nil {
		err.Add(toPath(InHeader, name), perr.Error())
		return false
	}

	return true
}

func ParsePath[T any](dst *T, ctx *web.Context, name string, err *ParseError) bool {
	text := ctx.PathParam(name)

	if perr := unmarshalText(dst, []byte(text)); perr != nil {
		err.Add(toPath(InPath, name), perr.Error())
		return false
	}

	return true
}

func ParseQuery[T any](dst *T, ctx *web.Context, name string, req bool, err *ParseError) bool {
	text := ctx.Request().URL.Query().Get(name)
	if text == "" {
		if req {
			err.AddRequired(toPath(InQuery, name))
		}

		return false
	}

	if perr := unmarshalText(dst, []byte(text)); perr != nil {
		err.Add(toPath(InQuery, name), perr.Error())
		return false
	}

	return true
}

type (
	In string
)

const (
	InBody   In = "body"
	InCookie In = "cookie"
	InHeader In = "header"
	InPath   In = "path"
	InQuery  In = "query"
)

func toPath(in In, name string) string {
	return fmt.Sprintf("%s/%s", in, name)
}

func unmarshalText(dest any, text []byte) error {
	text, _ = json.Marshal(string(text))
	return json.Unmarshal(text, dest)
}
