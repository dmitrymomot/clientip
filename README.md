# Real IP Lookup Module

[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/dmitrymomot/clientip)](https://github.com/dmitrymomot/clientip)
[![Go Reference](https://pkg.go.dev/badge/github.com/dmitrymomot/clientip.svg)](https://pkg.go.dev/github.com/dmitrymomot/clientip)
[![License](https://img.shields.io/github/license/dmitrymomot/clientip)](https://github.com/dmitrymomot/clientip/blob/main/LICENSE)

[![Tests](https://github.com/dmitrymomot/clientip/actions/workflows/tests.yml/badge.svg)](https://github.com/dmitrymomot/clientip/actions/workflows/tests.yml)
[![CodeQL Analysis](https://github.com/dmitrymomot/clientip/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/dmitrymomot/clientip/actions/workflows/codeql-analysis.yml)
[![GolangCI Lint](https://github.com/dmitrymomot/clientip/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/dmitrymomot/clientip/actions/workflows/golangci-lint.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/dmitrymomot/clientip)](https://goreportcard.com/report/github.com/dmitrymomot/clientip)

A lightweight, production-ready Go module for accurately determining client IP addresses in HTTP applications. Particularly useful for applications behind proxies, load balancers, or CDNs like Cloudflare, DigitalOcean App Platform, and Fastly.

## Features

-   **Accurate IP Detection**: Intelligently extracts client IPs from various standard and vendor-specific headers
-   **IPv4 & IPv6 Support**: Full support for both IPv4 and IPv6 addresses with proper canonicalization
-   **Flexible Header Configuration**: Use default headers or specify custom ones for your specific setup
-   **Multiple Integration Options**:
    -   Direct IP lookup function for manual integration
    -   Drop-in middleware for automatic IP detection
    -   Context-based IP storage for thread-safe access
-   **Production-Ready**: Used in production environments with comprehensive test coverage

## Installation

To install the module, use the following command:

```bash
go get github.com/dmitrymomot/clientip
```

## Usage

### Direct IP Lookup

```go
func handler(w http.ResponseWriter, r *http.Request) {
    // Use default headers
    ip := clientip.LookupFromRequest(r)

    // Or specify custom headers
    ip = clientip.LookupFromRequest(r, "X-Custom-IP", "X-Real-IP")

    fmt.Fprintf(w, "Your IP: %s", ip)
}
```

### Middleware Integration

```go
func main() {
    mux := http.NewServeMux()

    // Basic middleware usage with default headers
    handler := clientip.Middleware()(mux)

    // Or with custom headers
    handler = clientip.Middleware("X-Custom-IP", "CF-Connecting-IP")(mux)

    http.ListenAndServe(":8080", handler)
}
```

### Context-Based IP Access

```go
func main() {
    mux := http.NewServeMux()

    // Chain both middlewares
    handler := clientip.Middleware()(mux)
    handler = clientip.IpToContext(handler) // Store IP in context, it must be called after Middleware

    http.ListenAndServe(":8080", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
    // Get IP from context anywhere in your handler chain
    if ip := clientip.GetIPAddress(r.Context()); ip != "" {
        fmt.Fprintf(w, "Your IP from context: %s", ip)
    }
}
```

### Supported Headers

The module checks the following headers by default (in order):

-   `DO_Connecting-IP`, `DO-Connecting-IP` (DigitalOcean)
-   `True-Client-IP`
-   `X-Real-IP`
-   `CF-Connecting-IP` (Cloudflare)
-   `Fastly-Client-IP`
-   `X-Cluster-Client-IP`
-   `X-Client-IP`
-   `X-Forwarded-For` (first IP in the chain)

You can override these by providing your own headers to the functions.

## Customization

Both `LookupFromRequest` and `Middleware` functions accept an optional list of headers to consider when looking for the client IP address. By default, a predefined set of common headers used for forwarding client IPs in proxy setups are checked.

## Best Practices

-   **Security**: Always validate and sanitize IP addresses before using them in security-critical contexts
-   **Header Order**: Configure headers in order of trust (most trusted first)
-   **Performance**: Use the middleware approach for consistent IP handling across your application
-   **Context Usage**: Prefer context-based IP access when working with complex handler chains

## Troubleshooting

### Common Issues

1. **Empty IP Address**

    - Check if your proxy/load balancer is properly configured to forward IP headers
    - Verify the header names match your infrastructure setup

2. **Incorrect IP Address**

    - Ensure headers are being set in the correct order of precedence
    - Check if intermediate proxies are modifying the headers

3. **IPv6 Handling**
    - The module automatically canonicalizes IPv6 addresses to their /64 prefix
    - If you need the full address, use the raw header value instead

For more help, please open an issue on GitHub.

## Contributing

If you wish to contribute to this project, please fork the repository and submit a pull request.

## License

This project is licensed under the [Apache 2.0](LICENSE) - see the `LICENSE` file for details.
