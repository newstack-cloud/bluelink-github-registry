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
	Contents         []byte
	ContentsProvider func(url string) ([]byte, error)
}

// Do returns a stub response with the configured contents.
func (c *StubHTTPClient) Do(req *http.Request) (*http.Response, error) {
	var body io.ReadCloser
	if c.ContentsProvider != nil {
		contents, err := c.ContentsProvider(req.URL.String())
		if err != nil {
			return nil, err
		}
		body = io.NopCloser(bytes.NewReader(contents))
	} else {
		body = io.NopCloser(bytes.NewReader(c.Contents))
	}

	resp := &http.Response{
		StatusCode: 200,
		Body:       body,
	}
	return resp, nil
}
