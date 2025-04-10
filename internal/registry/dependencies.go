package registry

import (
	"github.com/two-hundred/celerity-github-registry/internal/core"
	"github.com/two-hundred/celerity-github-registry/internal/httputils"
	"github.com/two-hundred/celerity-github-registry/internal/plugins"
	"github.com/two-hundred/celerity-github-registry/internal/repos"
	"go.uber.org/zap"
)

// GetDependencies retrieves the dependencies for the registry application
// endpoint handlers.
func GetDependencies(
	config *core.Config,
	logger *zap.Logger,
) *registryDependencies {
	httpClient := httputils.NewNativeHTTPClient(
		httputils.WithNativeHTTPClientTimeout(config.HTTPClientTimeout),
	)
	repoService := repos.NewGitHubService()
	pluginService := plugins.NewDefaultService(
		repoService,
		httpClient,
		config,
		logger,
	)
	return &registryDependencies{
		pluginService: pluginService,
	}
}
