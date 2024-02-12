package clientip

import (
	"net"
	"net/http"
	"strings"
)

// LookupFromRequest retrieves the client IP address from the provided HTTP request.
// It accepts an optional list of header names to check for the IP address.
// If no headers are specified, it uses a default set of headers.
// The function returns the client IP address as a string.
func LookupFromRequest(r *http.Request, headers ...string) string {
	return lookupFromRequest(r, headers...)
}

// IP header names
var ipHeaders = []string{
	"DO_Connecting-IP", // DigitalOcean
	"DO-Connecting-IP",
	"True-Client-IP",
	"X-Real-IP",
	"CF-Connecting-IP", // Cloudflare
	"Fastly-Client-IP",
	"X-Cluster-Client-IP",
	"X-Client-IP",
}

// getRealIP returns the real IP address.
func lookupFromRequest(r *http.Request, headers ...string) string {
	// Check the provided headers first
	for _, header := range headers {
		if ip := r.Header.Get(header); ip != "" {
			return canonicalizeIP(ip)
		}
	}

	// Check the default headers
	for _, header := range ipHeaders {
		if ip := r.Header.Get(header); ip != "" {
			return canonicalizeIP(ip)
		}
	}

	// Fallback to the remote address
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		if i := strings.Index(xff, ", "); i != -1 {
			return canonicalizeIP(xff[:i])
		}
		return canonicalizeIP(xff)
	}

	if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return canonicalizeIP(host)
	}

	return ""
}

// canonicalizeIP returns a form of ip suitable for comparison to other IPs.
// For IPv4 addresses, this is simply the whole string.
// For IPv6 addresses, this is the /64 prefix.
func canonicalizeIP(ip string) string {
	isIPv6 := false
	// This is how net.ParseIP decides if an address is IPv6
	// https://cs.opensource.google/go/go/+/refs/tags/go1.17.7:src/net/ip.go;l=704
	for i := 0; !isIPv6 && i < len(ip); i++ {
		switch ip[i] {
		case '.':
			// IPv4
			return ip
		case ':':
			// IPv6
			isIPv6 = true
		}
	}
	if !isIPv6 {
		return "" // Not an IP address at all
	}

	ipv6 := net.ParseIP(ip)
	if ipv6 == nil {
		return "" // Invalid IP
	}

	return ipv6.Mask(net.CIDRMask(64, 128)).String()
}
