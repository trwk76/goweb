package gen

import (
	"fmt"
	"slices"

	"github.com/trwk76/goweb/openapi/spec"
)

func (a *API) Tag(name string, desc string) Tag {
	if slices.ContainsFunc(a.spc.Tags, func(t spec.Tag) bool { return t.Name == name }) {
		panic(fmt.Errorf("api already declares a tag named '%s'", name))
	}

	a.spc.Tags = append(a.spc.Tags, spec.Tag{
		Name:        name,
		Description: desc,
	})

	return Tag{name: name}
}

type (
	Tags []Tag

	Tag struct {
		name string
	}
)

func (t Tags) spec() []string {
	res := make([]string, 0, len(t))

	for _, tag := range t {
		if !slices.Contains(res, tag.name) {
			res = append(res, tag.name)
		}
	}

	return res
}
