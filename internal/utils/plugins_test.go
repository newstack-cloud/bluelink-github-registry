package utils

import (
	"context"
	"strings"
	"testing"

	"github.com/google/go-github/v70/github"
	"github.com/newstack-cloud/bluelink-github-registry/internal/testutils"
	"github.com/newstack-cloud/bluelink-github-registry/internal/types"
	"github.com/stretchr/testify/suite"
)

type PluginUtilsTestSuite struct {
	suite.Suite
}

func (s *PluginUtilsTestSuite) Test_extracts_plugin_versions_from_the_provided_releases() {
	pluginVersions, err := ExtractPluginVersions(
		context.Background(),
		"bluelink-provider-example",
		inputReleases1(),
		&testutils.StubHTTPClient{
			Contents: registryInfoContents(),
		},
		"test-token",
	)
	s.Require().NoError(err)
	s.Assert().Equal(
		expectedPluginVersions(),
		pluginVersions,
	)
}

func (s *PluginUtilsTestSuite) Test_finds_repository_for_provided_plugin() {
	pluginRepo := FindPluginRepo(
		reposToSearch(),
		"newstack-cloud",
		"example",
	)
	s.Assert().NotNil(pluginRepo)
	s.Assert().Equal(
		&github.Repository{
			Name: github.Ptr("bluelink-provider-example"),
			Owner: &github.User{
				Login: github.Ptr("newstack-cloud"),
			},
			Private: github.Ptr(true),
		},
		pluginRepo,
	)
}

func (s *PluginUtilsTestSuite) Test_extracts_plugin_package_info_for_the_provided_release() {
	signingKeys, err := testutils.GetSigningKeysFromEnv()
	s.Require().NoError(err)

	versionPackage, err := ExtractPluginVersionPackage(
		context.Background(),
		&ExtractPluginVersionPackageParams{
			Repository:            "bluelink-provider-example",
			Release:               packageInfoRelease(),
			Version:               "1.0.1",
			OS:                    "linux",
			Arch:                  "amd64",
			SigningKeysSerialised: signingKeys.Input,
		},
		&testutils.StubHTTPClient{
			ContentsProvider: providePackageInfoRequestContent,
		},
		"test-token",
	)
	s.Require().NoError(err)
	s.Assert().Equal(
		expectedVersionPackage(signingKeys.Expected),
		versionPackage,
	)
}

func expectedVersionPackage(
	expectedSigningKeys *types.PublicGPGSigningKeys,
) *types.PluginVersionPackage {
	return &types.PluginVersionPackage{
		SupportedProtocols: []string{"1.2", "2.0"},
		OS:                 "linux",
		Arch:               "amd64",
		Filename:           "bluelink-provider-example_1.0.1_linux_amd64.zip",
		// See the packageInfoRelease function for the URL in the source github releases.
		DownloadURL:         *testutils.GithubAssetURL(6),
		SHASumsURL:          packageInfoRegistrySHA256SumsURL(),
		SHASumsSignatureURL: *testutils.GithubAssetURL(8),
		SHASum:              "c635e6201021832cc1f4cfe5345",
		SigningKeys:         expectedSigningKeys,
		Dependencies: map[string]string{
			"bluelink/aws": "^1.0.0",
		},
	}
}

func expectedPluginVersions() *types.PluginVersions {
	return &types.PluginVersions{
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
	}
}

func reposToSearch() []*github.Repository {
	return []*github.Repository{
		{
			Name: github.Ptr("bluelink-provider-example"),
			Owner: &github.User{
				Login: github.Ptr("newstack-cloud"),
			},
			Private: github.Ptr(true),
		},
		{
			Name: github.Ptr("bluelink-transformer-example"),
			Owner: &github.User{
				Login: github.Ptr("newstack-cloud"),
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
					Name: github.Ptr("bluelink-provider-example_1.0.0_darwin_amd64.zip"),
					URL:  testutils.GithubAssetURL(1),
				},
				{
					Name: github.Ptr("bluelink-provider-example_1.0.0_linux_amd64.zip"),
					URL:  testutils.GithubAssetURL(2),
				},
				{
					Name: github.Ptr("bluelink-provider-example_1.0.0_windows_amd64.zip"),
					URL:  testutils.GithubAssetURL(3),
				},
				{
					Name: github.Ptr("bluelink-provider-example_1.0.0_registry_info.json"),
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
					Name: github.Ptr("bluelink-provider-example_1.0.1_darwin_amd64.zip"),
					URL:  testutils.GithubAssetURL(5),
				},
				{
					Name: github.Ptr("bluelink-provider-example_1.0.1_linux_amd64.zip"),
					URL:  testutils.GithubAssetURL(6),
				},
				{
					Name: github.Ptr("bluelink-provider-example_1.0.1_windows_amd64.zip"),
					URL:  testutils.GithubAssetURL(7),
				},
				{
					Name: github.Ptr("bluelink-provider-example_1.0.1_registry_info.json"),
					URL:  testutils.GithubAssetURL(8),
				},
			},
		},
	}
}

