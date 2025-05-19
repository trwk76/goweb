package form

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/trwk76/goweb/check"
	"github.com/trwk76/goweb/loc"
)

func UnmarshalURLValues[C any](ctx C, l loc.Loc, tag Tag, vals []string, dst reflect.Value) error {
	switch tag.Del {
	case TagDelComma, TagDelSpace, TagDelPipe:
		vals = splitDeclVals(vals, string(tag.Del))
	}

	switch dst.Kind() {
	case reflect.Array:
		if len(vals) < dst.Len() {
			return check.ValueCountError{Expected: dst.Len(), Actual: len(vals)}
		}

		errs := check.NewFieldErrors()
		dst = dst.Elem()

		for idx, val := range vals {
			tgt := dst.Index(idx).Addr()
			set := false

			if tgt.Kind() == reflect.Pointer {
				tgt = reflect.New(tgt.Type().Elem())
				set = true
			}

			// unmarshalValue will check the value and return an error if it is not valid
			if err := UnmarshalURLValue(ctx, l, val, tgt); err != nil {
				errs.Add(strconv.Itoa(idx), err)
			} else if set {
				dst.Index(idx).Set(tgt)
			}
		}

		return errs.Err()
	case reflect.Slice:
		errs := check.NewFieldErrors()
		dst = dst.Elem()

		if dst.Len() < len(vals) {
			// grow slice to accommodate all values
			if dst.Cap() < len(vals) {
				dst.Grow(len(vals))
			}

			dst.SetLen(len(vals))
		}

		for idx, val := range vals {
			tgt := dst.Index(idx).Addr()
			set := false

			if tgt.Kind() == reflect.Pointer {
				tgt = reflect.New(tgt.Type().Elem())
				set = true
			}

			// unmarshalValue will check the value and return an error if it is not valid
			if err := UnmarshalURLValue(ctx, l, val, tgt); err != nil {
				errs.Add(strconv.Itoa(idx), err)
			} else if set {
				dst.Index(idx).Set(tgt)
			}
		}

		return errs.Err()
	}

	// dst is neither a slice nor an array, so use first value only
	return UnmarshalURLValue(ctx, l, vals[0], dst)
}

func UnmarshalURLValue[C any](ctx C, l loc.Loc, val string, dst reflect.Value) error {
	if um, ok := dst.Interface().(URLUnmarshalerCtx[C]); ok {
		return um.UnmarshalURLCtx(ctx, l, val)
	} else if um, ok := dst.Interface().(URLUnmarshaler); ok {
		return um.UnmarshalURL(l, val)
	}

	if err := loc.UnmarshalText(l, val, dst.Addr().Interface()); err != nil {
		return err
	}

	return check.WithCtx(ctx, dst.Elem().Interface())
}

func splitDeclVals(vals []string, delim string) []string {
	res := make([]string, 0, len(vals))

	for _, val := range vals {
		res = append(res, strings.Split(val, delim)...)
	}

	return res
}
