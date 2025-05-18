package bind

import (
	"net/http"
	"reflect"
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
					if err := loc.UnmarshalText(l, cook.Value, ftgt.Interface()); err != nil {
						errs.Add("cookie/"+tag.Name, err)
					} else if set {
						fdst.Set(ftgt)
					}
				} else if !tag.Opt {
					errs.Add("cookie/"+tag.Name, check.RequiredError{})
				}
			case TagSrcHeader:
			case TagSrcPath:
				if val := r.PathValue(tag.Name); val != "" {
					if err := loc.UnmarshalText(l, val, ftgt.Interface()); err != nil {
						errs.Add("path/"+tag.Name, err)
					} else if set {
						fdst.Set(ftgt)
					}
				} else {
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

func Bind(r *http.Request, dest any) error {
	return BindCtx[any](nil, r, dest)
}
