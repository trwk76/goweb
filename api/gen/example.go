package gen

import "github.com/trwk76/goweb/api/spec"

type (
	Examples[T any] map[string]Example[T]

	Example[T any] struct {
		Summary     string
		Description string
		Value       T
	}
)

func (e Examples[T]) spec() spec.NamedExamples {
	res := make(spec.NamedExamples)

	for key, itm := range e {
		res[key] = itm.spec()
	}

	return res
}

func (e Example[T]) spec() spec.Example {
	return spec.Example{
		Summary:     e.Summary,
		Description: e.Description,
		Value:       e.Value,
	}
}
