package stack

import (
	"bytes"
	"regexp"
	"runtime/debug"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
)

var (
	fileLineRE = regexp.MustCompile(`^([^:]+):([0-9]+)`)
	truncateRE = regexp.MustCompile(`/(pkg|cmd)/[^:]+`)
)

// TraceFields builds a log.Fields object when all you need is a stack trace.
func TraceFields() log.Fields {

	stack := bytes.Split(debug.Stack(), []byte{'\n'})
	var f string
	var l int

	for i, s := range stack {
		if strings.Contains(string(s), "github.com/btm6084/utilities/stack.TraceFields") {
			if i+3 >= len(stack) {
				return log.Fields{"stacktrace": Strings()}
			}

			matches := fileLineRE.FindAllStringSubmatch(strings.TrimSpace(string(stack[i+3])), 1)

			if len(matches) < 1 {
				return log.Fields{"stacktrace": Strings()}
			}

			if fa := truncateRE.FindString(matches[0][1]); fa != "" {
				f = fa
			} else {
				f = matches[0][1]
			}

			l = cast.ToInt(matches[0][2])
		}
	}

	if f == "" {
		return log.Fields{"stacktrace": Strings()}
	}

	return log.Fields{"stacktrace": map[string]interface{}{"file": f, "line": l}}
}

func Trace(depth int) (string, int) {
	stack := bytes.Split(debug.Stack(), []byte{'\n'})
	start := 0
	for i, s := range stack {
		if strings.Contains(string(s), "github.com/btm6084/utilities/stack.Trace") {
			start = i
			break
		}
	}

	depth = depth + 1
	index := start + (2 * depth) + 1
	if index >= len(stack) {
		index = len(stack) - 1
	}

	tmp := strings.TrimSpace(string(stack[index]))
	tmp = strings.Split(tmp, " ")[0]
	pieces := strings.Split(tmp, ":")
	return string(pieces[0]), cast.ToInt(pieces[1])
}

func String() string {
	return strings.Join(Strings(), "\n")
}

func Strings() []string {
	stack := bytes.Split(debug.Stack(), []byte{'\n'})
	out := make([]string, len(stack))

	for i, s := range stack {
		out[i] = string(s)
	}

	return out
}
