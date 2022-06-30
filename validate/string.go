package validate

import "regexp"

var (
	// IDNumber matches a single integer, potentially negative.
	positiveNumeric    = regexp.MustCompile("^0$|^([1-9][0-9]*)$")
	numeric            = regexp.MustCompile("^-?[1-9][0-9]*$")
	numericLeadingZero = regexp.MustCompile("^-?[0-9]+$")

	// idWord matches a single word-character identifier. Allows spaces.
	idWord = regexp.MustCompile(`^[\w ]+$`)

	alpha    = regexp.MustCompile(`^[A-Za-z]+$`)
	alphaNum = regexp.MustCompile(`^[A-Za-z0-9]+$`)
)

func IsInteger(in string, allowLeadingZeros bool) bool {
	if in == "0" {
		return true
	}

	if allowLeadingZeros {
		return numericLeadingZero.MatchString(in)
	}

	return numeric.MatchString(in)
}

// IsPositiveInteger returns true if the input string represents a positive integer
func IsPositiveInteger(in string) bool {
	if in == "0" {
		return true
	}

	return positiveNumeric.MatchString(in)
}

func IsIdentifier(in string) bool {
	return idWord.MatchString(in)
}

func IsAlpha(in string) bool {
	return alpha.MatchString(in)
}

func IsAlphaNum(in string) bool {
	return alphaNum.MatchString(in)
}
