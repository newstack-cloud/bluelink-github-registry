package registry

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
	"github.com/two-hundred/celerity-github-registry/internal/core"
	"go.uber.org/zap"
)

type GetManifestHandlerTestSuite struct {
	suite.Suite
	server *httptest.Server
}

func (s *GetManifestHandlerTestSuite) SetupTest() {
	router := mux.NewRouter()

	getDeps := func(
		config *core.Config,
		logger *zap.Logger,
	) *registryDependencies {
		// Empty dependencies are fine for the get manifest handler
		// as it doesn't make use of the dependencies to generate
		// the manifest.
		return &registryDependencies{}
	}

	_, _, err := Setup(router, getDeps)
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
			ProviderV1: &PluginTypeManifestInfo{
				Endpoint:                  "http://gh-registry.celerity.local/plugins",
				DownloadAcceptContentType: DownloadContentType,
			},
			TransformerV1: &PluginTypeManifestInfo{
				Endpoint:                  "http://gh-registry.celerity.local/plugins",
				DownloadAcceptContentType: DownloadContentType,
			},
			AuthV1: &AuthManifestInfo{
				APIKeyHeader: "celerity-gh-registry-token",
				DownloadAuth: "bearer",
			},
		},
		manifest,
	)
}

func TestGetManifestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(GetManifestHandlerTestSuite))
}
