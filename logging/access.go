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
	"strings"
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

func firstNonEmptyString(list ...string) string {
	for _, v := range list {
		if v != "" {
			return v
		}
	}

	return ""
}

func UserAgent(req *http.Request) string {
	ua := firstNonEmptyString(req.Header.Get("X-User-Agent"), req.Header.Get("User-Agent"), "-")
	return strings.Trim(ua, `'" `)
}

func escape(in string) string {
	b := []byte(strconv.Quote(in))
	if len(b) > 2 {
		return string(b[1 : len(b)-1])
	}

	return "-"
}

type ResponseWriter struct {
	w      http.ResponseWriter
	status int
	length int
}

func (l *ResponseWriter) Header() http.Header { return l.w.Header() }

func (l *ResponseWriter) Write(data []byte) (int, error) {
	n, err := l.w.Write(data)
	l.length += n
	return n, err
}

func (l *ResponseWriter) WriteHeader(status int) {
	l.status = status
	l.w.WriteHeader(status)
}

type Logger interface {
	LogRequest(req *http.Request, start time.Time, dur time.Duration, rw *ResponseWriter, pretty bool)
}

type LogWriter struct {
	Hostname string
	IP       net.IP
	Port     string
	Logger   io.Writer
}

type AccessLog struct {
	ClientIP              string          `json:"cIP,omitempty"`
	ContentType           string          `json:"-"`
	Cookies               string          `json:"-"`
	Date                  string          `json:"date,omitempty"`
	Duration              time.Duration   `json:"dur,omitempty"`
	FromCache             bool            `json:"cache,omitempty"`
	HttpStatusCode        int             `json:"code,omitempty"`
	Method                string          `json:"method,omitempty"`
	Path                  string          `json:"path,omitempty"`
	QueryString           json.RawMessage `json:"qs,omitempty"`
	Referrer              string          `json:"rfer,omitempty"`
	RequestContentLength  int             `json:"rqCL,omitempty"`
	ResponseContentLength int             `json:"rsCL,omitempty"`
	ServerHN              string          `json:"sHN,omitempty"`
	ServerIP              string          `json:"sIP,omitempty"`
	ServerPort            string          `json:"sPT,omitempty"`
	Time                  string          `json:"time,omitempty"`
	TxnID                 string          `json:"txnID,omitempty"`
	UserAgent             string          `json:"ua,omitempty"`
	Username              string          `json:"-"`
}

func (l LogWriter) LogRequest(req *http.Request, start time.Time, dur time.Duration, rw *ResponseWriter, pretty bool) {
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

	out := AccessLog{
		ClientIP:              remoteip.Get(req),
		ContentType:           rw.w.Header().Get("Content-Type"),
		Cookies:               first(req.Header["Cookie"]),
		Date:                  start.Format("2006-01-02"),
		Duration:              dur / time.Millisecond,
		FromCache:             CacheStatusFromHeaders(rw.w.Header()),
		HttpStatusCode:        rw.status,
		Method:                req.Method,
		Path:                  req.URL.Path,
		QueryString:           json.RawMessage(q),
		Referrer:              first(req.Header["Referer"]),
		RequestContentLength:  int(req.ContentLength),
		ResponseContentLength: rw.length,
		ServerHN:              l.Hostname,
		ServerIP:              l.IP.String(),
		ServerPort:            l.Port,
		Time:                  start.Format("15:04:05.000"),
		TxnID:                 TransactionFromContext(req.Context()),
		UserAgent:             UserAgent(req),
		Username:              username,
	}

	raw, _ := json.Marshal(out)
	if pretty {
		raw, _ = json.MarshalIndent(out, "", "\t")
	}
	fmt.Fprintln(l.Logger, string(raw))
}

// CreateLogger creates an access logger
func CreateLogger(logger io.Writer, listenPort int, pretty bool) func(http.Handler) http.Handler {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "-"
	}

	lw := LogWriter{Hostname: hostname, Logger: logger, IP: GetOutboundIP(), Port: strconv.Itoa(listenPort)}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			start := time.Now()
			rw := &ResponseWriter{w, 0, 0}
			next.ServeHTTP(rw, req)

			// Check for a timeout error, and make sure its status gets logged.
			ctxErr := req.Context().Err()
			if ctxErr != nil && req.Context().Err().Error() == "context deadline exceeded" {
				rw.status = http.StatusServiceUnavailable
				rw.length = 0
			}

			lw.LogRequest(req, start, time.Since(start), rw, pretty)
		})
	}
}

// CustomAccessLog creates a custom access logger
func CustomAccessLog(logger Logger, pretty bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			start := time.Now()
			rw := &ResponseWriter{w, 0, 0}
			next.ServeHTTP(rw, req)

			// Check for a timeout error, and make sure its status gets logged.
			ctxErr := req.Context().Err()
			if ctxErr != nil && req.Context().Err().Error() == "context deadline exceeded" {
				rw.status = http.StatusServiceUnavailable
				rw.length = 0
			}

			logger.LogRequest(req, start, time.Since(start), rw, pretty)
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
