package openapi

import (
	"reflect"

	"github.com/trwk76/goweb"
)

func Response[T any](description string, examples map[string]Example[T]) ResponseSpec {
	return ResponseSpec{
		desc: description,
		t:    typeOf[T](),
		e:    Examples(examples),
	}
}

func DefaultResponse[T any](description string, examples map[string]Example[T], crea func(ctx *goweb.Context, status int) T) DefaultResponseSpec {
	return DefaultResponseSpec{
		crea: func(ctx *goweb.Context, status int) reflect.Value {
			return reflect.ValueOf(crea(ctx, status))
		},
		resp: Response(description, examples),
	}
}

type (
	Responses map[int]*ResponseSpec

	ResponseSpec struct {
		desc string
		t    reflect.Type
		e    ExampleSpecs
	}

	DefaultResponseSpec struct {
		crea func(ctx *goweb.Context, status int) reflect.Value
		resp ResponseSpec
	}
)

func (r Responses) Default() *ResponseSpec {
	return r[0]
}
