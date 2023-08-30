// Package sqltime provides support for Postgres time data type.
//
// # Examples
//
// ## GORM
//
// **NOTE**: Do not use gorm.DB.AutoMigrate(&Example{}) for any model containing
// sqltime.Time, because GORM will try to create a column of type timestamptz
// instead.
//
//	type Example struct {
//		gorm.Model
//		// ...
//		Clock sqltime.Time `gorm:"type:time"`
//	}
package sqltime

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Time is a one-to-one representation of Postgres time data type.
type Time struct {
	Hour   int
	Minute int
	Second int
}

type Part int

var _ fmt.Stringer = (Part)(0)

func (p Part) String() string {
	switch p {
	case PartHour:
		return "Hour"
	case PartMinute:
		return "Minute"
	case PartSecond:
		return "Second"
	default:
		return ""
	}
}

const (
	PartHour Part = iota + 1
	PartMinute
	PartSecond
)

type RangeError struct {
	Part Part
	Min  int
	Max  int
}

var _ error = (*RangeError)(nil)

func (e *RangeError) Error() string {
	return fmt.Sprintf("%s out of range [%d, %d]", e.Part, e.Min, e.Max)
}

type AtoiError struct {
	Part Part
	Err  error
}

var _ error = (*AtoiError)(nil)

func (e *AtoiError) Error() string {
	return fmt.Sprintf("%s %v", e.Part, e.Err)
}

var _ interface{ Unwrap() error } = (*AtoiError)(nil)

func (e *AtoiError) Unwrap() error {
	return e.Err
}

var (
	ErrTooManyParts = errors.New("too many parts")
)

func (t *Time) Parse(s string) error {
	parts := strings.Split(s, ":")
	if len(parts) > 3 {
		return ErrTooManyParts
	}

	hh, err := strconv.Atoi(parts[0])
	if err != nil {
		return &AtoiError{
			Part: PartHour,
			Err:  err,
		}
	}
	if hh < 0 || hh > 23 {
		return &RangeError{
			Part: PartHour,
			Min:  0,
			Max:  23,
		}
	}

	mm := 0
	if len(parts) > 1 {
		mm, err = strconv.Atoi(parts[1])
		if err != nil {
			return &AtoiError{
				Part: PartMinute,
				Err:  err,
			}
		}
		if mm < 0 || mm > 59 {
			return &RangeError{
				Part: PartMinute,
				Min:  0,
				Max:  59,
			}
		}
	}

	ss := 0
	if len(parts) > 2 {
		ss, err = strconv.Atoi(parts[2])
		if err != nil {
			return &AtoiError{
				Part: PartSecond,
				Err:  err,
			}
		}
		if ss < 0 || ss > 59 {
			return &RangeError{
				Part: PartSecond,
				Min:  0,
				Max:  59,
			}
		}
	}

	t.Hour, t.Minute, t.Second = hh, mm, ss
	return nil
}

var _ fmt.Stringer = Time{}

func (t Time) String() string {
	return fmt.Sprintf("%02d:%02d:%02d", t.Hour, t.Minute, t.Second)
}

var _ json.Unmarshaler = (*Time)(nil)

func (s *Time) UnmarshalJSON(data []byte) error {
	str := ""
	json.Unmarshal(data, &str)

	return s.Parse(str)
}

var _ json.Marshaler = (*Time)(nil)

func (s Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

var _ sql.Scanner = (*Time)(nil)

func (s *Time) Scan(src any) error {
	switch v := src.(type) {
	case []byte:
		return s.Parse(string(v))
	case string:
		return s.Parse(v)
	case time.Time:
		s.Hour, s.Minute, s.Second = v.Hour(), v.Minute(), v.Second()
		return nil
	case nil:
		*s = Time{}
		return nil
	default:
		return fmt.Errorf("type sqltime.Time is incompatible with value %v", v)
	}
}

var _ driver.Valuer = Time{}

func (s Time) Value() (driver.Value, error) {
	return driver.Value(s.String()), nil
}
