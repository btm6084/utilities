package conv

import (
	"errors"
	"fmt"
	"time"

	"github.com/spf13/cast"
)

func DaysSince(t time.Time) int {
	return cast.ToInt(time.Since(t).Hours() / 24)
}

func DiffInDays(a, b time.Time) int {
	return int(a.Sub(b).Hours() / 24)
}

// ParseTime attempts to extract a valid time.Time from an input string. We prefer RFC3999.
// Fallback should match MSSQL's CURRENT_TIMESTAMP.
func ParseTime(input string) (time.Time, error) {
	var layouts = []string{
		time.RFC3339,
		time.RFC1123,
		"2006-01-02T15:04:05",        // MSSQL Non-Space format.
		"2006-01-02 15:04:05",        // MSSQL With-Space format.
		"2006-01-02 15:04:05 -07:00", // MSSQL With-Space format and timezone.
		"2006-01-02",
		"2006.01.02",
		"2006/01/02",
		"01-02-2006",
		"01.02.2006",
		"01/02/2006",
		"06-01-02",
		"06.01.02",
		"06/01/02",
		"01-02-06",
		"01.02.06",
		"01/02/06",
	}

	for _, l := range layouts {
		out, err := time.Parse(l, input)
		if err == nil {
			return out, nil
		}
	}

	return time.Time{}, errors.New("Unable to parse timestamp '" + input + "'")
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
