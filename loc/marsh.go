package loc

import (
	"encoding"
	"reflect"
)

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
		if val, err := loc.ParseBool(data); err != nil {
			return err
		} else {
			dst.SetBool(val)
			return nil
		}
	case reflect.Float32, reflect.Float64:
		if val, err := loc.ParseFloat(data); err != nil {
			return err
		} else if dst.OverflowFloat(val) {
			return OverflowError{}
		} else {
			dst.SetFloat(val)
			return nil
		}
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8:
		if val, err := loc.ParseInt(data); err != nil {
			return err
		} else if dst.OverflowInt(val) {
			return OverflowError{}
		} else {
			dst.SetInt(val)
			return nil
		}
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
		if val, err := loc.ParseUint(data); err != nil {
			return err
		} else if dst.OverflowUint(val) {
			return OverflowError{}
		} else {
			dst.SetUint(val)
			return nil
		}
	case reflect.String:
		dst.SetString(data)
		return nil
	}

	panic("cannot unmarshal text into a " + dst.Type().String())
}

type (
	TextUnmarshaler interface {
		UnmarshalLocText(loc Loc, data string) error
	}
)
