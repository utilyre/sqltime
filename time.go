package sqltime

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var ErrInvalidLayout = errors.New("Invalid time layout")

type Time struct {
	Hour   int
	Minute int
	Second int
}

func (t *Time) Parse(s string) error {
	parts := strings.Split(s, ":")
	if len(parts) != 3 {
		return ErrInvalidLayout
	}
	if len(parts[0]) != 2 || len(parts[1]) != 2 || len(parts[2]) != 2 {
		return ErrInvalidLayout
	}

	hh, err := strconv.Atoi(parts[0])
	if err != nil {
		return ErrInvalidLayout
	}

	mm, err := strconv.Atoi(parts[1])
	if err != nil {
		return ErrInvalidLayout
	}

	ss, err := strconv.Atoi(parts[2])
	if err != nil {
		return ErrInvalidLayout
	}

	t.Hour, t.Minute, t.Second = hh, mm, ss
	return nil
}

func (t Time) String() string {
	return fmt.Sprintf("%02d:%02d:%02d", t.Hour, t.Minute, t.Second)
}

func (s *Time) UnmarshalJSON(data []byte) error {
	str := ""
	json.Unmarshal(data, &str)

	return s.Parse(str)
}

func (s Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

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
		return fmt.Errorf("cannot sql.Scan() Time from: %#v", v)
	}
}

func (s Time) Value() (driver.Value, error) {
	return driver.Value(s.String()), nil
}