func registryInfoContents() []byte {
	return []byte(`
	{
		"supportedProtocols": ["1.2", "2.0"],
		"dependencies": {
			"bluelink/aws": "^1.0.0"
		}
	}
	`)
}

func packageInfoRelease() *github.RepositoryRelease {
	return &github.RepositoryRelease{
		TagName: github.Ptr("v1.0.1"),
		Assets: []*github.ReleaseAsset{
			{
				Name: github.Ptr("bluelink-provider-example_1.0.1_darwin_amd64.zip"),
				URL:  testutils.GithubAssetURL(5),
			},
			{
				Name: github.Ptr("bluelink-provider-example_1.0.1_linux_amd64.zip"),
				URL:  testutils.GithubAssetURL(6),
			},
			{
				Name: github.Ptr("bluelink-provider-example_1.0.1_windows_amd64.zip"),
				URL:  testutils.GithubAssetURL(7),
			},
			// The following files get their own separate URL to allow the contents retriever
			// to easily identify them when fetching the file contents.
			// These are the only two files that are downloaded in the process of preparing
			// version package info.
			{
				Name: github.Ptr("bluelink-provider-example_1.0.1_registry_info.json"),
				URL:  github.Ptr(packageInfoRegistryInfoURL()),
			},
			{
				Name: github.Ptr("bluelink-provider-example_1.0.1_SHA256SUMS"),
				URL:  github.Ptr(packageInfoRegistrySHA256SumsURL()),
			},
			{
				Name: github.Ptr("bluelink-provider-example_1.0.1_SHA256SUMS.sig"),
				URL:  testutils.GithubAssetURL(8),
			},
		},
	}
}

func packageInfoRegistryInfoURL() string {
	return "https://artifacts.example.com/bluelink-provider-example/1.0.1/bluelink-provider-example_1.0.1_registry_info.json"
}

func packageInfoRegistrySHA256SumsURL() string {
	return "https://artifacts.example.com/bluelink-provider-example/1.0.1/bluelink-provider-example_1.0.1_SHA256SUMS"
}

func providePackageInfoRequestContent(url string) ([]byte, error) {
	if strings.HasPrefix(url, packageInfoRegistryInfoURL()) {
		return registryInfoContents(), nil
	}

	if strings.HasPrefix(url, packageInfoRegistrySHA256SumsURL()) {
		return packageSHASumContents(), nil
	}

	return []byte{}, nil
}

func packageSHASumContents() []byte {
	return []byte(`
		c3e51ec2a5857d4e2e48af02de97  bluelink-provider-example_1.0.1_darwin_amd64.zip
		ed370cc761421bfd60479d4f6214  bluelink-provider-example_1.0.1_darwin_arm64.zip
		03f5694b5a0fec5b328365bb294  bluelink-provider-example_1.0.1_docs.json
		34623f6a541be48b5314e6e2ebb  bluelink-provider-example_1.0.1_linux_386.zip
		c635e6201021832cc1f4cfe5345  bluelink-provider-example_1.0.1_linux_amd64.zip
		4cfc841b4582ad748133dba0fce  bluelink-provider-example_1.0.1_linux_arm.zip
		14a971e72106337503baa26cfe4  bluelink-provider-example_1.0.1_linux_arm64.zip
		02a95af4369f9f0edc1d4ef6deb  bluelink-provider-example_1.0.1_registry_info.json
	`)
}

func TestPluginUtilsTestSuite(t *testing.T) {
	suite.Run(t, new(PluginUtilsTestSuite))
}
