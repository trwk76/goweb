package openapi

import "github.com/trwk76/goweb/openapi/spec"

func New(basePath string, info spec.Info, defaultResponse DefaultResponseSpec) *API {
	return &API{
		path:            basePath,
		info:            info,
		defResp:         defaultResponse,
		Responses:       make(Responses),
		SecuritySchemes: nil,
		DefaultSecurity: AllowAnonymousRequirement,
		Root:            newPathSpec(),
	}
}

type (
	API struct {
		path    string
		info    spec.Info
		defResp DefaultResponseSpec

		Responses       Responses
		SecuritySchemes SecuritySchemes
		DefaultSecurity SecurityRequirements
		Tags            spec.Tags
		Root            PathSpec
	}
)
