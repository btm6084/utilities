package handler

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

// CreateOptionsHandler creates a generic options handler that takes a list of allowed methods.
func CreateOptionsHandler(methods ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		methodString := strings.Join(methods, ", ")

		header := w.Header()
		header.Set("Access-Control-Allow-Headers", "Authorization, Cache-Control, Content-Type")
		header.Set("Access-Control-Allow-Methods", methodString)

		age := int(time.Hour * 1 / time.Second)
		w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d, public", age))
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusNoContent)
	}
}
