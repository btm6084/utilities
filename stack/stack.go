package stack

import (
	"runtime"
	"strings"
)

func Trace(depth int) (string, int) {
	var pcs [16]uintptr
	n := runtime.Callers(0, pcs[:])

	if depth < 0 {
		return "", 0
	}

	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		n := fn.Name()

		switch {
		case strings.HasPrefix(n, "runtime."):
			continue
		case depth > 0:
			depth--
		default:
			return fn.FileLine(pc)
		}
	}

	return "", 0
}
