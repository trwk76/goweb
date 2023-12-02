package openapi

import (
	"fmt"
	"github.com/trwk76/goweb/openapi/spec"
)

type (
	Auths []*Auth

	Auth struct {
		name string
		spec spec.SecurityScheme
	}
)

func (a Auths) Name(name string) *Auth {
	for _, itm := range a {
		if itm.name == name {
			return itm
		}
	}

	return nil
}

func (a *Auths) Add(auth Auth) {
	if a.Name(auth.name) != nil {
		panic(fmt.Errorf("another authentication is already registered as '%s'", auth.name))
	}

	*a = append(*a, &auth)
}

func (a Auth) Name() string {
	return a.name
}
