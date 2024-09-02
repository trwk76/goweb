package marshal

import (
	"fmt"
	"regexp"
	"strconv"
)

func MarshalString[T ~string](value T) ([]byte, error) {
	return []byte(string(value)), nil
}

func UnmarshalString[T ~string](dest *T, raw []byte, minLength uint32, maxLength uint32, regex *regexp.Regexp) error {
	tmp := string(raw)
	l := len(tmp)

	if l < int(minLength) {
		return fmt.Errorf("value is shorter than %d character(s) (value: %s)", minLength, strconv.Quote(tmp))
	} else if maxLength > 0 && l > int(maxLength) {
		return fmt.Errorf("value is longer than %d character(s) (value: %s)", maxLength, strconv.Quote(tmp))
	}

	if regex != nil && regex.FindString(tmp) != tmp {
		return fmt.Errorf("value does not match pattern %s (value: %s)", strconv.Quote(tmp), strconv.Quote(tmp))
	}

	*dest = T(tmp)
	return nil
}
