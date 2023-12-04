package openapi

import "strings"

func (p *PathSpec) Sub(path string) *PathSpec {
	path = strings.TrimPrefix(path, "/")
}

type (
	PathSpec struct {
		Summary     string
		Description string
		GET         *OperationSpec
		PUT         *OperationSpec
		POST        *OperationSpec
		DELETE      *OperationSpec
		OPTIONS     *OperationSpec
		HEAD        *OperationSpec
		PATCH       *OperationSpec
		TRACE       *OperationSpec

		sub paths
	}

	OperationSpec struct {
		id   string
		summ string
		desc string
		tags []string
		depr bool
		req  RequestSpec
		res  ResponsesSpec
		sec  SecurityRequirements
	}

	paths struct {
		n map[string]*PathSpec
		p *paramPath
	}

	paramPath struct {
		s PathSpec
		n string
	}
)

func newPathSpec() PathSpec {
	return PathSpec{
		sub: newPaths(),
	}
}

func newPaths() paths {
	return paths{
		n: make(map[string]*PathSpec),
	}
}
