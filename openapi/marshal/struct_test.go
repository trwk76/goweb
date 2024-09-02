package marshal_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/trwk76/goweb/openapi/marshal"
)

func TestStruct(t *testing.T) {
	var val Record

	if err := json.Unmarshal([]byte(`{"label": null}`), &val); err == nil {
		t.Error(fmt.Errorf("property color is required"))
	}

	if err := json.Unmarshal([]byte(`{"label": null, "color": "blue"}`), &val); err != nil {
		t.Error(err)
	}

	if err := json.Unmarshal([]byte(`{"label": null, "time": "1976-04-30T13:31:00.123456789+02:00", "color": "blue"}`), &val); err != nil {
		t.Error(err)
	}
}

type Record struct {
	Label *string           `json:"label"`
	Color Color             `json:"color"`
	Time  *marshal.DateTime `json:"time,omitempty"`
}

func (v *Record) UnmarshalJSON(raw []byte) error {
	return marshal.UnmarshalStruct(v, raw, nil)
}
