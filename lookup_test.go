package clientip_test

import (
	"net/http"
	"testing"

	"github.com/dmitrymomot/clientip"
)

func TestLookupFromRequest(t *testing.T) {
	type args struct {
		r       *http.Request
		headers []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Default headers",
			args: args{
				r: func() *http.Request {
					r, _ := http.NewRequest("GET", "http://example.com", nil)
					r.Header.Set("X-Forwarded-For", "48.135.12.111")
					r.Header.Set("X-Real-IP", "48.135.12.111")
					r.RemoteAddr = "230.173.15.154:1234"
					return r
				}(),
			},
			want: "48.135.12.111",
		},
		{
			name: "Custom headers",
			args: args{
				r: func() *http.Request {
					r, _ := http.NewRequest("GET", "http://example.com", nil)
					r.Header.Set("X-Forwarded-For", "48.135.12.111")
					r.Header.Set("X-Real-IP", "48.135.12.111")
					r.Header.Set("Custom-IP", "212.207.103.215")
					r.RemoteAddr = "230.173.15.154:1234"
					return r
				}(),
				headers: []string{"Custom-IP"},
			},
			want: "212.207.103.215",
		},
		{
			name: "X-Forwarded-For",
			args: args{
				r: func() *http.Request {
					r, _ := http.NewRequest("GET", "http://example.com", nil)
					r.Header.Set("X-Forwarded-For", "48.135.12.111")
					r.RemoteAddr = "230.173.15.154:1234"
					return r
				}(),
			},
			want: "48.135.12.111",
		},
		{
			name: "X-Forwarded-For: multiple addresses",
			args: args{
				r: func() *http.Request {
					r, _ := http.NewRequest("GET", "http://example.com", nil)
					r.Header.Set("X-Forwarded-For", "48.135.12.111, 181.95.251.176")
					r.RemoteAddr = "230.173.15.154:1234"
					return r
				}(),
			},
			want: "48.135.12.111",
		},
		{
			name: "default remote address",
			args: args{
				r: func() *http.Request {
					r, _ := http.NewRequest("GET", "http://example.com", nil)
					r.RemoteAddr = "230.173.15.154:1234"
					return r
				}(),
			},
			want: "230.173.15.154",
		},
		{
			name: "X-Forwarded-For: ipv6 address",
			args: args{
				r: func() *http.Request {
					r, _ := http.NewRequest("GET", "http://example.com", nil)
					r.Header.Set("X-Forwarded-For", "42e7:9f02:2ced:9303:2691:cd2e:7f9d:8ae3")
					return r
				}(),
			},
			want: "42e7:9f02:2ced:9303::", // /64 prefix.
		},
		{
			name: "invalid ip address: parse error",
			args: args{
				r: func() *http.Request {
					r, _ := http.NewRequest("GET", "http://example.com", nil)
					r.Header.Set("X-Forwarded-For", "123:456789")
					return r
				}(),
			},
			want: "",
		},
		{
			name: "not an ip address at all: no ip address found in the request",
			args: args{
				r: func() *http.Request {
					r, _ := http.NewRequest("GET", "http://example.com", nil)
					r.Header.Set("X-Forwarded-For", "123456789")
					return r
				}(),
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := clientip.LookupFromRequest(tt.args.r, tt.args.headers...); got != tt.want {
				t.Errorf("LookupFromRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
