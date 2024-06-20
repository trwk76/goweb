package gen_test

import (
	"testing"

	web "github.com/trwk76/goweb"
	"github.com/trwk76/goweb/openapi/gen"
	"github.com/trwk76/goweb/openapi/spec"
)

func TestAPI(t *testing.T) {
	mtypes := gen.MediaTypes{web.ContentTypeJSON: gen.JSON, web.ContentTypeYAML: gen.YAML}

	api := gen.New("/api/v1", spec.Info{Title: "TestAPI", Version: "1.0.0"}, keyReqs)
	tagSec := api.Tag("Security", "Security related methods")
	tagUsr := api.Tag("User", "User related methods")

	api.APIKeySecurity(keyAuth, "API key based authentication", "apikey", spec.SecurityInCookie)

	paramID := api.ParamOrRef("entityID", gen.ParamFor(api, "id", spec.ParameterPath, func(s *spec.Parameter) {
		s.Description = "Entity identifier."
	}, gen.Examples[EntityID]{
		"System": {
			Summary:     "System",
			Description: "Identifier of the system",
			Value:       EntityID(1),
		},
	}))

	api.NamedPath("security", func(p *gen.Path) {
		p.OpPrefix = "Security"
		p.Tags = append(p.Tags, tagSec)
	})

	api.NamedPath("user", func(p *gen.Path) {
		p.OpPrefix = "User"
		p.Tags = append(p.Tags, tagUsr)

		p.ParamPath(paramID, func(p *gen.Path) {
			p.GET(spec.Operation{
				OperationID: "Fetch",
				Summary:     "Fetch user",
				Description: "Fetches a user by its ID.",
				Responses: spec.Responses{
					"200": api.ResponseOrRef("", gen.ResponseFor(api, "Information about the user.", nil, mtypes, gen.Examples[UserInfo]{
						"System": {
							Summary:     "System",
							Description: "System user information",
							Value: UserInfo{
								ID:    EntityID(1),
								Name:  "System",
								Email: "",
							},
						},
					})),
				},
			})
		})
	})

	if err := api.Generate("test", "api_gen.go", "api"); err != nil {
		t.Error(err)
	}
}

type (
	EntityID int64

	UserInfo struct {
		ID    EntityID `json:"id" yaml:"id"`
		Name  string   `json:"name" yaml:"name"`
		Email string   `json:"email,omitempty" yaml:"email,omitempty"`
	}
)

func (EntityID) Schema() spec.Schema {
	return spec.Schema{
		Type:    spec.TypeInteger,
		Format:  spec.FormatInt64,
		Minimum: 1,
	}
}

var keyReqs spec.SecurityRequirements = spec.SecurityRequirements{{keyAuth: {}}}

const (
	keyAuth string = "key"
)
