package cache

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/btm6084/utilities/logging"
	"github.com/btm6084/utilities/metrics"
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

// HandlerWrapper caches all interactions with the API based on Method, URI, and status code for the wrapped handler.
// cacheDuration is in seconds.
func HandlerWrapper(cacheDuration int, next http.Handler) http.HandlerFunc {
	d := time.Duration(cacheDuration) * time.Second

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m := metrics.GetRecorder(r.Context())

		if r.Method != "GET" {
			next.ServeHTTP(w, r)
			return
		}

		key := r.Method + r.RequestURI + r.Header.Get("range")

		var b string
		if err := Get(m, key, &b); err == nil {
			var s int
			if err := Get(m, key+"StatusCode", &s); err == nil {
				w.WriteHeader(s)
			} else {
				w.WriteHeader(200)
			}

			// Retain any headers.
			var h http.Header
			if err := Get(m, key+"headers", &h); err == nil {
				for k, v := range h {
					if forbiddenHeader(k) {
						continue
					}

					for i := 0; i < len(v); i++ {
						w.Header().Set(k, v[i])
					}
				}
			}

			req := logging.RequestWithCacheStatus(r, true)
			if r != nil && req != nil {
				*r = *req
			}

			w.Header().Set("X-Cache-Hit", "true")
			w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d, public", int(dur/time.Second)))
			w.Write([]byte(b))
			return
		}

		writer := ResponseWriterTee{w: w}
		w.Header().Set("X-Cache-Hit", "false")
		w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d, public", int(dur/time.Second)))

		req := logging.RequestWithCacheStatus(r, false)
		if r != nil && req != nil {
			*r = *req
		}

		next.ServeHTTP(&writer, r)

		sc := writer.StatusCode
		if sc >= 500 {
			return
		}

		SetWithDuration(m, key, string(writer.Buffer.Bytes()), d)
		SetWithDuration(m, key+"headers", w.Header(), d)
		SetWithDuration(m, key+"StatusCode", writer.StatusCode, d)
	})
}

func forbiddenHeader(k string) bool {
	_, isset := forbiddenHeaders[k]
	return isset
}
