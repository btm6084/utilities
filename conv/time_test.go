package conv

import (
	"testing"
	"time"

	"github.com/spf13/cast"
	"github.com/stretchr/testify/require"
)

func TestDaysSince(t *testing.T) {
	testCases := []struct {
		A        time.Time
		Expected int
	}{
		{time.Now(), 0},
	}

	for k, tc := range testCases {
		t.Run(cast.ToString(k+1), func(t *testing.T) {
			require.Equal(t, tc.Expected, DaysSince(tc.A))
		})
	}
}

func TestDiffInDays(t *testing.T) {
	testCases := []struct {
		A        time.Time
		B        time.Time
		Expected int
	}{
		{time.Now(), time.Now(), 0},
	}

	for k, tc := range testCases {
		t.Run(cast.ToString(k+1), func(t *testing.T) {
			require.Equal(t, tc.Expected, DiffInDays(tc.A, tc.B))
		})
	}
}
