package openapi

import "reflect"

type (
	Responses map[int]Response

	Response struct {
		desc string
		typ  reflect.Type
	}
)

func (r Responses) Default() Response {
	return r[0]
}

func (r Responses) Set(status int, resp Response) {
	r[status] = resp
}

func newResponses() Responses {
	return make(Responses)
}
