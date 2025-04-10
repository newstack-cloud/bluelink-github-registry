package testutils

import (
	"bytes"
	"io"
	"net/http"
)

// StubHTTPClient provides a stub client
// that returns the provided contents
// when the Do method is called.
// It implements the httputils.Client interface.
type StubHTTPClient struct {
	Contents []byte
}

// Do returns a stub response with the configured contents.
func (c *StubHTTPClient) Do(req *http.Request) (*http.Response, error) {
	resp := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewReader(c.Contents)),
	}
	return resp, nil
}
