package sqltime

// See https://github.com/golang/go/issues/448
func div(a, b int) int {
	c := a / b
	if a%b < 0 {
		c -= 1
	}

	return c
}

// See https://github.com/golang/go/issues/448
func mod(a, b int) int {
	return ((a % b) + b) % b
}
