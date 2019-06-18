package conv

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFToC(t *testing.T) {
	i := -33.5
	assert.Equal(t, -36.39, FToC(i))
	i = -12.0
	assert.Equal(t, -24.44, FToC(i))
	i = 0.0
	assert.Equal(t, -17.78, FToC(i))
	i = 32.0
	assert.Equal(t, 0.0, FToC(i))
	i = 74.0
	assert.Equal(t, 23.33, FToC(i))
	i = 212.0
	assert.Equal(t, 100.0, FToC(i))
}

func TestMilesToKm(t *testing.T) {
	i := 0.0
	assert.Equal(t, 0.0, MilesToKm(i))
	i = 10.5
	assert.Equal(t, 16.9, MilesToKm(i))
	i = 42
	assert.Equal(t, 67.59, MilesToKm(i))
	i = 62.14
	assert.Equal(t, 100.0, MilesToKm(i))
}

func TestFeetToMeters(t *testing.T) {
	i := 0.0
	assert.Equal(t, 0.0, FeetToMeters(i))
	i = 3.0
	assert.Equal(t, 0.91, FeetToMeters(i))
	i = 6.0
	assert.Equal(t, 1.83, FeetToMeters(i))
	i = 17.2
	assert.Equal(t, 5.24, FeetToMeters(i))
}

func TestRound(t *testing.T) {
	i := 10.5
	assert.Equal(t, 11, Round(i))
	i = 0.0
	assert.Equal(t, 0, Round(i))
	i = -1.1
	assert.Equal(t, -1, Round(i))
	i = -1.9
	assert.Equal(t, -2, Round(i))
}

func TestToFixed(t *testing.T) {
	i := 3.1415926
	assert.Equal(t, 3.14, ToFixed(i, 2))
	assert.Equal(t, 3.1415926, ToFixed(i, 8))
	assert.Equal(t, 3.0, ToFixed(i, 0))

	a := 0.1
	b := 0.2
	assert.NotEqual(t, 0.3, a+b)
	assert.Equal(t, 0.3, ToFixed(a+b, 1))
}

func TestFormatDuration(t *testing.T) {
	testCases := []struct {
		Label    string
		Input    time.Time
		Expected string
	}{
		{"1 Day Ago", time.Now().Round(0).Add(-(1 * 24) * time.Hour), "1d"},
		{"1 Day from Now", time.Now().Round(0).Add(24 * time.Hour), "now"},
		{"1 Hour Ago", time.Now().Round(0).Add(-(1) * time.Hour), "1h"},
		{"1 Hour from Now", time.Now().Round(0).Add(time.Hour), "now"},
		{"1 Minute Ago", time.Now().Round(0).Add(-(1) * time.Minute), "1m"},
		{"1 Minute from Now", time.Now().Round(0).Add(time.Minute), "now"},
		{"1 Month Ago", time.Now().Round(0).Add(-(1 * 30 * 24) * time.Hour), "30d"},
		{"1 Second Ago", time.Now().Round(0).Add(-(1) * time.Second), "now"},
		{"1 Second from Now", time.Now().Round(0).Add(time.Second), "now"},
		{"1 Year Ago", time.Now().Round(0).Add(-(1 * 24 * 365) * time.Hour), "1y"},
		{"1 Year from Now", time.Now().Round(0).Add(24 * 365 * time.Hour), "now"},
		{"1/2 Second Ago", time.Now().Round(0).Add(-(500) * time.Millisecond), "now"},
		{"100 Months Ago", time.Now().Round(0).Add(-(100 * 24 * 30) * time.Hour), "8y"},
		{"100 Years Ago", time.Now().Round(0).Add(-(100 * 24 * 365) * time.Hour), "100y"},
		{"11 Months Ago", time.Now().Round(0).Add(-(11 * 30 * 24) * time.Hour), "11mo"},
		{"12 Minutes Ago", time.Now().Round(0).Add(-(12) * time.Minute), "12m"},
		{"17 Day Ago", time.Now().Round(0).Add(-(42 * 24) * time.Hour), "42d"},
		{"17 Hours Ago", time.Now().Round(0).Add(-(17) * time.Hour), "17h"},
		{"1825 Days Ago", time.Now().Round(0).Add(-(1825 * 24) * time.Hour), "5y"},
		{"2 Months Ago", time.Now().Round(0).Add(-(2 * 30 * 24) * time.Hour), "2mo"},
		{"23 Hours Ago", time.Now().Round(0).Add(-(23) * time.Hour), "23h"},
		{"32 Seconds Ago", time.Now().Round(0).Add(-(32) * time.Second), "now"},
		{"63 Minutes Ago", time.Now().Round(0).Add(-(63) * time.Minute), "1h"},
		{"69 Seconds Ago", time.Now().Round(0).Add(-(69) * time.Second), "1m"},
		{"Not Quite 1 Year Ago", time.Now().Round(0).Add(-(1 * 24 * 356) * time.Hour), "11mo"},
	}

	for _, tc := range testCases {
		t.Run(tc.Label, func(t *testing.T) {
			assert.Equal(t, tc.Expected, RelativeTime(time.Since(tc.Input)))
		})
	}
}
