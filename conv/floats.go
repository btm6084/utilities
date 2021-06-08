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

func MaxFloat(a, b float64) float64 {
	if a > b {
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
