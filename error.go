package sqltime

import (
	"errors"
	"fmt"
)

const (
	PartNone = iota
	PartHour
	PartMinute
	PartSecond
)

var (
	ErrTooManyParts = errors.New("too many parts")
	ErrNegativePart = errors.New("negative part")
)

type TimeError struct {
	Func string
	Part int
	Err  error
}

func tooManyPartsErr(fn string) *TimeError {
	return &TimeError{
		Func: fn,
		Part: PartNone,
		Err:  ErrTooManyParts,
	}
}

func atoiErr(fn string, part int, err error) *TimeError {
	return &TimeError{
		Func: fn,
		Part: part,
		Err:  err,
	}
}

func negativePartErr(fn string, part int) *TimeError {
	return &TimeError{
		Func: fn,
		Part: part,
		Err:  ErrNegativePart,
	}
}

func scanErr(fn string, value any) *TimeError {
	return &TimeError{
		Func: fn,
		Part: PartNone,
		Err:  fmt.Errorf("type sqltime.Time is incompatible with %v", value),
	}
}

func (e *TimeError) Error() string {
	part := ""
	switch e.Part {
	case PartHour:
		part = "parsing `Hour`: "
	case PartMinute:
		part = "parsing `Minute`: "
	case PartSecond:
		part = "parsing `Second`: "
	}

	return fmt.Sprintf("sqltime.%s: %s%s", e.Func, part, e.Err)
}

func (e *TimeError) Unwrap() error {
	return e.Err
}
