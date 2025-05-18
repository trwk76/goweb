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
	if dst.Kind() != reflect.Struct {
		panic("dest is not a valid pointer to a structure")
	}

	if err := r.ParseForm(); err != nil {
		return err
	}

	// TODO

	return errs.Err()
}

func Parse(l loc.Loc, r *http.Request, dest any) error {
	return ParseCtx[any](nil, l, r, dest)
}
