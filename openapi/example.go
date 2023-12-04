package openapi

import "reflect"

func Examples[T any](items map[string]Example[T]) ExampleSpecs {
	var item *T

	itms := make(map[string]exampleSpec)

	for key, item := range items {
		itms[key] = exampleSpec{
			summ: item.Summary,
			desc: item.Description,
			val:  reflect.ValueOf(item.Value),
		}
	}

	return ExampleSpecs{
		t: reflect.TypeOf(item).Elem(),
		n: itms,
	}
}

type (
	Example[T any] struct {
		Summary     string
		Description string
		Value       T
	}

	ExampleSpecs struct {
		t reflect.Type
		n map[string]exampleSpec
	}

	exampleSpec struct {
		summ string
		desc string
		val  reflect.Value
	}
)
