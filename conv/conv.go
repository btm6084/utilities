// Package conv provides handy conversion functions.
package conv

import (
	"fmt"
	"math"
	"time"
)

// FToC converts fahrenheit to celsius
func FToC(f float64) float64 {
	return ToFixed((float64(f)-32.0)*(5.0/9.0), 2)
}

// MilesToKm converts miles to kilometers
func MilesToKm(m float64) float64 {
	return ToFixed(1.60934*m, 2)
}

// FeetToMeters converts feet to meters
func FeetToMeters(f float64) float64 {
	return ToFixed(0.3048*f, 2)
}

// Round a floating point number to the nearest integer.
// -0.5 becomes -1, 0.5 becomes 1.
func Round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

// ToFixed will set a fixed precision on a float.
func ToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(Round(num*output)) / output
}

// RelativeTime takes a time.Duration and returns the largest denomination
// relative time. eg. 1y for a duration over 400 days.
func RelativeTime(dur time.Duration) string {
	y := dur / (365 * 24 * time.Hour)
	if y > 0 {
		return fmt.Sprintf("%dy", y)
	}

	mo := dur / (30 * 24 * time.Hour)
	if mo > 1 {
		return fmt.Sprintf("%dmo", mo)
	}

	d := dur / (24 * time.Hour)
	if d > 0 {
		return fmt.Sprintf("%dd", d)
	}

	h := dur / time.Hour
	if h > 0 {
		return fmt.Sprintf("%dh", h)
	}

	m := dur / time.Minute
	if m > 0 {
		return fmt.Sprintf("%dm", m)
	}

	return "now"
}
