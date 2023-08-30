package sqltime

func Parse(s string) Time {
	t := Time{}
	t.Parse(s)
	return t
}
