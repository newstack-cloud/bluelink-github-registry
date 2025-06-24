package registry

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/newstack-cloud/bluelink-github-registry/internal/core"
	"github.com/newstack-cloud/bluelink-github-registry/internal/types"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type GetPluginVersionsHandlerTestSuite struct {
	suite.Suite
	server *httptest.Server
}

func (s *GetPluginVersionsHandlerTestSuite) SetupTest() {
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

func (s *GetPluginVersionsHandlerTestSuite) TearDownTest() {
	s.server.Close()
}

func (s *GetPluginVersionsHandlerTestSuite) Test_get_plugin_versions() {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/plugins/newstack-cloud/aws/versions", s.server.URL),
		nil,
	)
	s.Require().NoError(err)
	req.Header.Set("bluelink-gh-registry-token", "test-token")

	resp, err := http.DefaultClient.Do(req)
	s.Require().NoError(err)
	s.Require().Equal(200, resp.StatusCode)
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	s.Require().NoError(err)

	versions := &types.PluginVersions{}
	err = json.Unmarshal(respBytes, versions)
	s.Require().NoError(err)

	s.Require().Equal(
		expectedVersions,
		versions,
	)
}

func (s *GetPluginVersionsHandlerTestSuite) Test_returns_401_response_for_missing_token() {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/plugins/newstack-cloud/aws/versions", s.server.URL),
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

func (s *GetPluginVersionsHandlerTestSuite) Test_returns_404_response_for_missing_plugin_repo() {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/plugins/newstack-cloud/azure/versions", s.server.URL),
		nil,
	)
	s.Require().NoError(err)
	req.Header.Set("bluelink-gh-registry-token", "test-token")

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

func (s *GetPluginVersionsHandlerTestSuite) Test_returns_401_response_for_an_inaccessible_org() {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/plugins/other-org/aws/versions", s.server.URL),
		nil,
	)
	s.Require().NoError(err)
	req.Header.Set("bluelink-gh-registry-token", "test-token")

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

func (s *GetPluginVersionsHandlerTestSuite) Test_returns_403_response_for_a_forbidden_plugin() {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/plugins/newstack-cloud/forbidden-plugin/versions", s.server.URL),
		nil,
	)
	s.Require().NoError(err)
	req.Header.Set("bluelink-gh-registry-token", "test-token")

	resp, err := http.DefaultClient.Do(req)
	s.Require().NoError(err)
	s.Require().Equal(403, resp.StatusCode)
	defer resp.Body.Close()
	respBytes, err := io.ReadAll(resp.Body)
	s.Require().NoError(err)
	s.Require().Equal(
		`{"message":"Forbidden"}`,
		string(respBytes),
	)
}

func TestGetPluginVersionsHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(GetPluginVersionsHandlerTestSuite))
}
