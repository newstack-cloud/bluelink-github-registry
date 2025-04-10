package utils

import (
	"context"
	"testing"

	"github.com/google/go-github/v70/github"
	"github.com/stretchr/testify/suite"
	"github.com/two-hundred/celerity-github-registry/internal/testutils"
	"github.com/two-hundred/celerity-github-registry/internal/types"
)

type PluginUtilsTestSuite struct {
	suite.Suite
}

func (s *PluginUtilsTestSuite) Test_extracts_plugin_versions_from_the_provided_releases() {
	pluginVersions, err := ExtractPluginVersions(
		context.Background(),
		"celerity-provider-example",
		inputReleases1(),
		&testutils.StubHTTPClient{
			Contents: registryInfoContents(),
		},
		"test-token",
	)
	s.Require().NoError(err)
	s.Assert().Equal(
		&types.PluginVersions{
			Versions: []*types.PluginVersion{
				{
					Version:            "1.0.0",
					SupportedProtocols: []string{"1.2", "2.0"},
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
					SupportedProtocols: []string{"1.2", "2.0"},
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
		pluginVersions,
	)
}

func (s *PluginUtilsTestSuite) Test_finds_repository_for_provided_plugin() {
	pluginRepo := FindPluginRepo(
		reposToSearch(),
		"two-hundred",
		"example",
	)
	s.Assert().NotNil(pluginRepo)
	s.Assert().Equal(
		&github.Repository{
			Name: github.Ptr("celerity-provider-example"),
			Owner: &github.User{
				Login: github.Ptr("two-hundred"),
			},
			Private: github.Ptr(true),
		},
		pluginRepo,
	)
}

func reposToSearch() []*github.Repository {
	return []*github.Repository{
		{
			Name: github.Ptr("celerity-provider-example"),
			Owner: &github.User{
				Login: github.Ptr("two-hundred"),
			},
			Private: github.Ptr(true),
		},
		{
			Name: github.Ptr("celerity-transformer-example"),
			Owner: &github.User{
				Login: github.Ptr("two-hundred"),
			},
			Private: github.Ptr(true),
		},
	}
}

func inputReleases1() []*github.RepositoryRelease {
	return []*github.RepositoryRelease{
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
			// This release should be ignored as the tag
			// is not of the form vX.Y.Z.
			TagName: github.Ptr("some-other-tag-1.2"),
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
	}
}

func registryInfoContents() []byte {
	return []byte(`
	{
		"supportedProtocols": ["1.2", "2.0"]
	}
	`)
}

func TestPluginUtilsTestSuite(t *testing.T) {
	suite.Run(t, new(PluginUtilsTestSuite))
}
