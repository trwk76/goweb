package marshal

import (
	"encoding"
	"fmt"
	"math"
	"reflect"
	"strconv"
)

func MarshalText(value any) ([]byte, error) {
	val := reflect.ValueOf(value)

	if val.Kind() == reflect.Pointer {
		if val.IsNil() {
			return nil, nil
		}

		val = reflect.Indirect(val)
	}

	if impl, ok := val.Interface().(encoding.TextMarshaler); ok {
		return impl.MarshalText()
	}

	switch val.Kind() {
	case reflect.Bool:
		return []byte(strconv.FormatBool(val.Interface().(bool))), nil
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8:
		return MarshalInt(val.Int())
	case reflect.Float32, reflect.Float64:
		return MarshalFloat(val.Float())
	case reflect.String:
		return []byte(val.String()), nil
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
		return MarshalUint(val.Uint())
	}

	panic(fmt.Errorf("cannot marshal %s to text", val.Type().String()))
}

func UnmarshalText(dest any, text []byte) error {
	val := reflect.ValueOf(dest)

	if val.Kind() != reflect.Pointer {
		panic(fmt.Errorf("dest is not a pointer to the target"))
	}

	if impl, ok := val.Interface().(encoding.TextUnmarshaler); ok {
		return impl.UnmarshalText(text)
	}

	val = reflect.Indirect(val)

	switch val.Kind() {
	case reflect.Bool:
		var tmp Bool

		if err := tmp.UnmarshalText(text); err != nil {
			return err
		}

		val.Set(reflect.ValueOf(bool(tmp)))
		return nil
	case reflect.Int16:
		var tmp int16

		if err := UnmarshalInt(&tmp, text, math.MinInt16, math.MaxInt16, 0); err != nil {
			return err
		}

		val.SetInt(int64(tmp))
		return nil
	case reflect.Int32:
		var tmp int32

		if err := UnmarshalInt(&tmp, text, math.MinInt32, math.MaxInt32, 0); err != nil {
			return err
		}

		val.SetInt(int64(tmp))
		return nil
	case reflect.Int, reflect.Int64:
		var tmp int64

		if err := UnmarshalInt(&tmp, text, math.MinInt64, math.MaxInt64, 0); err != nil {
			return err
		}

		val.SetInt(tmp)
		return nil
	case reflect.Int8:
		var tmp int8

		if err := UnmarshalInt(&tmp, text, math.MinInt8, math.MaxInt8, 0); err != nil {
			return err
		}

		val.SetInt(int64(tmp))
		return nil
	case reflect.Float32:
		var tmp float32

		if err := UnmarshalFloat(&tmp, text, -math.MaxFloat32, math.MaxFloat32, 0); err != nil {
			return err
		}

		val.SetFloat(float64(tmp))
		return nil
	case reflect.Float64:
		var tmp float64

		if err := UnmarshalFloat(&tmp, text, -math.MaxFloat64, math.MaxFloat64, 0); err != nil {
			return err
		}

		val.SetFloat(tmp)
		return nil
	case reflect.String:
		val.SetString(string(text))
		return nil
	case reflect.Uint16:
		var tmp uint16

		if err := UnmarshalUint(&tmp, text, 0, math.MaxUint16, 0); err != nil {
			return err
		}

		val.SetUint(uint64(tmp))
		return nil
	case reflect.Uint32:
		var tmp uint32

		if err := UnmarshalUint(&tmp, text, 0, math.MaxUint32, 0); err != nil {
			return err
		}

		val.SetUint(uint64(tmp))
		return nil
	case reflect.Uint, reflect.Uint64:
		var tmp uint64

		if err := UnmarshalUint(&tmp, text, 0, math.MaxUint64, 0); err != nil {
			return err
		}

		val.SetUint(tmp)
		return nil
	case reflect.Uint8:
		var tmp uint8

		if err := UnmarshalUint(&tmp, text, 0, math.MaxUint8, 0); err != nil {
			return err
		}

		val.SetUint(uint64(tmp))
		return nil
	}

	panic(fmt.Errorf("cannot unmarshal %s from text", val.Type().String()))
}
