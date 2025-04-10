package httputils

import (
	"net/http"
	"time"
)

const (
	// DefaultHTTPTimeout is the default timeout for HTTP requests.
	DefaultHTTPTimeout = 30
)

// Client provides an abstraction for any type of http client
// that can send a http request.
type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

// NativeHTTPClientOptions is a function that configures a http.Client.
type NativeHTTPClientOptions func(*http.Client)

// WithNativeHTTPClientTimeout configures a http.Client instance
// with a timeout, where timeout is in seconds.
func WithNativeHTTPClientTimeout(timeout int) NativeHTTPClientOptions {
	return func(c *http.Client) {
		c.Timeout = time.Second * time.Duration(timeout)
	}
}

// NewNativeHTTPClient creates a new instance of a HTTP client
// configured with a timeout.
// This is to be used with packages that only have interoperability
// with the built-in http.Client.
func NewNativeHTTPClient(opts ...NativeHTTPClientOptions) *http.Client {
	client := &http.Client{
		Timeout:   DefaultHTTPTimeout * time.Second,
		Transport: http.DefaultTransport,
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}
