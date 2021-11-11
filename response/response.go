// Package response provides a set of utilities for writing responses to an http.ResponseWriter.
package response

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// ServeJSON serves the supplied data as JSON
func ServeJSON(w http.ResponseWriter, r *http.Request, statusCode int, data interface{}) {
	serveJSON(w, r, statusCode, data, false)
}

// ServeETagJSON serves the supplied data as JSON and includes an ETag for 2xx status codes
func ServeETagJSON(w http.ResponseWriter, r *http.Request, statusCode int, data interface{}) {
	serveJSON(w, r, statusCode, data, true)
}

func serveJSON(w http.ResponseWriter, r *http.Request, statusCode int, data interface{}, etag bool) {
	w.Header().Set("Content-Type", "application/json")

	buf := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)

	err := enc.Encode(data)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b := buf.Bytes()

	if etag {
		sum := md5.Sum(b)
		etag := hex.EncodeToString(sum[:])

		if statusCode/100 == 2 {
			match := r.Header.Get("If-None-Match")
			if etag == match {
				w.WriteHeader(http.StatusNotModified)
				return
			}

			w.Header().Set("ETag", etag)
		}
	}

	age := int(time.Minute * 1 / time.Second)
	w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d, public", age))

	w.WriteHeader(statusCode)
	w.Write(b)
}

// ServeNoContent serves the status code only, with no content body or Content-Type headers.
func ServeNoContent(w http.ResponseWriter, r *http.Request, statusCode int) {
	w.WriteHeader(statusCode)
}

// ServeString serves the data parameter as the output from the request with no processing on the data.
func ServeString(w http.ResponseWriter, r *http.Request, statusCode int, data string) {
	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "text/plain")
	}

	w.WriteHeader(statusCode)
	w.Write([]byte(data))
}
