package content

import (
	"fmt"
	"strings"
)

type (
	// TypeError represents a mismatch between potential content types and the one actually provided
	TypeError struct {
		Act string
		Exp []string
	}
)

func (e TypeError) Error() string {
	return fmt.Sprintf("cannot handle content type %s; expected are %s", e.Act, strings.Join(e.Exp, ", "))
}

var (
	_ error = TypeError{}
)
