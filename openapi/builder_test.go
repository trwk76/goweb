package openapi_test

import (
	"testing"

	"github.com/trwk76/goweb"
	"github.com/trwk76/goweb/openapi"
	"github.com/trwk76/goweb/openapi/spec"
)

func TestBuilder(t *testing.T) {
	bs := openapi.HTTPSecurity("basic", "Basic HTTP authentication.", "Basic", "")
	ks := openapi.APIKeySecurity("key", "API Key authentication.", spec.SecurityInHeader, "X-Token")

	b := openapi.New("/api/v1", spec.Info{
		Title:   "TestAPI",
		Version: "1.0",
	}, spec.Tags{
		{
			Name:        "Users",
			Description: "User management methods",
		},
	}, openapi.NewMediaTypes(),
		openapi.NewSecuritySchemes(&bs, &ks),
		openapi.SecurityReqs{{&ks: nil}},
	)

	b.Named("user", openapi.PathOpts{
		Tags: []string{"Users"},
	}, func(p *openapi.PathBuilder) {
		openapi.GET(p, openapi.OperationOpts{
			ID:          "ListUsers",
			Summary:     "Lists users",
			Description: "Lists all users matching the specified criteria.",
		},
			func(ctx *goweb.Context, req UserListReq, res *UserListRes) {

			},
		)

		p.Param("id", openapi.PathOpts{
			Summary:     "User specific methods",
			Description: "Methods specific to an identified user.",
		}, func(p *openapi.PathBuilder) {
			openapi.POST(p, openapi.OperationOpts{
				ID:          "UpdateUser",
				Summary:     "Update user",
				Description: "Updates the user with the provided data.",
			},
				func(ctx *goweb.Context, req UserUpdateReq, res *UserUpdateRes) {

				},
			)
		})
	})
}

type (
	Page[T any] struct {
		Count     uint64
		PageSize  uint32
		PageIndex uint32
		Items     []T
	}

	Error struct {
		ID      string `json:"id" yaml:"id"`
		Status  int    `json:"status" yaml:"status"`
		Message string `json:"message" yaml:"message"`
	}

	UserData struct {
		Name    string
		Country string
	}

	UserInfo struct {
		ID int64
		UserData
	}

	InternalErrorRes struct {
		BadRequest *Error `openapi:"500"`
	}

	BadRequestRes struct {
		BadRequest *Error `openapi:"400"`
	}

	UnauthorizedRes struct {
		Unauthorized *Error `openapi:"401"`
	}

	NotFoundRes struct {
		NotFound *Error `openapi:"404"`
	}

	UserListReq struct {
		Search    string `openapi:"search@query"`
		PageSize  uint32 `openapi:"page_size@query"`
		PageIndex uint32 `openapi:"page_index@query"`
	}

	UserListRes struct {
		InternalErrorRes
		BadRequestRes
		UnauthorizedRes
		OK *Page[UserInfo] `openapi:"200"`
	}

	UserUpdateReq struct {
		ID   int64    `openapi:"id@path"`
		Body UserData `openapi:"body"`
	}

	UserUpdateRes struct {
		InternalErrorRes
		BadRequestRes
		UnauthorizedRes
		NotFoundRes
		OK *Page[UserInfo] `openapi:"200"`
	}
)
