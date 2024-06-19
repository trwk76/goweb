package gen_test

import (
	"testing"

	"github.com/trwk76/goweb/openapi/gen"
	"github.com/trwk76/goweb/openapi/spec"
)

func TestAPI(t *testing.T) {
	// mtypes := gen.MediaTypes{web.ContentTypeJSON: gen.JSON, web.ContentTypeYAML: gen.YAML}

	api := gen.New("/api/v1", spec.Info{Title: "TestAPI", Version: "1.0.0"})
	tagSec := api.Tag("Security", "Security related methods")
	tagUsr := api.Tag("User", "User related methods")

	api.NamedPath("security", func(p *gen.Path) {
		p.OpPrefix = "Security"
		p.Tags = append(p.Tags, tagSec)
	})

	api.NamedPath("user", func(p *gen.Path) {
		p.OpPrefix = "User"
		p.Tags = append(p.Tags, tagUsr)
	})
}
