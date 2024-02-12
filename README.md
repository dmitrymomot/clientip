# Real IP Lookup Module

[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/dmitrymomot/clientip)](https://github.com/dmitrymomot/clientip)
[![Go Reference](https://pkg.go.dev/badge/github.com/dmitrymomot/clientip.svg)](https://pkg.go.dev/github.com/dmitrymomot/clientip)
[![License](https://img.shields.io/github/license/dmitrymomot/clientip)](https://github.com/dmitrymomot/clientip/blob/main/LICENSE)

[![Tests](https://github.com/dmitrymomot/clientip/actions/workflows/tests.yml/badge.svg)](https://github.com/dmitrymomot/clientip/actions/workflows/tests.yml)
[![CodeQL Analysis](https://github.com/dmitrymomot/clientip/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/dmitrymomot/clientip/actions/workflows/codeql-analysis.yml)
[![GolangCI Lint](https://github.com/dmitrymomot/clientip/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/dmitrymomot/clientip/actions/workflows/golangci-lint.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/dmitrymomot/clientip)](https://goreportcard.com/report/github.com/dmitrymomot/clientip)

This Go module provides functionality to accurately determine a client's IP address from HTTP requests, considering various headers that might contain the real IP address, especially when your application is behind a proxy or load balancer.

## Features

- **LookupFromRequest**: Retrieve the client IP address directly from an HTTP request with consideration for custom headers.
- **Middleware**: An easy-to-integrate middleware function for HTTP servers that automatically sets the `RemoteAddr` field of the request to the client's real IP address.

## Installation

To install the module, use the following command:

```bash
go get github.com/dmitrymomot/clientip
```

## Usage

### LookupFromRequest Function

`LookupFromRequest` retrieves the client IP address from the provided HTTP request. It checks a list of headers for the IP address and uses a default set if none are specified.

```go
import "github.com/dmitrymomot/clientip"

func handler(w http.ResponseWriter, r *http.Request) {
    clientIP := clientip.LookupFromRequest(r)
    fmt.Fprintf(w, "Client IP: %s", clientIP)
}
```

With custom headers:

```go
import "github.com/dmitrymomot/clientip"

func handler(w http.ResponseWriter, r *http.Request) {
    clientIP := clientip.LookupFromRequest(r, "X-Custom-Real-IP")
    fmt.Fprintf(w, "Client IP: %s", clientIP)
}
```

### Middleware

The `Middleware` function returns a middleware handler that modifies the request's `RemoteAddr` based on the client IP address found in the specified headers.

```go
import (
    "net/http"
    "github.com/dmitrymomot/clientip"
)

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", handler)

    // Init the middleware with custom header. If not specified, the default headers are used.
    mdw := clientip.Middleware("X-Custom-Real-IP")

    http.ListenAndServe(":8080", mdw(mux))
}

func handler(w http.ResponseWriter, r *http.Request) {
    // The RemoteAddr now contains the real client IP
    fmt.Fprintf(w, "Client IP: %s", r.RemoteAddr)
}
```

## Customization

Both `LookupFromRequest` and `Middleware` functions accept an optional list of headers to consider when looking for the client IP address. By default, a predefined set of common headers used for forwarding client IPs in proxy setups are checked.

## Contributing

If you wish to contribute to this project, please fork the repository and submit a pull request.

## License

This project is licensed under the [Apache 2.0](LICENSE) - see the `LICENSE` file for details.