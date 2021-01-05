package slice

import (
	"sort"

	"github.com/spf13/cast"
)

// UniqueInt takes a slice of int and returns a slice with only unique elements.
func UniqueInt(in []int) []int {
	uMap := make(map[int]bool)
	for _, v := range in {
		uMap[v] = true
	}

	out := make([]int, len(uMap))
	i := 0
	for k := range uMap {
		out[i] = k
		i++
	}

	return out
}

// UniqueString takes a slice of string and returns a slice with only unique elements.
func UniqueString(in []string, allowEmpty bool) []string {
	seen := make(map[string]bool)
	var out []string
	for _, v := range in {
		if _, isset := seen[v]; isset {
			continue
		}

		if !allowEmpty && v == "" {
			continue
		}

		seen[v] = true
		out = append(out, v)
	}

	return out
}

// StringToInt takes a slice of string and returns a slice of int.
func StringToInt(in []string) []int {
	out := make([]int, len(in))
	for i := 0; i < len(in); i++ {
		out[i] = cast.ToInt(in[i])
	}

	return out
}

// ContainsString returns true if string needle is an element in string slice haystack
func ContainsString(haystack []string, needle string) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}
	return false
}

// ContainsInt returns true if int needle is an element in int slice haystack
func ContainsInt(haystack []int, needle int) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}
	return false
}

// ContainsFloat returns true if float needle is an element in float slice haystack
func ContainsFloat(haystack []float64, needle float64) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}
	return false
}

// SumInts returns the sum of all elements in the given slice of ints.
func SumInts(in []int) int {
	sum := 0
	for _, v := range in {
		sum += v
	}

	return sum
}

// SumFloats returns the sum of all elements in the given slice of float64.
func SumFloats(in []float64) float64 {
	sum := 0.0
	for _, v := range in {
		sum += v
	}

	return sum
}

// MedianInts gets the median number in the given slice of ints.
func MedianInts(in []int) int {
	l := len(in)
	if l == 0 {
		return 0
	}

	if l == 1 {
		return in[0]
	}

	sort.Ints(in)

	if l%2 != 0 {
		return in[l/2]
	}

	pos := l / 2
	return (in[pos-1] + in[pos]) / 2.0
}

// MedianFloats gets the median number in the given slice of float64.
func MedianFloats(in []float64) float64 {
	l := len(in)
	if l == 0 {
		return 0
	}

	if l == 1 {
		return in[0]
	}

	sort.Float64s(in)

	if l%2 != 0 {
		return in[l/2]
	}

	pos := l / 2
	return (in[pos-1] + in[pos]) / 2.0
}
