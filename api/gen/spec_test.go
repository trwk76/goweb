package gen_test

import (
	"testing"

	web "github.com/trwk76/goweb"
	"github.com/trwk76/goweb/api/gen"
	"github.com/trwk76/goweb/api/spec"
)

func TestSpec(t *testing.T) {
	a := gen.New("/api/test", "tests", "api")

	a.Info.Title = "TestAPI"
	a.Info.Version = "1.0.0"
	a.Auths = spec.NamedSecuritySchemes{
		basicAuthKey: spec.SecurityScheme{
			Type:        spec.SecurityTypeHTTP,
			Description: "",
			Scheme:      spec.SecurityHTTPSchemeBasic,
		},
		keyAuthKey: spec.SecurityScheme{
			Type:        spec.SecurityTypeAPIKey,
			Description: "",
			Name:        keyAuthKey,
			In:          spec.SecurityAPIKeyInCookie,
		},
	}
	a.Tags = spec.Tags{
		{Name: "Security", Description: "Security related operations."},
		{Name: "User", Description: "User related operations."},
	}

	gen.Path(a, nil, "user", false, nil, []string{"User"}, func(sp *gen.PathSpec) {
		gen.Path(a, sp, "apiKey", false, nil, []string{"Security"}, func(sp *gen.PathSpec) {
			gen.PUT(sp, gen.Op[APIKeyCreateReq, APIKeyCreateRes](a, APIKeyCreate, gen.OpInfo{
				ID:       "APIKeyCreate",
				Summary:  "Create an API key.",
				Security: secReqBasic,
			}))

			gen.DELETE(sp, gen.Op(a, APIKeyDelete, gen.OpInfo{}))
		})
	})

	if err := a.Generate(); err != nil {
		t.Error(err)
	}
}

func APIKeyCreate(c *web.Context, req APIKeyCreateReq, res *APIKeyCreateRes) {

}

func APIKeyDelete(c *web.Context, req APIKeyDeleteReq, res *APIKeyDeleteRes) {

}

type (
	Void struct{}

	APIKeyCreateReq Void

	APIKeyCreateRes struct {
		Created UserInfo `api:"201"`
	}

	APIKeyDeleteReq Void

	APIKeyDeleteRes struct {
		OK Void `api:"200"`
	}

	UserInfo struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}
)

var (
	secReqBasic spec.SecurityRequirements = spec.SecurityRequirements{{basicAuthKey: {}}}
	secReqKey   spec.SecurityRequirements = spec.SecurityRequirements{{keyAuthKey: {}}}
)

const (
	basicAuthKey string = "basic"
	keyAuthKey   string = "key"
)
