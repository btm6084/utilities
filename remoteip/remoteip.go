package remoteip

import (
	"net"
	"net/http"
)

// Get returns the remote IP address from the request.
func Get(r *http.Request) string {
	var raw string

	switch {
	case r.Header.Get("X-Forwarded-For") != "":
		raw = r.Header.Get("X-Forwarded-For")
	case r.Header.Get("X-Client-IP") != "":
		raw = r.Header.Get("X-Client-IP")
	case r.RemoteAddr != "":
		raw = r.RemoteAddr
	}

	var host string
	host, _, err := net.SplitHostPort(raw)
	if err != nil {
		host = raw
	}

	ip := net.ParseIP(host)
	if ip == nil {
		return "-"
	}

	return ip.String()
}
