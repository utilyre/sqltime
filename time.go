// Package sqltime provides support for Postgres time data type.
package sqltime

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Time is a one-to-one representation of SQL time data type.
type Time struct {
	hour   int
	minute int
	second int
}

// Hour returns the hour specified by t, in the range [0, 23].
func (t Time) Hour() int {
	return t.hour
}

// Minute returns the minute specified by t, in the range [0, 59].
func (t Time) Minute() int {
	return t.minute
}

// Second returns the second specified by t, in the range [0, 59].
func (t Time) Second() int {
	return t.second
}

// Add returns the time t+d.
func (t Time) Add(d Time) Time {
	return Time{
		hour:   (t.hour + d.hour) % 24,
		minute: (t.minute + d.minute) % 60,
		second: (t.second + d.second) % 60,
	}
}

// Sub returns the time t-d.
func (t Time) Sub(d Time) Time {
	return Time{
		hour:   (t.hour - d.hour) % 24,
		minute: (t.minute - d.minute) % 60,
		second: (t.second - d.second) % 60,
	}
}

// Parse parses a formatted string, in time.TimeOnly (15:04:05) format.
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

	t.hour, t.minute, t.second = hh, mm, ss
	return nil
}

var _ fmt.Stringer = Time{}

func (t Time) String() string {
	return fmt.Sprintf("%02d:%02d:%02d", t.hour, t.minute, t.second)
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
		s.hour, s.minute, s.second = v.Hour(), v.Minute(), v.Second()
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
