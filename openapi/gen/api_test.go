package gen_test

import (
	"testing"

	"github.com/trwk76/goweb"
	"github.com/trwk76/goweb/openapi/gen"
	"github.com/trwk76/goweb/openapi/spec"
)

func TestAPI(t *testing.T) {
	api := gen.New(
		"/api/v1",
		spec.Info{
			Title:       "TestAPI 1",
			Version:     "1.0.0",
			Description: "API for testing purpose",
		},
		spec.SecurityRequirements{{apiKeySec: []string{}}},
		gen.SchemaTypeCamelCase,
	)

	api.Tags().Add("Authentication", "Token related methods.")
	api.Tags().Add("Token", "Token related methods.")

	resp500 := api.Resps().Add("respError", gen.RespFor[Error](&api, "An error occurred.", func(a *gen.API, item *spec.Response) {
		item.Content[goweb.ContentTypeJSON].Examples = spec.NamedExampleOrRefs{
			"500": spec.ExampleOrRef{Item: spec.Example{
				Summary:     "An error occurred",
				Description: "",
				Value: Error{
					RespBase: RespBase{
						Status: 500,
					},
					Message: "An unexpected error occurred.",
				},
			}},
		}
	}))

	api.NamedPath("token", func(a *gen.API, item *gen.PathItem) {
		item.PUT("tokenCreate", func(a *gen.API, item *spec.Operation) {
			item.Summary = "Create Token"
			item.Description = "Creates a token with the provided HTTP basic credentials."
			item.Tags = []string{"Authentication", "Token"}
			item.Security = spec.SecurityRequirements{{basicSec: []string{}}}
			item.Responses = spec.Responses{
				"201": gen.RespFor[TokenResponse](&api, "Token created.", func(api *gen.API, item *spec.Response) {

				}),
				"500": resp500,
			}
		})
	})

	t.Log(string(api.Spec().JSON()))
}

type (
	RespBase struct {
		Status int `json:"status"`
	}

	Error struct {
		RespBase
		Message string `json:"message"`
	}

	TokenResponse struct {
		RespBase
	}
)

const (
	apiKeySec string = "apiKey"
	basicSec  string = "basic"
)
