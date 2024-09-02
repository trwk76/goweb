package gen

import (
	"fmt"

	"github.com/trwk76/goweb/openapi/spec"
)

func (s *Securities) AddAPIKey(key string, name string, in spec.SecurityIn, setup func(s *spec.SecurityScheme)) {
	item := spec.SecurityScheme{
		Type: spec.SecurityAPIKey,
		Name: name,
		In:   in,
	}

	if setup != nil {
		setup(&item)
	}

	s.add(key, item)
}

func (s *Securities) AddHTTP(key string, scheme string, setup func(s *spec.SecurityScheme)) {
	item := spec.SecurityScheme{
		Type:   spec.SecurityHTTP,
		Scheme: scheme,
	}

	if setup != nil {
		setup(&item)
	}

	s.add(key, item)
}

func (s *Securities) AddOAuth2(key string, flows spec.OAuthFlows, setup func(s *spec.SecurityScheme)) {
	item := spec.SecurityScheme{
		Type:  spec.SecurityOAuth2,
		Flows: &flows,
	}

	if setup != nil {
		setup(&item)
	}

	s.add(key, item)
}

func (s *Securities) AddOpenIDConnect(key string, openIDConnectURL string, setup func(s *spec.SecurityScheme)) {
	item := spec.SecurityScheme{
		Type:             spec.SecurityOpenIDConnect,
		OpenIDConnectURL: openIDConnectURL,
	}

	if setup != nil {
		setup(&item)
	}

	s.add(key, item)
}

type (
	Securities spec.NamedSecuritySchemeOrRefs
)

func (s *Securities) add(key string, item spec.SecurityScheme) {
	if _, fnd := (*s)[key]; fnd {
		panic(fmt.Errorf("security scheme '%s' redefined", key))
	}

	(*s)[key] = spec.SecuritySchemeOrRef{Item: item}
}

func (s Securities) spec() spec.NamedSecuritySchemeOrRefs {
	return spec.NamedSecuritySchemeOrRefs(s)
}
