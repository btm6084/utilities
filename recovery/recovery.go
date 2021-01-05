package recovery

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
)

// PanicRecovery returns a general use Panic Recovery function to capture panics,
// optionally log them, and return them, as errors. A pointer to the error
// to populate will be passed in via the err parameter.
//
// usage: defer PanicRecovery(&err, true)()
func PanicRecovery(err *error, logErrors bool) func() {
	return func() {

		// If err is nil, we don't want to cause a panic. However, we can't actually set it
		// and have it persist beyond this function. Therefore, we force logging on so it
		// isn't just ignored completely if we receive a nil pointer.
		if err == nil {
			logErrors = true
			err = new(error)
		}

		if r := recover(); r != nil && err != nil {
			s, i, _ := identifyPanic()

			switch r.(type) {
			case error:
				if logErrors {
					log.WithFields(log.Fields{"panic": "error", "file": s, "line_num": i}).Error(r.(error))
				}

				// Create a new error and assign it to our pointer.
				*err = r.(error)
			case string:
				if logErrors {
					log.WithFields(log.Fields{"panic": "string", "file": s, "line_num": i}).Error(r.(string))
				}

				// Create a new error and assign it to our pointer.
				*err = errors.New(r.(string))
			default:
				msg := fmt.Sprintf("%+v", r)

				if logErrors {
					log.WithFields(log.Fields{"panic": "default", "file": s, "line_num": i}).Error(msg)
				}

				// Create a new error and assign it to our pointer.
				*err = errors.New("Panic: " + msg)
			}
		}
	}
}

// https://gist.github.com/swdunlop/9629168
func identifyPanic() (string, int, error) {
	var name, file string
	var line int
	var pcs [16]uintptr
	var pc uintptr

	n := runtime.Callers(3, pcs[:])

	// Find the first stack entry that isn't a runtime
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			continue
		}
		file, line = fn.FileLine(pc)
		name = fn.Name()
		if !strings.HasPrefix(name, "runtime.") {
			break
		}
	}

	if file != "" {
		return file, line, nil
	}

	return "", -1, fmt.Errorf("pc:%x", pc)
}

// PanicHandler creates an handler that intercepts panics
func PanicHandler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			defer PanicRecovery(nil, true)()

			next.ServeHTTP(w, req)
		})
	}
}
