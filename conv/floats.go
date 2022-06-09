package conv

func MinMaxFloat(val, min, max float64) float64 {
	if val < min {
		return min
	}

	if val > max {
		return max
	}

	return val
}

func DefaultMaxFloat(val, def, max float64) float64 {
	if val <= 0 {
		return def
	}

	if val > max {
		return max
	}

	return val
}

func MaxFloat(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

// MinFloat returns the larger of a or b
func MinFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func RatioPercent(a, b float64) float64 {
	if b == 0 {
		return 0
	}
	return float64(int(((a / b) * 10000))) / 100.00
}

func FirstNonZeroFloat(list ...float64) float64 {
	for _, v := range list {
		if v != 0 {
			return v
		}
	}

	return 0
}

func FirstPositiveFloat(list ...float64) float64 {
	for _, v := range list {
		if v > 0 {
			return v
		}
	}

	return 0
}
