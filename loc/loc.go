package loc

import "time"

type (
	Loc interface {
		FormatBool(v bool) string
		FormatInt(v int64) string
		FormatUint(v uint64) string
		FormatFloat(v float64) string
		FormatDate(v time.Time) string
		FormatTime(v time.Time) string
		FormatDateTime(v time.Time) string

		ParseBool(v string) (bool, error)
		ParseInt(v string) (int64, error)
		ParseUint(v string) (uint64, error)
		ParseFloat(v string) (float64, error)
		ParseDate(v string) (time.Time, error)
		ParseTime(v string) (time.Time, error)
		ParseDateTime(v string) (time.Time, error)
	}
)
