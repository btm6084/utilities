package validate

import (
	"testing"

	"github.com/spf13/cast"
	"github.com/stretchr/testify/require"
)

func TestIsAlpha(t *testing.T) {
	var testCases = []struct {
		Input    string
		Expected bool
	}{
		{"Hello", true},
		{"lskdhflshdfasiyhdflsakdhflaskdhvasdlfv", true},
		{"ls12kdhfl3shdfasiyhdfl32213123sakdh2flaskdh3vasdlfv", false},
		{"lskdhfls h d fasiyhdflsakdhflaskdhvasdlfv", false},
		{"Hello!", false},
		{"Alpha Numeric 123", false},
		{"3728123123", false},
		{"-3728123123", false},
		{"-0728123123", false},
		{"003728123123", false},
		{"_003728123123", false},
	}

	for n, tc := range testCases {
		t.Run("IsAlpha "+cast.ToString(n), func(t *testing.T) {
			require.Equal(t, tc.Expected, IsAlpha(tc.Input))
		})
	}
}

func TestIsAlphaNumeric(t *testing.T) {
	var testCases = []struct {
		Input    string
		Expected bool
	}{
		{"Hello", true},
		{"lskdhflshdfasiyhdflsakdhflaskdhvasdlfv", true},
		{"ls12kdhfl3shdfasiyhdfl32213123sakdh2flaskdh3vasdlfv", true},
		{"lskdhfls h d fasiyhdflsakdhflaskdhvasdlfv", false},
		{"Hello!", false},
		{"Alpha Numeric 123", false},
		{"3728123123", true},
		{"-3728123123", false},
		{"-0728123123", false},
		{"003728123123", true},
		{"_003728123123", false},
	}

	for n, tc := range testCases {
		t.Run("IsAlphaNum "+cast.ToString(n), func(t *testing.T) {
			require.Equal(t, tc.Expected, IsAlphaNum(tc.Input))
		})
	}
}

func TestIsIntegerNoLeadingZeros(t *testing.T) {
	var testCases = []struct {
		Input    string
		Expected bool
	}{
		{"Hello", false},
		{"lskdhflshdfasiyhdflsakdhflaskdhvasdlfv", false},
		{"ls12kdhfl3shdfasiyhdfl32213123sakdh2flaskdh3vasdlfv", false},
		{"lskdhfls h d fasiyhdflsakdhflaskdhvasdlfv", false},
		{"Hello!", false},
		{"Alpha Numeric 123", false},
		{"3728123123", true},
		{"-3728123123", true},
		{"-0728123123", false},
		{"003728123123", false},
		{"_003728123123", false},
	}

	for n, tc := range testCases {
		t.Run("IsIntegerNoZeros "+cast.ToString(n), func(t *testing.T) {
			require.Equal(t, tc.Expected, IsInteger(tc.Input, false))
		})
	}
}

func TestIsIntegerLeadingZeros(t *testing.T) {
	var testCases = []struct {
		Input    string
		Expected bool
	}{
		{"Hello", false},
		{"lskdhflshdfasiyhdflsakdhflaskdhvasdlfv", false},
		{"ls12kdhfl3shdfasiyhdfl32213123sakdh2flaskdh3vasdlfv", false},
		{"lskdhfls h d fasiyhdflsakdhflaskdhvasdlfv", false},
		{"Hello!", false},
		{"Alpha Numeric 123", false},
		{"3728123123", true},
		{"-3728123123", true},
		{"-0728123123", true},
		{"003728123123", true},
		{"_003728123123", false},
	}

	for n, tc := range testCases {
		t.Run("IsIntegerWithZeros "+cast.ToString(n), func(t *testing.T) {
			require.Equal(t, tc.Expected, IsInteger(tc.Input, true))
		})
	}
}

func TestIsPositiveInteger(t *testing.T) {
	var testCases = []struct {
		Input    string
		Expected bool
	}{
		{"Hello", false},
		{"lskdhflshdfasiyhdflsakdhflaskdhvasdlfv", false},
		{"ls12kdhfl3shdfasiyhdfl32213123sakdh2flaskdh3vasdlfv", false},
		{"lskdhfls h d fasiyhdflsakdhflaskdhvasdlfv", false},
		{"Hello!", false},
		{"Alpha Numeric 123", false},
		{"3728123123", true},
		{"0", true},
		{"1", true},
		{"01", false},
		{"-3728123123", false},
		{"-0728123123", false},
		{"003728123123", false},
		{"_003728123123", false},
	}

	for n, tc := range testCases {
		t.Run("IsIntegerWithZeros "+cast.ToString(n), func(t *testing.T) {
			require.Equal(t, tc.Expected, IsPositiveInteger(tc.Input), tc.Input)
		})
	}
}

func TestIsIdentifier(t *testing.T) {
	var testCases = []struct {
		Input    string
		Expected bool
	}{
		{"Hello", true},
		{"lskdhflshdfasiyhdflsakdhflaskdhvasdlfv", true},
		{"ls12kdhfl3shdfasiyhdfl32213123sakdh2flaskdh3vasdlfv", true},
		{"lskdhfls h d fasiyhdflsakdhflaskdhvasdlfv", true},
		{"Hello!", false},
		{"Alpha Numeric 123", true},
		{"3728123123", true},
		{"-3728123123", false},
		{"-0728123123", false},
		{"003728123123", true},
		{"_003728123123", true},
	}

	for n, tc := range testCases {
		t.Run("IsIdentifier "+cast.ToString(n), func(t *testing.T) {
			require.Equal(t, tc.Expected, IsIdentifier(tc.Input))
		})
	}
}
