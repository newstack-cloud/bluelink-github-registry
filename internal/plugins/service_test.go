package plugins

import (
	"context"
	"testing"

	"github.com/google/go-github/v70/github"
	"github.com/stretchr/testify/suite"
	"github.com/two-hundred/celerity-github-registry/internal/core"
	"github.com/two-hundred/celerity-github-registry/internal/testutils"
	"github.com/two-hundred/celerity-github-registry/internal/types"
	"go.uber.org/zap"
)

type DefaultServiceTestSuite struct {
	suite.Suite
	service Service
}

func (s *DefaultServiceTestSuite) SetupTest() {
	logger, err := zap.NewDevelopment()
	s.Require().NoError(err)

	s.service = NewDefaultService(
		testutils.NewStubRepoService(
			stubRepos(),
			stubRepoReleases(),
		),
		&testutils.StubHTTPClient{
			Contents: registryInfoContents(),
		},
		&core.Config{},
		logger,
	)
}

func (s *DefaultServiceTestSuite) TestListVersions() {
	versions, err := s.service.ListVersions(
		context.Background(),
		"two-hundred",
		"example",
		"test-token",
	)
	s.Require().NoError(err)
	s.Assert().Equal(
		&types.PluginVersions{
			Versions: []*types.PluginVersion{
				{
					Version:            "1.0.0",
					SupportedProtocols: []string{"1.4", "2.1"},
					SupportedPlatforms: []*types.PluginVersionPlatform{
						{
							OS:   "darwin",
							Arch: "amd64",
						},
						{
							OS:   "linux",
							Arch: "amd64",
						},
						{
							OS:   "windows",
							Arch: "amd64",
						},
					},
				},
				{
					Version:            "1.0.1",
					SupportedProtocols: []string{"1.4", "2.1"},
					SupportedPlatforms: []*types.PluginVersionPlatform{
						{
							OS:   "darwin",
							Arch: "amd64",
						},
						{
							OS:   "linux",
							Arch: "amd64",
						},
						{
							OS:   "windows",
							Arch: "amd64",
						},
					},
				},
			},
		},
		versions,
	)
}

func stubRepos() []*github.Repository {
	return []*github.Repository{
		{
			Name:        github.Ptr("celerity-provider-example"),
			FullName:    github.Ptr("two-hundred/celerity-provider-example"),
			Description: github.Ptr("A plugin for Celerity"),
			Private:     github.Ptr(true),
			Owner: &github.User{
				Login: github.Ptr("two-hundred"),
			},
		},
		{
			Name:        github.Ptr("celerity-transformer-exampleTransform"),
			FullName:    github.Ptr("two-hundred/celerity-transformer-exampleTransform"),
			Description: github.Ptr("A plugin for Celerity"),
			Private:     github.Ptr(true),
			Owner: &github.User{
				Login: github.Ptr("two-hundred"),
			},
		},
	}
}

func stubRepoReleases() map[string][]*github.RepositoryRelease {
	return map[string][]*github.RepositoryRelease{
		"celerity-provider-example": {
			{
				TagName: github.Ptr("v1.0.0"),
				Assets: []*github.ReleaseAsset{
					{
						Name: github.Ptr("celerity-provider-example_1.0.0_darwin_amd64.zip"),
						URL:  testutils.GithubAssetURL(1),
					},
					{
						Name: github.Ptr("celerity-provider-example_1.0.0_linux_amd64.zip"),
						URL:  testutils.GithubAssetURL(2),
					},
					{
						Name: github.Ptr("celerity-provider-example_1.0.0_windows_amd64.zip"),
						URL:  testutils.GithubAssetURL(3),
					},
					{
						Name: github.Ptr("celerity-provider-example_1.0.0_registry_info.json"),
						URL:  testutils.GithubAssetURL(4),
					},
				},
			},
			{
				TagName: github.Ptr("v1.0.1"),
				Assets: []*github.ReleaseAsset{
					{
						Name: github.Ptr("celerity-provider-example_1.0.1_darwin_amd64.zip"),
						URL:  testutils.GithubAssetURL(5),
					},
					{
						Name: github.Ptr("celerity-provider-example_1.0.1_linux_amd64.zip"),
						URL:  testutils.GithubAssetURL(6),
					},
					{
						Name: github.Ptr("celerity-provider-example_1.0.1_windows_amd64.zip"),
						URL:  testutils.GithubAssetURL(7),
					},
					{
						Name: github.Ptr("celerity-provider-example_1.0.1_registry_info.json"),
						URL:  testutils.GithubAssetURL(8),
					},
				},
			},
		},
		"celerity-transformer-exampleTransform": {
			{
				TagName: github.Ptr("v1.0.0"),
				Assets: []*github.ReleaseAsset{
					{
						Name: github.Ptr("celerity-transformer-exampleTransform_1.0.0_darwin_amd64.zip"),
						URL:  testutils.GithubAssetURL(9),
					},
					{
						Name: github.Ptr("celerity-transformer-exampleTransform_1.0.0_linux_amd64.zip"),
						URL:  testutils.GithubAssetURL(10),
					},
					{
						Name: github.Ptr("celerity-transformer-exampleTransform_1.0.0_windows_amd64.zip"),
						URL:  testutils.GithubAssetURL(11),
					},
					{
						Name: github.Ptr("celerity-transformer-exampleTransform_1.0.0_registry_info.json"),
						URL:  testutils.GithubAssetURL(12),
					},
				},
			},
			{
				TagName: github.Ptr("v1.1.0"),
				Assets: []*github.ReleaseAsset{
					{
						Name: github.Ptr("celerity-transformer-exampleTransform_1.1.0_darwin_amd64.zip"),
						URL:  testutils.GithubAssetURL(13),
					},
					{
						Name: github.Ptr("celerity-transformer-exampleTransform_1.1.0_linux_amd64.zip"),
						URL:  testutils.GithubAssetURL(14),
					},
					{
						Name: github.Ptr("celerity-transformer-exampleTransform_1.1.0_windows_amd64.zip"),
						URL:  testutils.GithubAssetURL(15),
					},
					{
						Name: github.Ptr("celerity-transformer-exampleTransform_1.1.0_registry_info.json"),
						URL:  testutils.GithubAssetURL(16),
					},
				},
			},
		},
	}
}

func registryInfoContents() []byte {
	return []byte(`
	{
		"supportedProtocols": ["1.4", "2.1"]
	}
	`)
}

func TestDefaultServiceTestSuite(t *testing.T) {
	suite.Run(t, new(DefaultServiceTestSuite))
}
