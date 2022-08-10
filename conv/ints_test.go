package conv

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

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

func TestFuzzInt(t *testing.T) {
	rand.Seed(time.Now().UnixMicro())
	var asPct = func(in int) float64 {
		return ToFixed(float64(in)/100, 2)
	}

	for i := 0; i < 10000; i++ {
		min := rand.Intn(100)
		max := rand.Intn(100-min) + min

		in := 100
		out := FuzzInt(in, min, max)

		var pct float64
		var diff float64
		if in > out { // Minus Branch
			diff = float64(in) - float64(out)
			pct = ToFixed(float64(diff)/float64(in), 2)
		} else { // Plus Branch
			diff = float64(out) - float64(in)
			pct = ToFixed(float64(diff)/float64(in), 2)
		}

		require.True(t, pct >= asPct(min) && pct <= asPct(max), fmt.Sprintf("(%f, %f) Got: %f", asPct(min), asPct(max), pct))
	}
}
