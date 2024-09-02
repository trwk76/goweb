package marshal

import (
	"encoding"
	"encoding/json"
	"fmt"

	"github.com/trwk76/goweb/openapi/spec"
)

type Bool bool

func (Bool) Schema() spec.Schema {
	return spec.Schema{Type: spec.TypeBoolean}
}

func (v Bool) MarshalText() ([]byte, error) {
	return json.Marshal(bool(v))
}

func (v Bool) MarshalJSON() ([]byte, error) {
	return json.Marshal(bool(v))
}

func (v *Bool) UnmarshalText(raw []byte) error {
	var tmp bool

	if err := json.Unmarshal(raw, &tmp); err != nil {
		return fmt.Errorf("boolean value expected (either 'true' or 'false')")
	}

	*v = Bool(tmp)
	return nil
}

func (v *Bool) UnmarshalJSON(raw []byte) error {
	var tmp bool

	if err := json.Unmarshal(raw, &tmp); err != nil {
		return fmt.Errorf("boolean value expected (either 'true' or 'false')")
	}

	*v = Bool(tmp)
	return nil
}

var (
	_ encoding.TextMarshaler   = Bool(false)
	_ json.Marshaler           = Bool(false)
	_ encoding.TextUnmarshaler = (*Bool)(nil)
	_ json.Unmarshaler         = (*Bool)(nil)
)
