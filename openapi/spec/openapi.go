package spec

type (
	OpenAPI struct {
		OpenAPI           string               `json:"openapi" yaml:"openapi"`
		Info              Info                 `json:"info" yaml:"info"`
		JSONSchemaDialect string               `json:"jsonSchemaDialect,omitempty" yaml:"jsonSchemaDialect,omitempty"`
		Servers           Servers              `json:"servers,omitempty" yaml:"servers,omitempty"`
		Paths             Paths                `json:"paths,omitempty" yaml:"paths,omitempty"`
		Components        *Components          `json:"components,omitempty" yaml:"components,omitempty"`
		Security          SecurityRequirements `json:"security,omitempty" yaml:"security,omitempty"`
		Tags              Tags                 `json:"tags,omitempty" yaml:"tags,omitempty"`
		ExternalDocs      *ExternalDoc         `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	}

	Info struct {
		Title          string   `json:"title" yaml:"title"`
		Version        string   `json:"version" yaml:"version"`
		Summary        string   `json:"summary,omitempty" yaml:"summary,omitempty"`
		Description    string   `json:"description,omitempty" yaml:"description,omitempty"`
		TermsOfService string   `json:"termsOfService,omitempty" yaml:"termsOfService,omitempty"`
		Contact        *Contact `json:"contact,omitempty" yaml:"contact,omitempty"`
		License        *License `json:"license,omitempty" yaml:"license,omitempty"`
	}

	Contact struct {
		Name  string `json:"name,omitempty" yaml:"name,omitempty"`
		URL   string `json:"url,omitempty" yaml:"url,omitempty"`
		Email string `json:"email,omitempty" yaml:"email,omitempty"`
	}

	License struct {
		Name       string `json:"name" yaml:"name"`
		Identifier string `json:"identifier,omitempty" yaml:"identifier,omitempty"`
		URL        string `json:"url,omitempty" yaml:"url,omitempty"`
	}

	Servers []Server

	Server struct {
		URL         string          `json:"url" yaml:"url"`
		Description string          `json:"description,omitempty" yaml:"description,omitempty"`
		Variables   ServerVariables `json:"variables,omitempty" yanl:"variables,omitempty"`
	}

	ServerVariables map[string]ServerVariable

	ServerVariable struct {
		Default     string   `json:"default" yaml:"default"`
		Description string   `json:"description,omitempty" yaml:"description,omitempty"`
		Enum        []string `json:"enum,omitempty" yaml:"enum,omitempty"`
	}

	Components struct {
		Schemas         NamedSchemas         `json:"schemas,omitempty" yaml:"schemas,omitempty"`
		Responses       NamedResponses       `json:"responses,omitempty" yaml:"responses,omitempty"`
		Parameters      NamedParameters      `json:"parameters,omitempty" yaml:"parameters,omitempty"`
		Examples        NamedExamples        `json:"examples,omitempty" yaml:"examples,omitempty"`
		RequestBodies   NamedRequestBodies   `json:"requestBodies,omitempty" yaml:"requestBodies,omitempty"`
		Headers         NamedHeaders         `json:"headers,omitempty" yaml:"headers,omitempty"`
		SecuritySchemes NamedSecuritySchemes `json:"securitySchemes,omitempty" yaml:"securitySchemes,omitempty"`
	}

	NamedSchemas         map[string]Schema
	NamedResponses       map[string]Response
	NamedParameters      map[string]Parameter
	NamedExamples        map[string]Example
	NamedRequestBodies   map[string]RequestBody
	NamedHeaders         map[string]Header
	NamedSecuritySchemes map[string]SecurityScheme

	Tags []Tag

	Tag struct {
		Name         string       `json:"name" yaml:"name"`
		Description  string       `json:"description,omitempty" yaml:"description,omitempty"`
		ExternalDocs *ExternalDoc `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	}

	ExternalDoc struct {
		URL         string `json:"url" yaml:"url"`
		Description string `json:"description,omitempty" yaml:"description,omitempty"`
	}

	Example struct {
		Ref           string `json:"$ref,omitempty" yaml:"$ref,omitempty"`
		Summary       string `json:"summary,omitempty" yaml:"summary,omitempty"`
		Description   string `json:"description,omitempty" yaml:"description,omitempty"`
		Value         any    `json:"value,omitempty" yaml:"value,omitempty"`
		ExternalValue string `json:"externalValue,omitempty" yaml:"externalValue,omitempty"`
	}
)

const Version string = "3.0.3"
