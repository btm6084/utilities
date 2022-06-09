package conv

// MinMaxInt returns the min if val is < min, max if val is > max.
func MinMaxInt(val, min, max int) int {
	if val < min {
		return min
	}

	if val > max {
		return max
	}

	return val
}

// DefaultMaxInt returns the default if val is < 0, max if val is > max.
func DefaultMaxInt(val, def, max int) int {
	if val <= 0 {
		return def
	}

	if val > max {
		return max
	}

	return val
}

// MaxInt returns the larger of a or b
func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// MaxInt returns the larger of a or b
func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func FirstNonZeroInt(list ...int) int {
	for _, v := range list {
		if v != 0 {
			return v
		}
	}

	return 0
}

func FirstPositiveInt(list ...int) int {
	for _, v := range list {
		if v > 0 {
			return v
		}
	}

	return 0
}
