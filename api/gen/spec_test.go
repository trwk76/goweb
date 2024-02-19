package gen_test

import (
	"testing"

	"github.com/trwk76/goweb/api/gen"
)

func TestSpec(t *testing.T) {
	a := gen.New("/api/test", "tests", "api")

	if err := a.Generate(); err != nil {
		t.Error(err)
	}
}
