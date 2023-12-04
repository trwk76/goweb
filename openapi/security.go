package openapi

import (
	"log"

	"github.com/trwk76/goweb/openapi/spec"
)

type (
	SecuritySchemes []*SecurityScheme

	SecurityScheme struct {
		name string
		spec spec.SecurityScheme
	}

	SecurityReqs []SecurityReq

	SecurityReq map[*SecurityScheme][]string
)

func NewSecuritySchemes(items ...*SecurityScheme) SecuritySchemes {
	res := make(SecuritySchemes, 0)

	for _, itm := range items {
		res = res.Add(itm)
	}

	return res
}

func (s SecuritySchemes) Name(name string) *SecurityScheme {
	for _, itm := range s {
		if itm.name == name {
			return itm
		}
	}

	return nil
}

func (s SecuritySchemes) Add(item *SecurityScheme) SecuritySchemes {
	if s.Name(item.name) != nil {
		log.Fatalf("another security scheme is already registered as '%s'", item.name)
	}

	return append(s, item)
}

func APIKeySecurity(name string, description string, keyIn spec.SecurityIn, keyName string) SecurityScheme {
	return SecurityScheme{
		name: name,
		spec: spec.SecurityScheme{
			Description: description,
			Type:        spec.SecurityTypeAPIKey,
			In:          keyIn,
			Name:        keyName,
		},
	}
}

func HTTPSecurity(name string, description string, scheme string, bearerFormat string) SecurityScheme {
	return SecurityScheme{
		name: name,
		spec: spec.SecurityScheme{
			Description:  description,
			Type:         spec.SecurityTypeHTTP,
			Scheme:       scheme,
			BearerFormat: bearerFormat,
		},
	}
}

func (s SecurityScheme) Name() string {
	return s.name
}

func (s SecurityScheme) Description() string {
	return s.spec.Description
}

var (
	AnonymousReq SecurityReqs = SecurityReqs{{}}
)
