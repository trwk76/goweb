package openapi

import "reflect"

type (
	Responses map[int]Response

	Response struct {
		desc string
		typ  reflect.Type
	}
)
