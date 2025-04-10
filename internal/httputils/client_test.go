package httputils

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
)

type NativeHTTPClientTestSuite struct {
	suite.Suite
}

func (s *NativeHTTPClientTestSuite) Test_creates_a_usable_native_http_client() {
	client := NewNativeHTTPClient(
		WithNativeHTTPClientTimeout(30),
	)
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/plain")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Success"))
			},
		),
	)
	defer server.Close()

	req, err := http.NewRequest(
		http.MethodGet,
		server.URL,
		nil,
	)
	s.Require().NoError(err)

	resp, err := client.Do(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	s.Require().NoError(err)
	s.Require().Equal("Success", string(bodyBytes))
}

func TestNativeHTTPClientTestSuite(t *testing.T) {
	suite.Run(t, new(NativeHTTPClientTestSuite))
}
