package clientip

import (
	"net/http"
)

// Middleware is a function that returns a middleware handler.
// It takes a variadic number of headers as input parameters.
// The returned middleware handler modifies the request's RemoteAddr
// based on the values of the specified headers.
// If any of the headers are found in the request, the RemoteAddr is updated
// with the corresponding header value.
// The modified request is then passed to the next handler in the chain.
func Middleware(headers ...string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if rip := LookupFromRequest(r, headers...); rip != "" {
				r.RemoteAddr = rip
			}
			h.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

// IpToContext is a middleware that sets the client's IP address in the request context.
// This IP address can be used in the next handler.
func IpToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get user ip address and set it to the request context
		// This ip address can be used in the next handler
		next.ServeHTTP(w, r.WithContext(SetIPAddress(r.Context(), r.RemoteAddr)))
	})
}
