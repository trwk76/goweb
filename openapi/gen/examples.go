package gen

import "github.com/trwk76/goweb/openapi/spec"

type (
	Examples[T any] map[string]Example[T]

	Example[T any] struct {
		Summary     string
		Description string
		Value       T
	}
)

func (e Examples[T]) spec(m MediaType) spec.NamedExampleOrRefs {
	res := make(spec.NamedExampleOrRefs)

	for key, itm := range e {
		res[key] = spec.ExampleOrRef{Item: itm.spec(m)}
	}

	return res
}

func (e Example[T]) spec(m MediaType) spec.Example {
	return spec.Example{
		Summary:     e.Summary,
		Description: e.Description,
		Value:       m.Example(e.Value),
	}
}
