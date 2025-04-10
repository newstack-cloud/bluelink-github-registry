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

type GetPluginVersionstHandlerTestSuite struct {
	suite.Suite
	server *httptest.Server
}

func (s *GetPluginVersionstHandlerTestSuite) SetupTest() {
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

func (s *GetPluginVersionstHandlerTestSuite) TearDownTest() {
	s.server.Close()
}

func (s *GetPluginVersionstHandlerTestSuite) TestGetPluginVersions() {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/plugins/two-hundred/aws/versions", s.server.URL),
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

	versions := &types.PluginVersions{}
	err = json.Unmarshal(respBytes, versions)
	s.Require().NoError(err)

	s.Require().Equal(
		expectedVersions,
		versions,
	)
}

func TestGetPluginVersionstHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(GetPluginVersionstHandlerTestSuite))
}
