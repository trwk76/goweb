package spec

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

type (
	NamedSchemas map[string]Schema

	Schema struct {
		Ref                  string                `json:"$ref,omitempty" yaml:"$ref,omitempty"`
		Title                string                `json:"title,omitempty" yaml:"title,omitempty"`
		Description          string                `json:"description,omitempty" yaml:"description,omitempty"`
		Type                 Type                  `json:"type,omitempty" yaml:"type,omitempty"`
		Format               Format                `json:"format,omitempty" yaml:"format,omitempty"`
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
		AdditionalProperties *AdditionalProperties `json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty"`
		PatternProperties    NamedSchemas          `json:"patternProperties,omitempty" yaml:"patternProperties,omitempty"`
		MinProperties        uint32                `json:"minProperties,omitempty" yaml:"minProperties,omitempty"`
		MaxProperties        uint32                `json:"maxProperties,omitempty" yaml:"maxProperties,omitempty"`
	}

	Type   string
	Format string

	AdditionalProperties struct {
		Schema *Schema
		Bool   bool
	}
)

const (
	TypeBoolean Type = "boolean"
	TypeInteger Type = "integer"
	TypeNumber  Type = "number"
	TypeString  Type = "string"
	TypeArray   Type = "array"
	TypeObject  Type = "object"
)

const (
	FormatNone     Format = ""
	FormatInt8     Format = "int8"
	FormatInt16    Format = "int16"
	FormatInt32    Format = "int32"
	FormatInt64    Format = "int64"
	FormatUint8    Format = "uint8"
	FormatUint16   Format = "uint16"
	FormatUint32   Format = "uint32"
	FormatUint64   Format = "uint64"
	FormatFloat32  Format = "float"
	FormatFloat64  Format = "double"
	FormatDate     Format = "date"
	FormatDateTime Format = "date-time"
	FormatByte     Format = "byte"
	FormatBinary   Format = "binary"
)

func (a AdditionalProperties) MarshalJSON() ([]byte, error) {
	if a.Schema != nil {
		return json.Marshal(a.Schema)
	}

	return json.Marshal(a.Bool)
}

func (a AdditionalProperties) MarshalYAML() (any, error) {
	var node yaml.Node
	var err error

	if a.Schema != nil {
		err = node.Encode(a.Schema)
	} else {
		err = node.Encode(a.Bool)
	}

	return node, err
}

func (a *AdditionalProperties) UnmarshalJSON(raw []byte) error {
	sch := &Schema{}

	if err := json.Unmarshal(raw, sch); err == nil {
		a.Schema = sch
		return nil
	}

	a.Schema = nil
	return json.Unmarshal(raw, &a.Bool)
}

func (a *AdditionalProperties) UnmarshalYAML(node *yaml.Node) error {
	sch := &Schema{}

	if node.Decode(sch) == nil {
		a.Schema = sch
		return nil
	}

	a.Schema = nil
	return node.Decode(&a.Bool)
}

var (
	_ json.Marshaler   = AdditionalProperties{}
	_ yaml.Marshaler   = AdditionalProperties{}
	_ json.Unmarshaler = (*AdditionalProperties)(nil)
	_ yaml.Unmarshaler = (*AdditionalProperties)(nil)
)
