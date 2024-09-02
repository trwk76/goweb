package marshal

import (
	"encoding"
	"fmt"
	"time"

	"github.com/trwk76/goweb/openapi/spec"
)

func ParseDate(str string) (Date, error) {
	res, err := time.Parse(dateLayout, str)
	if err != nil {
		return Date{}, fmt.Errorf("not a date value (ex: 2024-12-31)")
	}

	return Date(res), nil
}

func (v Date) String() string {
	return time.Time(v).Format(dateLayout)
}

func ParseTime(str string) (Time, error) {
	res, err := time.Parse(timeLayout, str)
	if err != nil {
		return Time{}, fmt.Errorf("not a date value (ex: 15:59:59.9999999Z)")
	}

	return Time(res), nil
}

func (v Time) String() string {
	return time.Time(v).Format(timeLayout)
}

func ParseDateTime(str string) (DateTime, error) {
	res, err := time.Parse(dateTimeLayout, str)
	if err != nil {
		return DateTime{}, fmt.Errorf("not a date/time value (ex: 2024-12-31T15:59:59.9999999+0000)")
	}

	return DateTime(res), nil
}

func (v DateTime) String() string {
	return time.Time(v).Format(dateTimeLayout)
}

type (
	Date     time.Time
	Time     time.Time
	DateTime time.Time
)

func (Date) Schema() spec.Schema {
	return spec.Schema{Type: spec.TypeString, Format: spec.Format("date"), Pattern: "^\\d{4}-\\d{2}-\\d{2}$"}
}

func (Time) Schema() spec.Schema {
	return spec.Schema{Type: spec.TypeString, Format: spec.Format("time"), Pattern: "^\\d{2}%3A\\d{2}%3A\\d{2}(?:%2E\\d+)?[A-Z]?(?:[+.-](?:08%3A\\d{2}|\\d{2}[A-Z]))?$"}
}

func (DateTime) Schema() spec.Schema {
	return spec.Schema{Type: spec.TypeString, Format: spec.Format("date-time"), Pattern: "^\\d{4}-\\d{2}-\\d{2}T\\d{2}%3A\\d{2}%3A\\d{2}(?:%2E\\d+)?[A-Z]?(?:[+.-](?:08%3A\\d{2}|\\d{2}[A-Z]))?$"}
}

func (v Date) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

func (v *Date) UnmarshalText(raw []byte) error {
	tmp, err := ParseDate(string(raw))
	if err != nil {
		return err
	}

	*v = tmp
	return nil
}

func (v Time) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

func (v *Time) UnmarshalText(raw []byte) error {
	tmp, err := ParseTime(string(raw))
	if err != nil {
		return err
	}

	*v = tmp
	return nil
}

func (v DateTime) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

func (v *DateTime) UnmarshalText(raw []byte) error {
	tmp, err := ParseDateTime(string(raw))
	if err != nil {
		return err
	}

	*v = tmp
	return nil
}

const (
	dateLayout     string = time.DateOnly
	timeLayout     string = "15:04:05.999999999Z07:00"
	dateTimeLayout string = time.RFC3339Nano
)

var (
	_ encoding.TextMarshaler = Date{}
	_ encoding.TextMarshaler = Time{}
	_ encoding.TextMarshaler = DateTime{}

	_ encoding.TextUnmarshaler = (*Date)(nil)
	_ encoding.TextUnmarshaler = (*Time)(nil)
	_ encoding.TextUnmarshaler = (*DateTime)(nil)
)
