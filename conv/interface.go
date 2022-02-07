package conv

import (
	"reflect"

	"github.com/btm6084/utilities/recovery"
)

// Returns true if the given interface is nil.
func IsNil(i interface{}) (isNil bool) {
	var err error

	// IsNil panics on non-nilables, such as atomics. This will catch those and mark them not nil.
	defer func() {
		if err != nil {
			isNil = false
			return
		}
	}()
	defer recovery.PanicRecovery(&err, false, "")()

	if i == nil {
		isNil = true
		return
	}
	isNil = reflect.ValueOf(i).IsNil()
	return
}
