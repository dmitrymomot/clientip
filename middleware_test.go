package clientip_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dmitrymomot/clientip"
)

func TestMiddleware(t *testing.T) {
	tests := []struct {
		name              string
		headers           []string
		requestHeaders    map[string]string
		defaultRemoteAddr string
		expectedIP        string
	}{
		{
			name:           "Default headers",
			headers:        nil,
			requestHeaders: map[string]string{"X-Forwarded-For": "48.135.12.111", "X-Real-IP": "48.135.12.111"},
			expectedIP:     "48.135.12.111",
		},
		{
			name:           "Custom headers",
			headers:        []string{"Custom-IP"},
			requestHeaders: map[string]string{"X-Forwarded-For": "48.135.12.111", "X-Real-IP": "48.135.12.111", "Custom-IP": "212.207.103.215"},
			expectedIP:     "212.207.103.215",
		},
		{
			name:           "X-Forwarded-For",
			headers:        nil,
			requestHeaders: map[string]string{"X-Forwarded-For": "48.135.12.111"},
			expectedIP:     "48.135.12.111",
		},
		{
			name:           "X-Forwarded-For: multiple addresses",
			headers:        nil,
			requestHeaders: map[string]string{"X-Forwarded-For": "48.135.12.111, 181.95.251.176"},
			expectedIP:     "48.135.12.111",
		},
		{
			name:              "default remote address",
			headers:           nil,
			requestHeaders:    map[string]string{},
			defaultRemoteAddr: "249.212.158.68",
			expectedIP:        "249.212.158.68",
		},
		{
			name:           "X-Forwarded-For: ipv6 address",
			headers:        nil,
			requestHeaders: map[string]string{"X-Forwarded-For": "42e7:9f02:2ced:9303:2691:cd2e:7f9d:8ae3"},
			expectedIP:     "42e7:9f02:2ced:9303::",
		},
		{
			name:           "invalid ip address: parse error",
			headers:        nil,
			requestHeaders: map[string]string{"X-Forwarded-For": "123:456789"},
			expectedIP:     "",
		},
		{
			name:           "not an ip address at all: no ip address found in the request",
			headers:        nil,
			requestHeaders: map[string]string{"X-Forwarded-For": "123456789"},
			expectedIP:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new request with the specified headers
			req := httptest.NewRequest("GET", "http://example.com", nil)
			req.RemoteAddr = tt.defaultRemoteAddr
			for key, value := range tt.requestHeaders {
				req.Header.Set(key, value)
			}

			// Create a new response recorder
			res := httptest.NewRecorder()

			// Create a middleware handler with the specified headers
			handler := clientip.Middleware(tt.headers...)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Check if the remote address is set correctly
				if r.RemoteAddr != tt.expectedIP {
					t.Errorf("Unexpected remote address, got %s, want %s", r.RemoteAddr, tt.expectedIP)
				}
			}))

			// Serve the request using the middleware handler
			handler.ServeHTTP(res, req)
		})
	}
}
