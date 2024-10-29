package clientip

import "context"

// contextKey is a type used for keys in the context.
// These keys are used to store and retrieve values from the context.
// The keys are unique and are used to avoid conflicts with other keys.
type contextKey struct{ name string }

// requestIPKey is the key used to store the client's IP address in the request context.
var requestIPKey = contextKey{"request_ip"}

// SetIPAddress sets the client's IP address in the request context.
// This IP address can be used in the next handler.
func SetIPAddress(ctx context.Context, ip string) context.Context {
	return context.WithValue(ctx, requestIPKey, ip)
}

// GetIPAddress gets the client's IP address from the request context.
// If the IP address is not found in the context, it returns an empty string.
func GetIPAddress(ctx context.Context) string {
	if ip, ok := ctx.Value(requestIPKey).(string); ok {
		return ip
	}
	return ""
}
