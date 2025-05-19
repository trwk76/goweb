package bind

import (
	"net/http"
	"reflect"

	"github.com/trwk76/goweb/check"
	"github.com/trwk76/goweb/content/form"
	"github.com/trwk76/goweb/loc"
)

// WithCtx binds the request to the destination structure.
// It uses the "bind" tag to determine the source of the value and potential options. Please refer to the Tag struct for more details.
// dest must be a valid pointer to a structure otherwise the function will panic.
// The function will return an error if the binding or checks fail.
func WithCtx[C any](ctx C, l loc.Loc, r *http.Request, dest any) error {
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
					if err := form.UnmarshalURLValue(ctx, l, cook.Value, ftgt); err != nil {
						errs.Add("cookie/"+tag.Name, err)
					} else if set {
						fdst.Set(ftgt)
					}
				} else if !tag.Opt {
					errs.Add("cookie/"+tag.Name, check.RequiredError{})
				}
			case TagSrcHeader:
				if vals := r.Header.Values(tag.Name); len(vals) > 0 {
					if err := form.UnmarshalURLValues(ctx, l, tag.Tag, vals, ftgt); err != nil {
						errs.Add("header/"+tag.Name, err)
					} else if set {
						fdst.Set(ftgt)
					}
				} else if !tag.Opt {
					errs.Add("header/"+tag.Name, check.RequiredError{})
				}
			case TagSrcPath:
				if val := r.PathValue(tag.Name); val != "" {
					if err := form.UnmarshalURLValue(ctx, l, val, ftgt); err != nil {
						errs.Add("path/"+tag.Name, err)
					} else if set {
						fdst.Set(ftgt)
					}
				} else if !tag.Opt {
					errs.Add("path/"+tag.Name, check.RequiredError{})
				}
			case TagSrcQuery:
				if vals := r.URL.Query()[tag.Name]; len(vals) > 0 {
					if err := form.UnmarshalURLValues(ctx, l, tag.Tag, vals, ftgt); err != nil {
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
