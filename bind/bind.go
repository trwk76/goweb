package bind

import (
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/trwk76/goweb/check"
	"github.com/trwk76/goweb/loc"
)

func BindCtx[C any](ctx C, l loc.Loc, r *http.Request, dest any) error {
	if b, ok := dest.(BindableCtx[C]); ok {
		return b.BindCtx(ctx, r)
	} else if b, ok := dest.(Bindable); ok {
		return b.Bind(r)
	}

	dst := reflect.ValueOf(dest)
	if dst.Kind() != reflect.Pointer || dst.IsNil() {
		panic("dest is not a valid pointer")
	}

	dst = dst.Elem()
	if dst.Kind() != reflect.Struct {
		panic("dest is not a valid pointer to a structure")
	}

	errs := check.NewFieldErrors()
	typ := dst.Type()

	for _, fld := range reflect.VisibleFields(typ) {
		if rtag := fld.Tag.Get("bind"); rtag != "" {
			tag := MustParseTag(typ, fld, rtag)
			fdst := dst.FieldByIndex(fld.Index)
			ftgt := fdst.Addr()
			set := false

			if fdst.Kind() == reflect.Pointer {
				ftgt = reflect.New(fdst.Type().Elem())
				set = true
			}

			switch tag.Src {
			case TagSrcBody:

			case TagSrcCookie:
				if cook, err := r.Cookie(tag.Name); cook != nil && err == nil {
					if err := unmarshalValue(ctx, l, cook.Value, ftgt); err != nil {
						errs.Add("cookie/"+tag.Name, err)
					} else if set {
						fdst.Set(ftgt)
					}
				} else if !tag.Opt {
					errs.Add("cookie/"+tag.Name, check.RequiredError{})
				}
			case TagSrcHeader:
				if vals := r.Header.Values(tag.Name); len(vals) > 0 {
					if err := unmarshalValues(ctx, l, tag, vals, ftgt); err != nil {
						errs.Add("header/"+tag.Name, err)
					} else if set {
						fdst.Set(ftgt)
					}
				} else if !tag.Opt {
					errs.Add("header/"+tag.Name, check.RequiredError{})
				}
			case TagSrcPath:
				if val := r.PathValue(tag.Name); val != "" {
					if err := unmarshalValue(ctx, l, val, ftgt); err != nil {
						errs.Add("path/"+tag.Name, err)
					} else if set {
						fdst.Set(ftgt)
					}
				} else if !tag.Opt {
					errs.Add("path/"+tag.Name, check.RequiredError{})
				}
			case TagSrcQuery:
				if vals := r.URL.Query()[tag.Name]; len(vals) > 0 {
					if err := unmarshalValues(ctx, l, tag, vals, ftgt); err != nil {
						errs.Add("query/"+tag.Name, err)
					} else if set {
						fdst.Set(ftgt)
					}
				} else if !tag.Opt {
					errs.Add("query/"+tag.Name, check.RequiredError{})
				}
			}
		}
	}

	return errs.Err()
}

func Bind(l loc.Loc, r *http.Request, dest any) error {
	return BindCtx[any](nil, l, r, dest)
}

func unmarshalValues[C any](ctx C, l loc.Loc, tag Tag, vals []string, dst reflect.Value) error {
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
			if err := unmarshalValue(ctx, l, val, tgt); err != nil {
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
			if err := unmarshalValue(ctx, l, val, tgt); err != nil {
				errs.Add(strconv.Itoa(idx), err)
			} else if set {
				dst.Index(idx).Set(tgt)
			}
		}

		return errs.Err()
	}

	// dst is neither a slice nor an array, so use first value only
	return unmarshalValue(ctx, l, vals[0], dst)
}

func unmarshalValue[C any](ctx C, l loc.Loc, val string, dst reflect.Value) error {
	if um, ok := dst.Interface().(URLUnmarshalerCtx[C]); ok {
		return um.UnmarshalURLCtx(ctx, l, val)
	} else if um, ok := dst.Interface().(URLUnmarshaler); ok {
		return um.UnmarshalURL(l, val)
	}

	if err := loc.UnmarshalText(l, val, dst.Addr().Interface()); err != nil {
		return err
	}

	return check.CheckCtx(ctx, dst.Elem().Interface())
}

func splitDeclVals(vals []string, delim string) []string {
	res := make([]string, 0, len(vals))

	for _, val := range vals {
		res = append(res, strings.Split(val, delim)...)
	}

	return res
}
