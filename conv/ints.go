package conv

import (
	"math/rand"
)

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

// FuzzInt adds or subtracts a certain percentage from the input.
// Ex: FuzzInt(100, 1, 100) will have an output range between 0 and 100
// Ex: FuzzInt(100, 1, 10) will have an output range between 90 and 110, and will add/remove at least 1%
// Ex: FuzzInt(100, 8, 10) will have an output range between 90 and 110, and will add/remove at least 8%
func FuzzInt(input int, min, max int) int {
	if input < 2 {
		return 0
	}

	max = MaxInt(min, max)

	pct := float64(rand.Intn((max-min)+1)+min) / 100

	if rand.Intn(2) == 1 {
		pct += 1.0
	} else {
		pct = 1.0 - pct
	}

	return Round(float64(input) * pct)
}
