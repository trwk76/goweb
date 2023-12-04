package openapi

import "reflect"

func typeOf[T any]() reflect.Type {
	var d *T
	return reflect.TypeOf(d).Elem()
}
