package marshal

import (
	"encoding/json"
	"fmt"
	"math"
)

func MarshalInt[T integer](value T) ([]byte, error) {
	return json.Marshal(int64(value))
}

func UnmarshalInt[T integer](dest *T, raw []byte, min T, max T, multipleOf T) error {
	var tmp int64

	if err := json.Unmarshal(raw, &tmp); err != nil {
		return fmt.Errorf("integer value expected")
	}

	if tmp < int64(min) {
		return fmt.Errorf("value is smaller than %d (value: %d)", min, tmp)
	} else if tmp > int64(max) {
		return fmt.Errorf("value is greater than %d (value: %d)", max, tmp)
	} else if multipleOf != 0 && tmp%int64(multipleOf) != 0 {
		return fmt.Errorf("value is not a multiple of %d (value: %d)", multipleOf, tmp)
	}

	*dest = T(tmp)
	return nil
}

func MarshalUint[T uinteger](value T) ([]byte, error) {
	return json.Marshal(uint64(value))
}

func UnmarshalUint[T uinteger](dest *T, raw []byte, min T, max T, multipleOf T) error {
	var tmp uint64

	if err := json.Unmarshal(raw, &tmp); err != nil {
		return fmt.Errorf("unsigned integer value expected")
	}

	if tmp < uint64(min) {
		return fmt.Errorf("value is smaller than %d (value: %d)", min, tmp)
	} else if tmp > uint64(max) {
		return fmt.Errorf("value is greater than %d (value: %d)", max, tmp)
	} else if multipleOf != 0 && tmp%uint64(multipleOf) != 0 {
		return fmt.Errorf("value is not a multiple of %d (value: %d)", multipleOf, tmp)
	}

	*dest = T(tmp)
	return nil
}

func MarshalFloat[T float](value T) ([]byte, error) {
	return json.Marshal(float64(value))
}

func UnmarshalFloat[T float](dest *T, raw []byte, min T, max T, multipleOf T) error {
	var tmp float64

	if err := json.Unmarshal(raw, &tmp); err != nil {
		return fmt.Errorf("numeric value expected")
	}

	if tmp < float64(min) {
		return fmt.Errorf("value is smaller than %g (value: %g)", min, tmp)
	} else if tmp > float64(max) {
		return fmt.Errorf("value is greater than %g (value: %g)", max, tmp)
	} else if multipleOf != 0 && math.Remainder(tmp, float64(multipleOf)) != 0 {
		return fmt.Errorf("value is not a multiple of %g (value: %g)", multipleOf, tmp)
	}

	*dest = T(tmp)
	return nil
}

type (
	integer interface {
		~int8 | ~int16 | ~int32 | ~int64 | ~int
	}

	uinteger interface {
		~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uint
	}

	float interface {
		~float32 | ~float64
	}
)
