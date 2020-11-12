package permutation

// Strings returns a set of all the permutations of the input slice.
func Strings(in []string) [][]string {
	if len(in) < 2 {
		return [][]string{in}
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

// // Strings returns a set of all the permutations of the input slice.
// func Strings(in []string) [][]string {
// 	if len(in) < 2 {
// 		return [][]string{in}
// 	}

// 	n := len(in)
// 	p := n
// 	for i := n - 1; i > 0; i-- {
// 		p *= i
// 	}

// 	output := make([][]string, p)
// 	for i := 0; i < p; i++ {
// 		output[i] = make([]string, n)
// 	}

// 	// The smallest base permutation is 2x2
// 	output[0] = []string{in[0], in[1]}
// 	output[1] = []string{in[1], in[0]}

// 	// // Fill all slots with repeating pattern. This creates matching pairs, half the solution set.
// 	// for j := 0; j < n; j++ {
// 	// 	slots := int(p / n)

// 	// 	for i := 0; i < p; i++ {
// 	// 		k := (i/slots + j) % n
// 	// 		output[i][k] = in[j]
// 	// 	}
// 	// }

// 	// // Mirror one of every matching pair. This is the second half of the solution set.
// 	// for i := 1; i < p; i += 2 {
// 	// 	for k := 0; k < n; k++ {
// 	// 		output[i-1][k] = output[i][n-k-1]
// 	// 	}
// 	// }

// 	// for i := 0; i < len(output); i++ {
// 	// 	fmt.Println(output[i])
// 	// }

// 	// os.Exit(1)

// 	return output
// }
