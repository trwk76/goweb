package openapi

import (
	"net/http"

	"github.com/trwk76/goweb/openapi/spec"
)

func New[DR any](path string, info spec.Info, defResp DefaultResponder[DR]) Builder {
	return Builder{
		Path:       path,
		Info:       info,
		MediaTypes: make(MediaTypes, 0),
		Auths:      make(Auths, 0),
		Responses:  newResponses(),
	}
}

type (
	DefaultResponder[T any] func(c *Context, status int) T

	Builder struct {
		Path       string
		Info       spec.Info
		MediaTypes MediaTypes
		Auths      Auths
		Responses  Responses

		defResp func(c *Context, status int, w http.ResponseWriter)
	}

	Context struct {

	}
)
