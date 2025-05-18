package loc

import (
	"encoding"
	"reflect"
)

// UnmarshalText unmarshals a string into the destination value using the provided localization.
func UnmarshalText(loc Loc, data string, dest any) error {
	if um, ok := dest.(TextUnmarshaler); ok {
		return um.UnmarshalLocText(loc, data)
	} else if um, ok := dest.(encoding.TextUnmarshaler); ok {
		return um.UnmarshalText([]byte(data))
	}

	dst := reflect.ValueOf(dest)
	if dst.Kind() != reflect.Pointer || dst.IsNil() {
		panic("dest is not a valid pointer")
	}

	dst = dst.Elem()

	switch dst.Kind() {
	case reflect.Bool:
		val, err := loc.ParseBool(data)
		if err != nil {
			return err
		}

		dst.SetBool(val)
		return nil
	case reflect.Float32, reflect.Float64:
		val, err := loc.ParseFloat(data)
		if err != nil {
			return err
		} else if dst.OverflowFloat(val) {
			return OverflowError{}
		}

		dst.SetFloat(val)
		return nil
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8:
		val, err := loc.ParseInt(data)
		if err != nil {
			return err
		} else if dst.OverflowInt(val) {
			return OverflowError{}
		}

		dst.SetInt(val)
		return nil
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
		val, err := loc.ParseUint(data)
		if err != nil {
			return err
		} else if dst.OverflowUint(val) {
			return OverflowError{}
		}

		dst.SetUint(val)
		return nil
	case reflect.String:
		dst.SetString(data)
		return nil
	}

	panic("cannot unmarshal text into a " + dst.Type().String())
}

type (
	// TextMarshaler is an interface that defines a method for marshaling a value into a string using localization.
	TextMarshaler interface {
		MarshalLocText(loc Loc) (string, error)
	}

	// TextUnmarshaler is an interface that defines a method for unmarshaling a string into a value using localization.
	TextUnmarshaler interface {
		UnmarshalLocText(loc Loc, data string) error
	}
)
