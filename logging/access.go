package logging

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/btm6084/utilities/remoteip"
)

// NoopWriter satisfies the io.Writer interface, but does nothing.
type NoopWriter struct{}

func (*NoopWriter) Write([]byte) (int, error) {
	return 0, nil
}

func first(in []string) string {
	if len(in) > 0 {
		return in[0]
	}

	return "-"
}

func escape(in string) string {
	b := []byte(strconv.Quote(in))
	if len(b) > 2 {
		return string(b[1 : len(b)-1])
	}

	return "-"
}

type responseWriter struct {
	w      http.ResponseWriter
	status int
	length int
}

func (l *responseWriter) Header() http.Header { return l.w.Header() }

func (l *responseWriter) Write(data []byte) (int, error) {
	n, err := l.w.Write(data)
	l.length += n
	return n, err
}

func (l *responseWriter) WriteHeader(status int) {
	l.status = status
	l.w.WriteHeader(status)
}

type logWriter struct {
	hostname string
	ip       net.IP
	port     string
	logger   io.Writer
}

func (l logWriter) logRequest(req *http.Request, start time.Time, dur time.Duration, rw *responseWriter, pretty bool) {
	var username = "-"
	if req.URL.User != nil {
		if name := req.URL.User.Username(); name != "" {
			username = name
		}
	}

	// Parse the query string to JSON
	q := `"` + req.URL.RawQuery + `"`
	b, err := json.Marshal(req.URL.Query())
	if err == nil {
		q = string(b)
	}

	out := map[string]interface{}{
		"clientIP":              escape(remoteip.Get(req)),
		"contentType":           rw.w.Header().Get("Content-Type"),
		"cookies":               escape(first(req.Header["Cookie"])),
		"date":                  start.Format("2006-01-02"),
		"duration":              dur / time.Millisecond,
		"fromCache":             CacheStatusFromHeaders(rw.w.Header()),
		"httpStatusCode":        rw.status,
		"method":                escape(req.Method),
		"path":                  escape(req.URL.Path),
		"queryString":           json.RawMessage(q),
		"referrer":              escape(first(req.Header["Referer"])),
		"requestContentLength":  req.ContentLength,
		"responseContentLength": rw.length,
		"serverIP":              escape(l.ip.String()),
		"serverPort":            escape(l.port),
		"time":                  start.Format("15:04:05.000"),
		"txnID":                 TransactionFromContext(req.Context()),
		"userAgent":             escape(first(req.Header["User-Agent"])),
		"username":              escape(username),
	}

	raw, _ := json.Marshal(out)
	if pretty {
		raw, _ = json.MarshalIndent(out, "", "\t")
	}
	fmt.Fprintln(l.logger, string(raw))
}

// CreateLogger creates an access logger
func CreateLogger(logger io.Writer, listenPort int, pretty bool) func(http.Handler) http.Handler {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "-"
	}

	lw := logWriter{hostname: hostname, logger: logger, ip: GetOutboundIP(), port: strconv.Itoa(listenPort)}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			start := time.Now()
			rw := &responseWriter{w, 0, 0}
			next.ServeHTTP(rw, req)

			// Check for a timeout error, and make sure its status gets logged.
			ctxErr := req.Context().Err()
			if ctxErr != nil && req.Context().Err().Error() == "context deadline exceeded" {
				rw.status = http.StatusServiceUnavailable
				rw.length = 0
			}

			lw.logRequest(req, start, time.Since(start), rw, pretty)
		})
	}
}

// GetOutboundIP returns the outbound ip address.
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	addr := conn.LocalAddr().(*net.UDPAddr)

	return addr.IP
}
