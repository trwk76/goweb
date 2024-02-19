package spec

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

type (
	// Spec is an OpenAPI specification
	// [Reference]: https://spec.openapis.org/oas/v3.0.3#openapi-object
	Spec struct {
		OpenAPI    string               `json:"openapi" yaml:"openapi"`
		Info       Info                 `json:"info" yaml:"info"`
		Paths      NamedPaths           `json:"paths" yaml:"paths"`
		Components *Components          `json:"components,omitempty" yaml:"components,omitempty"`
		Servers    Servers              `json:"servers,omitempty" yaml:"servers,omitempty"`
		Security   SecurityRequirements `json:"security,omitempty" yaml:"security,omitempty"`
		Tags       Tags                 `json:"tags,omitempty" yaml:"tags,omitempty"`
	}

	Info struct {
		Title       string   `json:"title" yaml:"title"`
		Version     string   `json:"version" yaml:"version"`
		Description string   `json:"description,omitempty" yaml:"description,omitempty"`
		Contact     *Contact `json:"contact,omitempty" yaml:"contact,omitempty"`
		License     *License `json:"license,omitempty" yaml:"license,omitempty"`
	}

	Contact struct {
		Name  string `json:"name,omitempty" yaml:"name,omitempty"`
		Email string `json:"email,omitempty" yaml:"email,omitempty"`
		URL   string `json:"url,omitempty" yaml:"url,omitempty"`
	}

	License struct {
		Name string `json:"name" yaml:"name"`
		URL  string `json:"url,omitempty" yaml:"url,omitempty"`
	}

	Components struct {
		Schemas         NamedSchemas         `json:"schemas,omitempty" yaml:"schemas,omitempty"`
		Parameters      NamedParameters      `json:"parameters,omitempty" yaml:"parameters,omitempty"`
		RequestBodies   NamedRequestBodies   `json:"requestBodies,omitempty" yaml:"requestBodies,omitempty"`
		Responses       NamedResponses       `json:"responses,omitempty" yaml:"responses,omitempty"`
		Headers         NamedHeaders         `json:"headers,omitempty" yaml:"headers,omitempty"`
		SecuritySchemes NamedSecuritySchemes `json:"securitySchemes,omitempty" yaml:"securitySchemes,omitempty"`
		Examples        NamedExamples        `json:"examples,omitempty" yaml:"examples,omitempty"`
	}

	Servers []Server

	Server struct {
		URL         string                    `json:"url" yaml:"url"`
		Description string                    `json:"description,omitempty" yaml:"description,omitempty"`
		Variables   map[string]ServerVariable `json:"variables,omitempty" yaml:"variables,omitempty"`
	}

	ServerVariable struct {
		Default     string   `json:"default" yaml:"default"`
		Description string   `json:"description,omitempty" yaml:"description,omitempty"`
		Enum        []string `json:"enum,omitempty" yaml:"enum,omitempty"`
	}

	Tags []Tag

	Tag struct {
		Name        string `json:"name" yaml:"name"`
		Description string `json:"description,omitempty" yaml:"description,omitempty"`
	}
)

func ParseJSON(bytes []byte) (*Spec, error) {
	res := &Spec{}

	if err := json.Unmarshal(bytes, res); err != nil {
		return nil, err
	}

	return res, nil
}

func ParseYAML(bytes []byte) (*Spec, error) {
	res := &Spec{}

	if err := yaml.Unmarshal(bytes, res); err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Spec) JSON() []byte {
	res, _ := json.MarshalIndent(s, "", "\t")
	return res
}

func (s *Spec) YAML() []byte {
	res, _ := yaml.Marshal(s)
	return res
}

const OpenAPI string = "3.0.3"
