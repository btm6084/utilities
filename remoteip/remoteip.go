package remoteip

import (
	"net"
	"net/http"
	"strings"
)

// Get returns the remote IP address from the request.
func Get(r *http.Request) string {
	var raw string

	switch {
	case r.Header.Get("X-Real-IP") != "":
		raw = r.Header.Get("X-Real-IP")
	case r.Header.Get("X-Forwarded-For") != "":
		raw = r.Header.Get("X-Forwarded-For")
	case r.Header.Get("X-Client-IP") != "":
		raw = r.Header.Get("X-Client-IP")
	case r.RemoteAddr != "":
		raw = r.RemoteAddr
	}

	pieces := strings.Split(raw, ",")
	candidate := strings.Trim(pieces[0], `'" `)

	var host string
	host, _, err := net.SplitHostPort(candidate)
	if err != nil {
		host = candidate
	}

	ip := net.ParseIP(host)
	if ip == nil {
		return "-"
	}

	return ip.String()
}
