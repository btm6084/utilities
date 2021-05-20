package cache

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/btm6084/utilities/logging"
	"github.com/btm6084/utilities/metrics"
	"github.com/spf13/cast"
)

var (
	forbiddenHeaders = map[string]bool{
		"Access-Control-Allow-Credentials": true,
		"Access-Control-Allow-Headers":     true,
		"Access-Control-Allow-Methods":     true,
		"Access-Control-Allow-Origin":      true,
		"X-Cache-Hit":                      true,
	}
)

type CacheEvent struct {
	Content    string      `json:"content"`
	StatusCode int         `json:"statusCode"`
	Headers    http.Header `json:"headers"`
}

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

// Middleware provides a cache-layer middleware for caching the input/output for GET requests.
func Middleware(cacheDuration int, excludedPaths []string) func(http.Handler) http.Handler {
	d := time.Duration(cacheDuration) * time.Second

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !Enabled || r.Method != "GET" || excludedFromCache(r, excludedPaths) || cast.ToBool(r.URL.Query().Get("noCache")) {
				next.ServeHTTP(w, r)
				w.Header().Set("Cache-Control", "no-cache")
				return
			}

			key := r.Method + r.RequestURI + r.Header.Get("range")
			m := metrics.GetRecorder(r.Context())

			if handlerTryCache(w, r, m, key, d) {
				return
			}

			handleCacheableRequest(next, w, r, m, key, d)
		})
	}
}

// HandlerWrapper provices a cache-layer wrapper for a single API route.
func HandlerWrapper(cacheDuration int, next http.Handler) http.HandlerFunc {
	d := time.Duration(cacheDuration) * time.Second

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !Enabled || r.Method != "GET" {
			next.ServeHTTP(w, r)
			return
		}

		m := metrics.GetRecorder(r.Context())
		key := r.Method + r.RequestURI + r.Header.Get("range")

		if handlerTryCache(w, r, m, key, d) {
			return
		}

		handleCacheableRequest(next, w, r, m, key, d)
	})
}

func handlerTryCache(w http.ResponseWriter, r *http.Request, m metrics.Recorder, key string, d time.Duration) bool {

	var ce CacheEvent
	if err := Get(m, key, &ce); err == nil {
		// Retain any headers.
		for k, v := range ce.Headers {
			if forbiddenHeader(k) {
				continue
			}

			for i := 0; i < len(v); i++ {
				w.Header().Set(k, v[i])
			}
		}

		w.Header().Set("X-Cache-Hit", "true")
		w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d, public", int(d/time.Second)))
		w.WriteHeader(ce.StatusCode)
		w.Write([]byte(ce.Content))
		return true
	}

	return false
}

func handleCacheableRequest(next http.Handler, w http.ResponseWriter, r *http.Request, m metrics.Recorder, key string, d time.Duration) {
	writer := ResponseWriterTee{w: w}
	w.Header().Set("X-Cache-Hit", "false")
	w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d, public", int(d/time.Second)))

	req := logging.RequestWithCacheStatus(r, false)
	if r != nil && req != nil {
		*r = *req
	}

	next.ServeHTTP(&writer, r)

	sc := writer.StatusCode
	if sc >= 500 {
		return
	}

	ce := CacheEvent{
		Content:    writer.Buffer.String(),
		Headers:    w.Header(),
		StatusCode: writer.StatusCode,
	}

	SetWithDuration(m, key, ce, d)
}

func excludedFromCache(r *http.Request, excludePaths []string) bool {
	for _, p := range excludePaths {
		if strings.HasPrefix(r.URL.Path, p) {
			return true
		}
	}

	return false
}

func forbiddenHeader(k string) bool {
	_, isset := forbiddenHeaders[k]
	return isset
}
