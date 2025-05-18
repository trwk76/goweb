package loc

import "strconv"

// Default is the default implementation of the Loc interface using the standard functions.
var Default Loc = defLoc{}

type (
	// Loc is an interface that defines methods for formatting and parsing various types in a specific localization.
	Loc interface {
		FormatBool(v bool) string
		FormatInt(v int64) string
		FormatUint(v uint64) string
		FormatFloat(v float64) string

		ParseBool(v string) (bool, error)
		ParseInt(v string) (int64, error)
		ParseUint(v string) (uint64, error)
		ParseFloat(v string) (float64, error)
	}

	defLoc struct{}
)

func (defLoc) FormatBool(v bool) string {
	return strconv.FormatBool(v)
}

func (defLoc) FormatInt(v int64) string {
	return strconv.FormatInt(v, 10)
}

func (defLoc) FormatUint(v uint64) string {
	return strconv.FormatUint(v, 10)
}

func (defLoc) FormatFloat(v float64) string {
	return strconv.FormatFloat(v, 'g', -1, 64)
}

func (defLoc) ParseBool(v string) (bool, error) {
	return strconv.ParseBool(v)
}

func (defLoc) ParseInt(v string) (int64, error) {
	return strconv.ParseInt(v, 10, 64)
}

func (defLoc) ParseUint(v string) (uint64, error) {
	return strconv.ParseUint(v, 10, 64)
}

func (defLoc) ParseFloat(v string) (float64, error) {
	return strconv.ParseFloat(v, 64)
}
