package marshal

import (
	"fmt"
	"strings"
)

type (
	PathItem interface {
		fmt.Stringer
		pathItem()
	}

	PathMember string
	PathIndex  int

	Path []PathItem

	PathError struct {
		Path Path
		Err  error
	}

	PathErrors []PathError
)

func NewPathErrors(item PathItem, err error) PathErrors {
	if err == nil {
		return nil
	}

	switch terr := err.(type) {
	case PathError:
		return PathErrors{PathError{Path: append(Path{item}, terr.Path...), Err: terr.Err}}
	case PathErrors:
		res := make(PathErrors, len(terr))

		for idx, itm := range terr {
			res[idx] = PathError{Path: append(Path{item}, itm.Path...), Err: itm.Err}
		}

		return res
	}

	return PathErrors{PathError{Path: Path{item}, Err: err}}
}

func (m PathMember) String() string {
	return string(m)
}

func (m PathIndex) String() string {
	return fmt.Sprintf("[%d]", int(m))
}

func (m Path) String() string {
	res := strings.Builder{}

	for idx, itm := range m {
		switch titm := itm.(type) {
		case PathMember:
			if idx > 0 {
				res.WriteByte('.')
			}

			res.WriteString(titm.String())
		case PathIndex:
			res.WriteString(titm.String())
		}
	}

	return res.String()
}

func (e PathError) Error() string {
	return fmt.Sprintf("%s: %s", e.Path.String(), e.Err.Error())
}

func (e PathErrors) Error() string {
	return fmt.Sprintf("%d error(s) occurred", len(e))
}

func (PathMember) pathItem() {}
func (PathIndex) pathItem()  {}

var (
	_ PathItem = PathMember("")
	_ PathItem = PathIndex(0)
	_ error    = PathError{}
	_ error    = PathErrors{}
)
