package permutation

import "testing"

var (
	inputStrings = []string{"a", "b", "c", "d", "e", "f", "g", "h"}
	inputInts    = []int{1, 2, 3, 4, 5, 6, 7, 8}
)

/*

Thu 12 Nov 2020 12:53:36 PM CST

go test -bench=. -benchmem
goos: linux
goarch: amd64
pkg: github.com/btm6084/utilities/permutation
BenchmarkPermuteLoopStrings-4        	   63946	     17348 ns/op	   12672 B/op	     121 allocs/op
BenchmarkPermuteRecursiveStrings-4   	   27374	     41694 ns/op	   31168 B/op	     418 allocs/op
BenchmarkPermuteLoopInts-4           	   97945	     10705 ns/op	    8832 B/op	     121 allocs/op
BenchmarkPermuteRecursiveInts-4      	   50305	     24007 ns/op	   18808 B/op	     418 allocs/op
PASS
ok  	github.com/btm6084/utilities/permutation	5.523s

*/

func BenchmarkPermuteLoopStrings(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Strings(inputStrings)
	}
}

func BenchmarkPermuteRecursiveStrings(b *testing.B) {
	for i := 0; i < b.N; i++ {
		StringsRecursive(inputStrings)
	}
}

func BenchmarkPermuteLoopInts(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Ints(inputInts)
	}
}

func BenchmarkPermuteRecursiveInts(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IntsRecursive(inputInts)
	}
}
