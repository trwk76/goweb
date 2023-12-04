package spec_test

import (
	"reflect"
	"testing"

	"github.com/trwk76/goweb/openapi/spec"
)

func TestOpenAPI(t *testing.T) {
	for _, tst := range apiTests {
		tst.test(t)
	}
}

type apiTest struct {
	api  spec.OpenAPI
	yaml string
}

func (a apiTest) test(t *testing.T) {
	raw := string(a.api.YAML())

	if raw != a.yaml {
		t.Errorf("serialization failed:\n%s\n\n%s\n\n", raw, a.yaml)
		return
	}

	api, err := spec.ParseYAML([]byte(raw))
	if err != nil {
		t.Errorf("parsing failed: %s", err.Error())
		return
	}

	if !reflect.DeepEqual(a.api, api) {
		t.Errorf("APIs are not equal")
		return
	}
}

var apiTests []apiTest = []apiTest{
	{
		api: spec.OpenAPI{
			OpenAPI: spec.Version,
			Info: spec.Info{
				Title:       "TestAPI",
				Version:     "1.0.0",
				Summary:     "API for testing",
				Description: "API for testing purpose only",
			},
		},
		yaml: `openapi: 3.0.3
info:
    title: TestAPI
    version: 1.0.0
    summary: API for testing
    description: API for testing purpose only
`,
	},
}
