package conv

import (
	"testing"

	"github.com/spf13/cast"
	"github.com/stretchr/testify/require"
)

func TestMinMaxInts(t *testing.T) {
	testCases := []struct {
		Min      int
		Max      int
		Val      int
		Expected int
	}{
		{0, 100, 50, 50},
		{0, 75, 50, 50},
		{25, 75, 50, 50},
		{55, 75, 50, 55},
		{55, 75, 100, 75},
		{-1, -2, 100, -2},
		{-1, -2, -100, -1},
		{-2, -1, -100, -2},
	}

	for k, tc := range testCases {
		t.Run(cast.ToString(k+1), func(t *testing.T) {
			require.Equal(t, tc.Expected, MinMaxInt(tc.Val, tc.Min, tc.Max))
		})
	}
}

func TestMaxInts(t *testing.T) {
	testCases := []struct {
		A        int
		B        int
		Expected int
	}{
		{0, 100, 100},
		{101, 100, 101},
		{100, 0, 100},
		{-100, 100, 100},
		{-1, -2, -1},
	}

	for k, tc := range testCases {
		t.Run(cast.ToString(k+1), func(t *testing.T) {
			require.Equal(t, tc.Expected, MaxInt(tc.A, tc.B))
		})
	}
}

func TestMinInts(t *testing.T) {
	testCases := []struct {
		A        int
		B        int
		Expected int
	}{
		{0, 100, 0},
		{101, 100, 100},
		{100, 0, 0},
		{-100, 100, -100},
		{-1, -2, -2},
	}

	for k, tc := range testCases {
		t.Run(cast.ToString(k+1), func(t *testing.T) {
			require.Equal(t, tc.Expected, MinInt(tc.A, tc.B))
		})
	}
}

func TestDefaultMaxInts(t *testing.T) {
	testCases := []struct {
		Default  int
		Max      int
		Val      int
		Expected int
	}{
		{10, 100, 50, 50},
		{10, 75, 50, 50},
		{25, 75, 50, 50},
		{25, 75, 0, 25},
		{25, 75, 76, 75},
		{55, 75, 50, 50},
		{55, 75, 100, 75},
		{-1, -2, 100, -2},
		{-1, -2, -100, -1},
		{-2, -1, -100, -2},
	}

	for k, tc := range testCases {
		t.Run(cast.ToString(k+1), func(t *testing.T) {
			require.Equal(t, tc.Expected, DefaultMaxInt(tc.Val, tc.Default, tc.Max))
		})
	}
}
