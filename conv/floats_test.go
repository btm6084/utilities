package conv

import (
	"testing"

	"github.com/spf13/cast"
	"github.com/stretchr/testify/require"
)

func TestMinMaxFloats(t *testing.T) {
	testCases := []struct {
		Min      float64
		Max      float64
		Val      float64
		Expected float64
	}{
		{0, 100, 50, 50},
		{0, 75, 50, 50},
		{25, 75, 50, 50},
		{55, 75, 50, 55},
		{55, 75, 100, 75},
	}

	for k, tc := range testCases {
		t.Run(cast.ToString(k+1), func(t *testing.T) {
			require.Equal(t, tc.Expected, MinMaxFloat(tc.Val, tc.Min, tc.Max))
		})
	}
}

func TestMaxFloats(t *testing.T) {
	testCases := []struct {
		A        float64
		B        float64
		Expected float64
	}{
		{0, 100, 100},
		{101, 100, 101},
		{100, 0, 100},
		{-100, 100, 100},
		{-1, -2, -1},
	}

	for k, tc := range testCases {
		t.Run(cast.ToString(k+1), func(t *testing.T) {
			require.Equal(t, tc.Expected, MaxFloat(tc.A, tc.B))
		})
	}
}

func TestRoundPercent(t *testing.T) {
	testCases := []struct {
		A        float64
		B        float64
		Expected float64
	}{
		{0, 100, 0},
		{0, 0, 0},
		{17.53, 100, 17.53},
		{17.53, 87.231927, 20.09},
		{22.23, 17, 130.76},
		{43.88876643562436523456, 111, 39.53},
		{-1, -2, 50},
		{-1.234234234, -2.3421, 52.69},
		{-1.234234234, -1.234234234, 100.00},
	}

	for k, tc := range testCases {
		t.Run(cast.ToString(k+1), func(t *testing.T) {
			require.Equal(t, tc.Expected, RatioPercent(tc.A, tc.B))
		})
	}
}
