package permutation

import "testing"

var (
	input = []string{"a", "b", "c", "d", "e"}
)

/*

Thu 12 Nov 2020 12:53:36 PM CST

go test -bench=. -benchmem -benchtime=10s
goos: linux
goarch: amd64
pkg: github.com/btm6084/utilities/permutation
BenchmarkPermuteLoop-4              	  660688	     30304 ns/op	   12672 B/op	     121 allocs/op
BenchmarkPermuteRecursive-4         	  240336	     54207 ns/op	   31168 B/op	     418 allocs/op
BenchmarkPermuteRecursiveAppend-4   	  265472	     46140 ns/op	   30848 B/op	     413 allocs/op
PASS
ok  	github.com/btm6084/utilities/permutation	46.458s

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

func BenchmarkPermuteRecursiveAppend(b *testing.B) {
	for i := 0; i < b.N; i++ {
		StringsRecursiveAppend(input)
	}
}
