// Package cors provides functionality for adding CORS response headers into a handler chain.
package cors

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// AllowList allows the user to craft a custom hostname matcher.
// An AllowList returns true if the requested origin should include valid CORS headers
// based on an incoming hostname.
type AllowList func(host string) bool

// HostByRE matches the hostname to a given regular expression
func HostByRE(allowListRE *regexp.Regexp) AllowList {
	return func(host string) bool {
		return allowListRE.MatchString(host)
	}
}

// HostByAllowList matches the hostname to a given map of allowed hosts.
func HostByAllowList(allowList map[string]bool) AllowList {
	return func(host string) bool {
		return allowList[host]
	}
}

// OriginHost extracts the origin and origin hostname from a request.
func OriginHost(r *http.Request) (string, string, error) {
	origin := r.Header.Get("origin")
	url, err := url.Parse(origin)
	if err != nil {
		return "", "", err
	}

	host := strings.TrimRight(url.Host, ".")
	host, _, err = net.SplitHostPort(host)
	if err != nil {
		host = strings.TrimRight(url.Host, ".")
	}

	return origin, host, nil
}

// Handler inserts Access-Control-Allow-Origin headers into the request response.
// allowListRE matches hostnames against the given regexp, and provides CORS headers where allowed.
func Handler(allowed AllowList) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			origin, host, err := OriginHost(r)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			if allowed(host) {
				header := w.Header()
				if header.Get("Access-Control-Allow-Credentials") == "" {
					header.Set("Access-Control-Allow-Credentials", "true")
				}
				if header.Get("Access-Control-Allow-Headers") == "" {
					header.Set("Access-Control-Allow-Headers", "Authorization, Cache-Control, Content-Type")
				}
				if header.Get("Access-Control-Allow-Methods") == "" {
					header.Set("Access-Control-Allow-Methods", r.Method)
				}
				if header.Get("Access-Control-Allow-Origin") == "" {
					header.Set("Access-Control-Allow-Origin", origin)
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

// CreateOptionsHandler creates a generic options handler that takes a list of allowed methods.
// allowListRE matches hostnames against the given regexp, and provides CORS headers where allowed.
func CreateOptionsHandler(allowed AllowList, methods ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		methodString := strings.Join(methods, ", ")

		header := w.Header()
		header.Set("Access-Control-Allow-Headers", "Authorization, Cache-Control, Content-Type")
		header.Set("Access-Control-Allow-Methods", methodString)
		header.Set("Access-Control-Allow-Credentials", "true")

		origin, host, err := OriginHost(r)
		if err != nil {
			log.Println(err.Error())
			return
		}

		if allowed(host) {
			header.Set("Access-Control-Allow-Origin", origin)
		}

		age := int(time.Hour * 1 / time.Second)
		w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d, public", age))
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusNoContent)
	}
}
