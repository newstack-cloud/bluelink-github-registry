package registry

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
	"github.com/two-hundred/celerity-github-registry/internal/core"
	"github.com/two-hundred/celerity-github-registry/internal/types"
	"go.uber.org/zap"
)

type GetPluginPackageHandlerTestSuite struct {
	suite.Suite
	server *httptest.Server
}

func (s *GetPluginPackageHandlerTestSuite) SetupTest() {
	router := mux.NewRouter()

	getDeps := func(
		config *core.Config,
		logger *zap.Logger,
	) *registryDependencies {
		// Empty dependencies are fine for the get manifest handler
		// as it doesn't make use of the dependencies to generate
		// the manifest.
		return &registryDependencies{
			pluginService: &stubPluginService{},
		}
	}

	_, _, err := Setup(router, getDeps)
	s.Require().NoError(err)

	server := httptest.NewServer(router)
	s.server = server
}

func (s *GetPluginPackageHandlerTestSuite) TearDownTest() {
	s.server.Close()
}

func (s *GetPluginPackageHandlerTestSuite) Test_get_plugin_versions() {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/plugins/two-hundred/aws/1.0.1/package/linux/amd64", s.server.URL),
		nil,
	)
	s.Require().NoError(err)
	req.Header.Set("celerity-gh-registry-token", "test-token")

	resp, err := http.DefaultClient.Do(req)
	s.Require().NoError(err)
	s.Require().Equal(200, resp.StatusCode)
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	s.Require().NoError(err)

	pkg := &types.PluginVersionPackage{}
	err = json.Unmarshal(respBytes, pkg)
	s.Require().NoError(err)

	s.Require().Equal(
		expectedVersionPackage,
		pkg,
	)
}

func (s *GetPluginPackageHandlerTestSuite) Test_returns_401_response_for_missing_token() {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/plugins/two-hundred/aws/1.0.1/package/linux/amd64", s.server.URL),
		nil,
	)
	s.Require().NoError(err)
	// No token set in the request header.

	resp, err := http.DefaultClient.Do(req)
	s.Require().NoError(err)
	s.Require().Equal(401, resp.StatusCode)
	defer resp.Body.Close()
	respBytes, err := io.ReadAll(resp.Body)
	s.Require().NoError(err)
	s.Require().Equal(
		`{"message":"Unauthorized"}`,
		string(respBytes),
	)
}

func (s *GetPluginPackageHandlerTestSuite) Test_returns_404_response_for_missing_plugin_repo() {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/plugins/two-hundred/azure/1.0.1/package/linux/amd64", s.server.URL),
		nil,
	)
	s.Require().NoError(err)
	req.Header.Set("celerity-gh-registry-token", "test-token")

	resp, err := http.DefaultClient.Do(req)
	s.Require().NoError(err)
	s.Require().Equal(404, resp.StatusCode)
	defer resp.Body.Close()
	respBytes, err := io.ReadAll(resp.Body)
	s.Require().NoError(err)
	s.Require().Equal(
		`{"message":"Plugin repository not found"}`,
		string(respBytes),
	)
}

func TestGetPluginPackageHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(GetPluginPackageHandlerTestSuite))
}
