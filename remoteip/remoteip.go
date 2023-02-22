package remoteip

import (
	"net"
	"net/http"
	"strings"
)

var (
	ExcludeIPS = map[string]bool{}
)

type CIDRRange []*net.IPNet

// NewCIDRRange converts a slice of IP address strings into an IPList of IPNets (CIDRs)
func NewCIDRRange(in []string) (CIDRRange, error) {
	out := make(CIDRRange, len(in))
	for k, v := range in {
		_, network, err := net.ParseCIDR(v)
		if err != nil {
			return nil, err
		}

		out[k] = network
	}
	return out, nil
}

var (
	privateIPs CIDRRange
	loopIPs    CIDRRange
)

func init() {
	loopIPs, _ = NewCIDRRange([]string{
		"127.0.0.0/8",         //IPv4 loopback
		"0:0:0:0:0:0:0:1/128", //IPv6 loopback
	})
	privateIPs, _ = NewCIDRRange([]string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
	})
}

func isLoopbackIP(ip net.IP) bool {
	for _, cidr := range loopIPs {
		if cidr.Contains(ip) {
			return true
		}
	}
	return false
}

func isPrivateIP(ip net.IP) bool {
	for _, cidr := range privateIPs {
		if cidr.Contains(ip) {
			return true
		}
	}
	return false
}

// Get returns the remote IP address from the request.
func Get(r *http.Request) string {
	var ips []string
	if r.Header.Get("X-Forwarded-For") != "" {
		ips = append(ips, strings.Split(r.Header.Get("X-Forwarded-For"), ",")...)
	}
	ips = append([]string{r.RemoteAddr}, ips...)

	for i := 0; i < len(ips); i++ {
		var host string
		host, _, err := net.SplitHostPort(ips[i])
		if err != nil {
			if !strings.Contains(err.Error(), "missing port in address") {
				continue
			}
			host = strings.Trim(ips[i], `"' ,`)
		}

		ip := net.ParseIP(host)
		if ip == nil {
			continue
		}

		if isPrivateIP(ip) || isLoopbackIP(ip) {
			continue
		}

		return ip.String()
	}

	var host string
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		if !strings.Contains(err.Error(), "missing port in address") {
			return "-"
		}
		host = strings.Trim(r.RemoteAddr, `"' ,`)
	}

	ip := net.ParseIP(host)
	if ip == nil {
		return "-"
	}

	return ip.String()
}

// Get returns the remote IP address from the request.
func RemoteIPs(r *http.Request) map[string]string {
	candidates := make(map[string]string)

	if r.Header.Get("X-Real-IP") != "" {
		candidates["X-Real-IP"] = r.Header.Get("X-Real-IP")
	}
	if r.Header.Get("X-Forwarded-For") != "" {
		candidates["X-Forwarded-For"] = r.Header.Get("X-Forwarded-For")
	}
	if r.Header.Get("X-Client-IP") != "" {
		candidates["X-Client-IP"] = r.Header.Get("X-Client-IP")
	}
	if r.RemoteAddr != "" {
		candidates["rmtAddr"] = r.RemoteAddr
	}

	return candidates
}
