package sqltime

func Parse(s string) (*Time, error) {
	t := &Time{}
	if err := t.Parse(s); err != nil {
		return nil, err
	}

	return t, nil
}
