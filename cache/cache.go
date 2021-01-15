package cache

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/spf13/cast"
)

// ResponseWriterTee captures input to an http.ResponseWriter
type ResponseWriterTee struct {
	Buffer     bytes.Buffer
	StatusCode int
	w          http.ResponseWriter
}

// Header proxies http.ResponseWriter Header
func (w *ResponseWriterTee) Header() http.Header {
	return w.w.Header()
}

// WriteHeader proxies http.ResponseWriter WriteHeader
func (w *ResponseWriterTee) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
	w.w.WriteHeader(statusCode)
}

// Write proxies http.ResponseWriter Write
func (w *ResponseWriterTee) Write(b []byte) (int, error) {
	w.Buffer.Write(b)
	return w.w.Write(b)
}

// Handler caches all interactions with the API based on Method, URI, and status code.
func Handler(dur int, excludedPaths []string) func(http.Handler) http.Handler {
	t := time.Duration(dur)
	c := cache.New(t*time.Second, (2*t)*time.Second)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "GET" || excludedFromCache(r, excludedPaths) || cast.ToBool(r.URL.Query().Get("noCache")) {
				next.ServeHTTP(w, r)
				w.Header().Set("Cache-Control", "no-cache")
				return
			}

			key := r.Method + r.RequestURI

			if b, ok := c.Get(key); ok {
				if s, ok := c.Get(key + "StatusCode"); ok {
					w.WriteHeader(s.(int))
				} else {
					w.WriteHeader(200)
				}

				w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d, public", dur))
				w.Write(b.([]byte))
				return
			}

			writer := ResponseWriterTee{w: w}
			w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d, public", dur))
			next.ServeHTTP(&writer, r)

			sc := writer.StatusCode
			if sc >= 500 {
				return
			}

			c.Set(key, writer.Buffer.Bytes(), 0)
			c.Set(key+"StatusCode", writer.StatusCode, 0)
		})
	}
}

func excludedFromCache(r *http.Request, excludePaths []string) bool {
	for _, p := range excludePaths {
		if strings.HasPrefix(r.URL.Path, p) {
			return true
		}
	}

	return false
}
