package gen

import (
	"fmt"
	"go/token"

	"github.com/trwk76/goweb/openapi/spec"
)

func (a *API) APIKeySecurity(key string, desc string, name string, in spec.SecurityIn) {
	a.addAuth(key, spec.SecurityScheme{
		Type:        spec.SecurityAPIKey,
		Description: desc,
		Name:        name,
		In:          in,
	})
}

func (a *API) HTTPBasicSecurity(key string, desc string) {
	a.addAuth(key, spec.SecurityScheme{
		Type:        spec.SecurityHTTP,
		Description: desc,
		Scheme:      "Basic",
	})
}

func (a *API) addAuth(key string, spc spec.SecurityScheme) {
	if !token.IsIdentifier(key) {
		panic(fmt.Errorf("'%s' is not a valid identifier", key))
	}

	if _, fnd := a.spc.Components.SecuritySchemes[key]; fnd {
		panic(fmt.Errorf("security scheme '%s' is already declared", key))
	}

	a.spc.Components.SecuritySchemes[key] = spec.SecuritySchemeOrRef{Item: spc}
}
