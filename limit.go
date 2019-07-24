package bgo

func stringLimit(s string, l int) string {
	if len(s) > l {
		return s[:l]
	}
	return s
}
