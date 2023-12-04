package openapi

import (
	"fmt"

	"github.com/trwk76/goweb/openapi/spec"
)

type (
	SecuritySchemes []*SecurityScheme

	SecurityScheme struct {
		name string
		spec spec.SecurityScheme
	}

	SecurityRequirements []SecurityRequirement

	SecurityRequirement map[string][]string
)

func (s SecuritySchemes) Name(name string) *SecurityScheme {
	for _, itm := range s {
		if itm.name == name {
			return itm
		}
	}

	return nil
}

func (s *SecuritySchemes) Add(item SecurityScheme) error {
	if s.Name(item.name) != nil {
		return fmt.Errorf("another security scheme is already registered as '%s'", item.name)
	}

	*s = append(*s, &item)
	return nil
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
	AllowAnonymousRequirement SecurityRequirements = SecurityRequirements{{}}
)
