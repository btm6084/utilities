package conv

func MinMaxInt(val, min, max int) int {
	if val < min {
		return min
	}

	if val > max {
		return max
	}

	return val
}

func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
