package spec

type (
	SecurityRequirements []SecurityRequirement

	SecurityRequirement map[string][]string

	NamedSecuritySchemes map[string]SecurityScheme

	SecurityScheme struct {
		Type         SecurityType       `json:"type" yaml:"type"`
		Description  string             `json:"description,omitempty" yaml:"description,omitempty"`
		Name         string             `json:"name,omitempty" yaml:"name,omitempty"`
		In           SecurityAPIKeyIn   `json:"in,omitempty" yaml:"in,omitempty"`
		Scheme       SecurityHTTPScheme `json:"scheme,omitempty" yaml:"scheme,omitempty"`
		BearerFormat string             `json:"bearerFormat,omitempty" yaml:"bearerFormat,omitempty"`
	}

	SecurityType       string
	SecurityAPIKeyIn   string
	SecurityHTTPScheme string
)

const (
	SecurityTypeAPIKey        SecurityType = "apiKey"
	SecurityTypeHTTP          SecurityType = "http"
	SecurityTypeOAuth2        SecurityType = "oauth2"
	SecurityTypeOpenIDConnect SecurityType = "openIdConnect"
)

const (
	SecurityAPIKeyInCookie SecurityAPIKeyIn = "cookie"
	SecurityAPIKeyInHeader SecurityAPIKeyIn = "header"
	SecurityAPIKeyInQuery  SecurityAPIKeyIn = "query"
)

const (
	SecurityHTTPSchemeBasic SecurityHTTPScheme = "Basic"
)

var AnonymousSecurity SecurityRequirements = SecurityRequirements{{}}
