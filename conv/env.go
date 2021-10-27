package conv

import (
	"os"
	"strings"
)

func Env(key string) string {
	if os.Getenv(key) != "" {
		return os.Getenv(key)
	}

	key = strings.ReplaceAll(key, "_", "")
	return os.Getenv(key)
}
