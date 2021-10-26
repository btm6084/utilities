package conv

// QuoteString adds quotes to a string.
func QuoteString(s string) string {
	switch {
	case len(s) == 0:
		return `""`
	case s == `"`:
		return `""`
	case len(s) == 1:
		return `"` + s + `"`
	case s[0] == '"' && s[len(s)-1] == '"':
		return s
	default:
		return `"` + s + `"`
	}
}

func FirstNonEmptyString(list ...string) string {
	for _, v := range list {
		if v != "" {
			return v
		}
	}

	return ""
}
