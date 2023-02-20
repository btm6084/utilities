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

	// Trim the final one out of the chain, as that is the thing talking to us from the internet.
	candidate := strings.Trim(pieces[len(pieces)-1], `'" `)

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
