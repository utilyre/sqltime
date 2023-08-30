package sqltime

import "fmt"

// Part indicates a part of time (e.g. hour).
type Part int

const (
	PartHour Part = iota + 1
	PartMinute
	PartSecond
)

var _ fmt.Stringer = (Part)(0)

func (p Part) String() string {
	switch p {
	case PartHour:
		return "Hour"
	case PartMinute:
		return "Minute"
	case PartSecond:
		return "Second"
	}

	return ""
}
