package marshal_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/trwk76/goweb/openapi/marshal"
)

func TestArray(t *testing.T) {
	var val Colors

	if err := json.Unmarshal([]byte(`["black", "blue"]`), &val); err == nil {
		t.Error(fmt.Errorf("black is not in the enumeration"))
	}

	if err := json.Unmarshal([]byte(`["blue", "green"]`), &val); err != nil {
		t.Error(err)
	}

	if err := json.Unmarshal([]byte(`["blue", "green", "red"]`), &val); err == nil {
		t.Error(fmt.Errorf("too many items"))
	}

	if err := json.Unmarshal([]byte(`["blue", "blue"]`), &val); err == nil {
		t.Error(fmt.Errorf("not unique"))
	}
}

type Colors []Color

func (v *Colors) UnmarshalJSON(raw []byte) error {
	return marshal.UnmarshalArrayUniqueJSON(v, raw, 1, 2)
}
