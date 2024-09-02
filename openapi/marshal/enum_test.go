package marshal_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/trwk76/goweb/openapi/marshal"
)

func TestEnum(t *testing.T) {
	var val Color

	if err := json.Unmarshal([]byte(`"black"`), &val); err == nil {
		t.Error(fmt.Errorf("black is not in the enumeration"))
	}

	if err := json.Unmarshal([]byte(`"blue"`), &val); err != nil {
		t.Error(err)
	}
}

type Color string

const (
	Red   Color = "red"
	Green Color = "green"
	Blue  Color = "blue"
)

func (v *Color) UnmarshalText(raw []byte) error {
	return marshal.UnmarshalStringEnum(v, raw, Red, Green, Blue)
}
