package marshal

import (
	"encoding/json"
	"fmt"
	"slices"
)

func MarshalArrayJSON[T ~[]I, I any](items T) ([]byte, error) {
	return json.Marshal(items)
}

func UnmarshalArrayJSON[T ~[]I, I any](dest *T, raw []byte, minLength uint32, maxLength uint32) error {
	var items []json.RawMessage

	if err := json.Unmarshal(raw, &items); err != nil {
		return fmt.Errorf("array expected")
	}

	l := len(items)
	if l < int(minLength) {
		return fmt.Errorf("array has less than %d item(s) (%d item(s) provided)", minLength, l)
	} else if maxLength > 0 && l > int(maxLength) {
		return fmt.Errorf("array has more than %d item(s) (%d item(s) provided)", maxLength, l)
	}

	res := make([]I, l)

	for idx, itm := range items {
		var item I

		if err := json.Unmarshal(itm, &item); err != nil {
			return NewPathErrors(PathIndex(idx), err)
		}

		res[idx] = item
	}

	*dest = T(res)
	return nil
}

func UnmarshalArrayUniqueJSON[T ~[]I, I comparable](dest *T, raw []byte, minLength uint32, maxLength uint32) error {
	var tmp T

	if err := UnmarshalArrayJSON(&tmp, raw, minLength, maxLength); err != nil {
		return err
	}

	for idx, itm := range tmp {
		if slices.Index(tmp, itm) != idx {
			return NewPathErrors(PathIndex(idx), fmt.Errorf("item appears more than once"))
		}
	}

	*dest = T(tmp)
	return nil
}
