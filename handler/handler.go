package handler

import (
	"fmt"
	"net/http"
	"strconv"
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

// HeadResponseWriter captures input to an http.ResponseWriter
type HeadResponseWriter struct {
	Bytes      int
	StatusCode int
	w          http.ResponseWriter
}

// Header proxies http.ResponseWriter Header
func (w *HeadResponseWriter) Header() http.Header {
	return w.w.Header()
}

// WriteHeader proxies http.ResponseWriter WriteHeader
func (w *HeadResponseWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
	w.w.WriteHeader(statusCode)
}

// Write prevents http.ResponseWriter Write from writing anything to the response body for HEAD requests.
func (w *HeadResponseWriter) Write(b []byte) (int, error) {
	w.Bytes += len(b)
	return len(b), nil
}

// HeadHandler returns a handler that prevents a body from being returned for any wrapped request handler if the HTTP method is HEAD.
func HeadHandler(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "HEAD" {
			next.ServeHTTP(w, r)
			return
		}

		writer := HeadResponseWriter{w: w}
		next.ServeHTTP(&writer, r)
		writer.Header().Set("Content-Length", strconv.Itoa(writer.Bytes))
	})
}
