package permutation

var (
	overflow = 10
)

// Strings returns a set of all the permutations of the input slice.
func Strings(in []string) [][]string {
	if len(in) < 2 {
		return [][]string{in}
	}

	// Overflow
	if len(in) > overflow {
		return nil
	}

	n := len(in)
	p := n
	for i := n - 1; i > 0; i-- {
		p *= i
	}

	output := make([][]string, p)
	for i := 0; i < p; i++ {
		output[i] = make([]string, n)
	}

	// The smallest base permutation is 2x2
	output[0][0], output[0][1] = in[0], in[1]
	output[1][0], output[1][1] = in[1], in[0]

	// acc is an accumulator, which aggregates the current set length.
	acc := 2
	for i := 2; i < n; i++ {
		// For each N, copy the last set to the new home
		k := 0
		setLen := acc
		acc *= i + 1
		for j := setLen; j < acc; j++ {
			copy(output[j], output[k])
			k++
			k = k % setLen
		}

		// For each copy set, insert 1 set of the new element into each position.
		for j := 0; j < acc; j++ {
			k := j / setLen
			copy(output[j][k+1:], output[j][k:])
			output[j][k] = in[i]
		}
	}

	return output
}

// StringsRecursive builds the permutation list recursively
func StringsRecursive(in []string) [][]string {
	// Overflow
	if len(in) > overflow {
		return nil
	}

	return permStrings(in, []string{}, [][]string{})
}

func permStrings(in []string, path []string, acc [][]string) [][]string {
	if acc == nil {
		return acc
	}

	if len(in) == 0 {
		p := make([]string, len(path))
		copy(p, path)
		acc = append(acc, p)
	}

	for i := 0; i < len(in); i++ {
		acc = permStrings(deleteStrings(in, i), append(path, in[i]), acc)
	}

	return acc
}

// delete makes a copy of the incoming slice that excludes the given index.
func deleteStrings(in []string, idx int) []string {
	var out = make([]string, len(in)-1)

	b := 0
	for i := 0; i < len(in); i++ {
		if i == idx {
			continue
		}

		out[b] = in[i]
		b++
	}

	return out
}
