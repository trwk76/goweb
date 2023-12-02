package spec

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

type (
	Schema struct {
		Ref                  string                `json:"$ref,omitempty" yaml:"$ref,omitempty"`
		Title                string                `json:"title,omitempty" yaml:"title,omitempty"`
		Description          string                `json:"description,omitempty" yaml:"description,omitempty"`
		OneOf                []Schema              `json:"oneOf,omitempty" yaml:"oneOf,omitempty"`
		AllOf                []Schema              `json:"allOf,omitempty" yaml:"allOf,omitempty"`
		Type                 SchemaType            `json:"type,omitempty" yaml:"type,omitempty"`
		Format               SchemaFormat          `json:"format,omitempty" yaml:"format,omitempty"`
		Nullable             bool                  `json:"nullable,omitempty" yaml:"nullable,omitempty"`
		Enum                 []any                 `json:"enum,omitempty" yaml:"enum,omitempty"`
		Minimum              *float64              `json:"minimum,omitempty" yaml:"minimum,omitempty"`
		ExclusiveMinimum     bool                  `json:"exclusiveMinimum,omitempty" yaml:"exclusiveMinimum,omitempty"`
		Maximum              *float64              `json:"maximum,omitempty" yaml:"maximum,omitempty"`
		ExclusiveMaximum     bool                  `json:"exclusiveMaximum,omitempty" yaml:"exclusiveMaximum,omitempty"`
		MultipleOf           float64               `json:"multipleOf,omitempty" yaml:"multipleOf,omitempty"`
		MinLength            uint32                `json:"minLength,omitempty" yaml:"minLength,omitempty"`
		MaxLength            uint32                `json:"maxLength,omitempty" yaml:"maxLength,omitempty"`
		Pattern              string                `json:"pattern,omitempty" yaml:"pattern,omitempty"`
		Items                *Schema               `json:"items,omitempty" yaml:"items,omitempty"`
		MinItems             uint32                `json:"minItems,omitempty" yaml:"minItems,omitempty"`
		MaxItems             uint32                `json:"maxItems,omitempty" yaml:"maxItems,omitempty"`
		UniqueItems          bool                  `json:"uniqueItems,omitempty" yaml:"uniqueItems,omitempty"`
		Required             []string              `json:"required,omitempty" yaml:"required,omitempty"`
		Properties           NamedSchemas          `json:"properties,omitempty" yaml:"properties,omitempty"`
		MinProperties        uint32                `json:"minProperties,omitempty" yaml:"minProperties,omitempty"`
		MaxProperties        uint32                `json:"maxProperties,omitempty" yaml:"maxProperties,omitempty"`
		AdditionalProperties *AdditionalProperties `json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty"`
		PatternProperties    NamedSchemas          `json:"patternProperties,omitempty" yaml:"patternProperties,omitempty"`
		Discriminator        *Discriminator        `json:"discriminator,omitempty" yaml:"discriminator,omitempty"`
		XML                  *XML                  `json:"xml,omitempty" yaml:"xml,omitempty"`
		Example              any                   `json:"example,omitempty" yaml:"example,omitempty"`
		Examples             NamedExamples         `json:"examples,omitempty" yaml:"examples,omitempty"`
	}

	SchemaType   string
	SchemaFormat string

	AdditionalProperties struct {
		Schema *Schema
		Bool   bool
	}

	Discriminator struct {
		PropertyName string            `json:"propertyName" yaml:"propertyName"`
		Mapping      map[string]string `json:"mapping,omitempty" yaml:"mapping,omitempty"`
	}

	XML struct {
		Name      string `json:"name,omitempty" yaml:"name,omitempty"`
		Namespace string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
		Prefix    string `json:"prefix,omitempty" yaml:"prefix,omitempty"`
		Attribute bool   `json:"attribute,omitempty" yaml:"attribute,omitempty"`
		Wrapped   bool   `json:"wrapped,omitempty" yaml:"wrapped,omitempty"`
	}
)

const (
	SchemaTypeBoolean SchemaType = "boolean"
	SchemaTypeInteger SchemaType = "integer"
	SchemaTypeNumber  SchemaType = "number"
	SchemaTypeString  SchemaType = "string"
	SchemaTypeArray   SchemaType = "array"
	SchemaTypeObject  SchemaType = "object"
)

const (
	SchemaFormatBinary              SchemaFormat = "binary"
	SchemaFormatBase64              SchemaFormat = "base64"
	SchemaFormatDate                SchemaFormat = "date"
	SchemaFormatTime                SchemaFormat = "time"
	SchemaFormatDateTime            SchemaFormat = "date-time"
	SchemaFormatDuration            SchemaFormat = "duration"
	SchemaFormatEmail               SchemaFormat = "email"
	SchemaFormatIDNEmail            SchemaFormat = "idn-email"
	SchemaFormatHostname            SchemaFormat = "hostname"
	SchemaFormatIDNHostname         SchemaFormat = "idn-hostname"
	SchemaFormatIPv4                SchemaFormat = "ipv4"
	SchemaFormatIPv6                SchemaFormat = "ipv6"
	SchemaFormatURI                 SchemaFormat = "uri"
	SchemaFormatURIReference        SchemaFormat = "uri-reference"
	SchemaFormatIRI                 SchemaFormat = "iri"
	SchemaFormatIRIReference        SchemaFormat = "iri-reference"
	SchemaFormatUUID                SchemaFormat = "uuid"
	SchemaFormatURITemplate         SchemaFormat = "uri-template"
	SchemaFormatJSONPointer         SchemaFormat = "json-pointer"
	SchemaFormatRelativeJSONPointer SchemaFormat = "relative-json-pointer"
	SchemaFormatRegex               SchemaFormat = "regex"
	SchemaFormatInt32               SchemaFormat = "int32"
	SchemaFormatInt64               SchemaFormat = "int64"
	SchemaFormatFloat               SchemaFormat = "float"
	SchemaFormatDouble              SchemaFormat = "double"
	SchemaFormatPassword            SchemaFormat = "password"
)

func (a AdditionalProperties) MarshalJSON() ([]byte, error) {
	if a.Schema != nil {
		return json.MarshalIndent(a.Schema, "", "\t")
	}

	return json.MarshalIndent(a.Bool, "", "\t")
}

func (a AdditionalProperties) MarshalYAML() (any, error) {
	var res yaml.Node
	var err error

	if a.Schema != nil {
		err = res.Encode(a.Schema)
	} else {
		err = res.Encode(a.Bool)
	}

	return res, err
}

func (a *AdditionalProperties) UnmarshalJSON(raw []byte) error {
	sch := &Schema{}

	if err := json.Unmarshal(raw, sch); err == nil {
		a.Schema = sch
		return nil
	}

	return json.Unmarshal(raw, &a.Bool)
}

func (a *AdditionalProperties) UnmarshalYAML(node *yaml.Node) error {
	sch := &Schema{}

	if err := node.Decode(sch); err == nil {
		a.Schema = sch
		return nil
	}

	return node.Decode(&a.Bool)
}

var (
	_ json.Marshaler   = AdditionalProperties{}
	_ yaml.Marshaler   = AdditionalProperties{}
	_ json.Unmarshaler = (*AdditionalProperties)(nil)
	_ yaml.Unmarshaler = (*AdditionalProperties)(nil)
)
