package form

import (
	"net/http"
	"reflect"

	"github.com/trwk76/goweb/check"
	"github.com/trwk76/goweb/loc"
)

func ParseCtx[C any](ctx C, l loc.Loc, r *http.Request, dest any) error {
	dst := reflect.ValueOf(dest)
	errs := check.NewFieldErrors()

	if dst.Kind() != reflect.Pointer || dst.IsNil() {
		panic("dest is not a valid pointer")
	}

	dst = dst.Elem()
	typ := dst.Type()
	if typ.Kind() != reflect.Struct {
		panic("dest is not a valid pointer to a structure")
	}

	if err := r.ParseForm(); err != nil {
		return err
	}

	for _, fld := range reflect.VisibleFields(typ) {
		if tag := fld.Tag.Get("form"); tag != "-" {
			t := MustParseTag(typ, fld, tag)

			if t.Name == "" {
				t.Name = fld.Name
			}

			fdst := dst.FieldByIndex(fld.Index)
			ftgt := fdst.Addr()
			set := false

			if fld.Type.Kind() == reflect.Pointer {
				ftgt = reflect.New(fld.Type.Elem())
				set = true
			}

			if vals := r.PostForm[t.Name]; len(vals) > 0 {
				if err := UnmarshalURLValues(ctx, l, t, vals, ftgt); err != nil {
					errs.Add(t.Name, err)
				} else if set {
					fdst.Set(ftgt)
				}
			} else if !t.Opt {
				errs.Add(t.Name, check.RequiredError{})
			}
		}
	}

	return errs.Err()
}

func Parse(l loc.Loc, r *http.Request, dest any) error {
	return ParseCtx[any](nil, l, r, dest)
}
