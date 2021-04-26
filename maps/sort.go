package maps

import "sort"

type stringInt struct {
	m map[string]int
	s []string
}

func (sm *stringInt) Len() int {
	return len(sm.m)
}

func (sm *stringInt) Less(i, j int) bool {
	return sm.m[sm.s[i]] > sm.m[sm.s[j]]
}

func (sm *stringInt) Swap(i, j int) {
	sm.s[i], sm.s[j] = sm.s[j], sm.s[i]
}

// SortStringIntValue sorts map[string]int keys in descending order of its values, returning the keys.
func SortStringIntValue(m map[string]int) []string {
	var sm stringInt

	sm.m = m
	sm.s = make([]string, len(m))

	i := 0
	for key := range m {
		sm.s[i] = key
		i++
	}

	sort.Sort(&sm)
	return sm.s
}
