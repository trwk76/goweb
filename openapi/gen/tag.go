package gen

import (
	"fmt"
	"slices"

	"github.com/trwk76/goweb/openapi/spec"
)

func (t *Tags) Add(name string, desc string) {
	if slices.ContainsFunc(*t, func(itm spec.Tag) bool { return itm.Name == name }) {
		panic(fmt.Errorf("tag '%s' redefined", name))
	}

	*t = append(*t, spec.Tag{Name: name, Description: desc})
}

type (
	Tags []spec.Tag
)
