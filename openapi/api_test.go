package openapi_test

import (
	"net/http"
	"testing"

	"github.com/trwk76/goweb"
	"github.com/trwk76/goweb/openapi"
	"github.com/trwk76/goweb/openapi/spec"
)

func TestAPI(t *testing.T) {
	a := openapi.New("/api/v1", spec.Info{
		Title:   "TestAPI",
		Version: "1.0.0",
		Summary: "Test API",
		Contact: &spec.Contact{
			Name:  "Designer",
			Email: "designer@company.org",
		},
	},
		openapi.DefaultResponse(
			"An error occured.",
			map[string]openapi.Example[Error]{
				"internalError": {
					Summary:     "Internal error",
					Description: "An internal error occurred.",
					Value: Error{
						ID:      "00000000-0000-0000-0000-000000000000",
						Status:  http.StatusInternalServerError,
						Message: "An internal server error occurred.",
					},
				},
			},
			func(ctx *goweb.Context, status int) Error {
				return Error{
					ID:      ctx.ID(),
					Status:  status,
					Message: http.StatusText(status),
				}
			},
		))

	basic := openapi.HTTPSecurity("basic", "HTTP Basic authentication for creating an API key.", "Basic", "")
	key := openapi.APIKeySecurity("key", "API key used to authenticate API methods.", spec.SecurityInHeader, "X-Token")

	a.SecuritySchemes.Add(basic)
	a.SecuritySchemes.Add(key)
	a.DefaultSecurity = openapi.SecurityRequirements{{key.Name(): []string{}}}
	a.Tags = spec.Tags{
		{
			Name:        "Authentication",
			Description: "Methods for handling user authentication",
		},
		{
			Name:        "User",
			Description: "Methods for handling user management",
		},
	}
}

type (
	Error struct {
		ID      string `json:"id"`
		Status  int    `json:"status"`
		Message string `json:"message"`
	}
)
