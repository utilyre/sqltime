package sqltime

// Parse parses a formatted string, in time.TimeOnly (15:04:05) format, and
// returns the time value it represents.
func Parse(s string) (Time, error) {
	t := Time{}
	if err := t.Parse(s); err != nil {
		return Time{}, err
	}

	return t, nil
}
