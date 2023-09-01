package sqltime

// See https://github.com/golang/go/issues/448
func div(a, b int) int {
	return (a - mod(a, b)) / b
}

// See https://github.com/golang/go/issues/448
func mod(a, b int) int {
	return ((a % b) + b) % b
}
