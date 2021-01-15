package cache

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

var (
	forbiddenHeaders = map[string]bool{
		"Access-Control-Allow-Credentials": true,
		"Access-Control-Allow-Headers":     true,
		"Access-Control-Allow-Methods":     true,
		"Access-Control-Allow-Origin":      true,
		"X-Cache-Status":                   true,
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
func HandlerWrapper(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			next.ServeHTTP(w, r)
			return
		}

		key := r.Method + r.RequestURI

		if b, ok := Get(key); ok {
			if s, ok := Get(key + "StatusCode"); ok {
				w.WriteHeader(s.(int))
			} else {
				w.WriteHeader(200)
			}

			// Retain any headers.
			if h, ok := Get(key + "headers"); ok {
				if headers, ok := h.(http.Header); ok {
					for k, v := range headers {
						if forbiddenHeader(k) {
							continue
						}

						for i := 0; i < len(v); i++ {
							w.Header().Set(k, v[i])
						}
					}
				}
			}

			w.Header().Set("X-Cache-Status", "hit")
			w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d, public", int(dur/time.Second)))
			w.Write(b.([]byte))
			return
		}

		writer := ResponseWriterTee{w: w}
		w.Header().Set("X-Cache-Status", "miss")
		w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d, public", int(dur/time.Second)))
		next.ServeHTTP(&writer, r)

		sc := writer.StatusCode
		if sc >= 500 {
			return
		}

		Set(key, writer.Buffer.Bytes())
		Set(key+"headers", w.Header())
		Set(key+"StatusCode", writer.StatusCode)
	})
}

func forbiddenHeader(k string) bool {
	_, isset := forbiddenHeaders[k]
	return isset
}
