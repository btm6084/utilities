package permutation

import "testing"

var (
	input = []string{"a", "b", "c", "d", "e"}
)

/*

Thu 12 Nov 2020 12:47:45 PM CST
go test -bench=. -benchmem -benchtime=10s
goos: linux
goarch: amd64
pkg: github.com/btm6084/utilities/permutation
BenchmarkPermuteLoop-4        	  630086	     20750 ns/op	   12672 B/op	     121 allocs/op
BenchmarkPermuteRecursive-4   	  243061	     49218 ns/op	   31168 B/op	     418 allocs/op

*/

func BenchmarkPermuteLoop(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Strings(input)
	}
}

func BenchmarkPermuteRecursive(b *testing.B) {
	for i := 0; i < b.N; i++ {
		StringsRecursive(input)
	}
}
