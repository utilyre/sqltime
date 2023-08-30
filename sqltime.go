package sqltime

func Parse(s string) (Time, error) {
	t := Time{}
	if err := t.Parse(s); err != nil {
		return Time{}, err
	}

	return t, nil
}
