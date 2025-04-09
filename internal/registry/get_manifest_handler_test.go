package registry

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
)

type GetManifestHandlerTestSuite struct {
	suite.Suite
	server *httptest.Server
}

func (s *GetManifestHandlerTestSuite) SetupTest() {
	router := mux.NewRouter()
	_, _, err := Setup(router)
	s.Require().NoError(err)

	server := httptest.NewServer(router)
	s.server = server
}

func (s *GetManifestHandlerTestSuite) TearDownTest() {
	s.server.Close()
}

func (s *GetManifestHandlerTestSuite) TestGetManifest() {
	resp, err := http.Get(s.server.URL + "/.well-known/celerity-services.json")
	s.Require().NoError(err)
	s.Require().Equal(200, resp.StatusCode)
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	s.Require().NoError(err)

	manifest := &Manifest{}
	err = json.Unmarshal(respBytes, manifest)
	s.Require().NoError(err)

	s.Require().Equal(
		&Manifest{
			ProviderV1:    "http://gh-registry.celerity.local/plugins",
			TransformerV1: "http://gh-registry.celerity.local/plugins",
			AuthV1: &AuthManifestInfo{
				APIKeyHeader: "celerity-gh-registry-token",
			},
		},
		manifest,
	)
}

func TestGetManifestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(GetManifestHandlerTestSuite))
}
