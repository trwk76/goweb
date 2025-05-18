package content

import "mime"

type (
	// Type represents a Content Type and its parameters
	Type struct {
		Type   string
		Params map[string]string
	}
)

func ParseType(v string) (Type, error) {
	t, params, err := mime.ParseMediaType(v)
	if err != nil {
		return Type{}, err
	}

	return Type{Type: t, Params: params}, nil
}

func (v Type) String() string {
	return mime.FormatMediaType(v.Type, v.Params)
}

const (
	// TypeApplicationJSON the value of the JSON content type.
	TypeApplicationJSON string = "application/json"
	// TypeApplicationFormURL the value of the form-urlencode content type.
	TypeApplicationFormURL string = "application/x-www-form-urlencoded"
)
