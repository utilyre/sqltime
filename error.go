package sqltime

import (
	"errors"
	"fmt"
)

var ErrTooManyParts = errors.New("too many parts")

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
