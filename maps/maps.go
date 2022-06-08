package maps

func Keys[K comparable, V any](in map[K]V) []K {
	out := make([]K, len(in))

	i := 0
	for k := range in {
		out[i] = k
		i++
	}

	return out
}

func Values[K comparable, V any](in map[K]V) []V {
	out := make([]V, len(in))

	i := 0
	for _, v := range in {
		out[i] = v
		i++
	}

	return out
}
