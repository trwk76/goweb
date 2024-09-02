package marshal

import (
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

func MarshalIntEnum[T integer](value T) ([]byte, error) {
	return json.Marshal(int64(value))
}

func UnmarshalIntEnum[T integer](dest *T, raw []byte, items ...T) error {
	var tmp int64

	if err := json.Unmarshal(raw, &tmp); err != nil {
		return fmt.Errorf("integer value expected")
	}

	if !slices.Contains(items, T(tmp)) {
		return fmt.Errorf("value is not part of the enumeration (value: %d; valid options: %s)", tmp, EnumOptionsString(items, func(item T) string { return strconv.FormatInt(int64(item), 10) }))
	}

	*dest = T(tmp)
	return nil
}

func MarshalUintEnum[T uinteger](value T) ([]byte, error) {
	return json.Marshal(uint64(value))
}

func UnmarshalUintEnum[T uinteger](dest *T, raw []byte, items ...T) error {
	var tmp uint64

	if err := json.Unmarshal(raw, &tmp); err != nil {
		return fmt.Errorf("integer value expected")
	}

	if !slices.Contains(items, T(tmp)) {
		return fmt.Errorf("value is not part of the enumeration (value: %d; valid options: %s)", tmp, EnumOptionsString(items, func(item T) string { return strconv.FormatUint(uint64(item), 10) }))
	}

	*dest = T(tmp)
	return nil
}

func MarshalFloatEnum[T integer](value T) ([]byte, error) {
	return json.Marshal(float64(value))
}

func UnmarshalFloatEnum[T float](dest *T, raw []byte, items ...T) error {
	var tmp float64

	if err := json.Unmarshal(raw, &tmp); err != nil {
		return fmt.Errorf("numeric value expected")
	}

	if !slices.Contains(items, T(tmp)) {
		return fmt.Errorf("value is not part of the enumeration (value: %g; valid options: %s)", tmp, EnumOptionsString(items, func(item T) string { return strconv.FormatFloat(float64(item), 'g', -1, 64) }))
	}

	*dest = T(tmp)
	return nil
}

func MarshalStringEnum[T ~string](value T) ([]byte, error) {
	return json.Marshal(string(value))
}

func UnmarshalStringEnum[T ~string](dest *T, raw []byte, items ...T) error {
	tmp := string(raw)

	if !slices.Contains(items, T(tmp)) {
		return fmt.Errorf("value is not part of the enumeration (value: %s; valid options: %s)", tmp, EnumOptionsString(items, func(item T) string { return strconv.Quote(string(item)) }))
	}

	*dest = T(tmp)
	return nil
}

func EnumOptionsString[T any](items []T, format func(item T) string) string {
	strs := make([]string, len(items))

	for idx, itm := range items {
		strs[idx] = format(itm)
	}

	return strings.Join(strs, ", ")
}
