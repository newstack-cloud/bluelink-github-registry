package plugins

import (
	"context"
	"strings"
	"testing"

	"github.com/google/go-github/v70/github"
	"github.com/newstack-cloud/bluelink-github-registry/internal/core"
	"github.com/newstack-cloud/bluelink-github-registry/internal/testutils"
	"github.com/newstack-cloud/bluelink-github-registry/internal/types"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type DefaultServiceTestSuite struct {
	suite.Suite
	service Service
}

func (s *DefaultServiceTestSuite) SetupTest() {
	logger, err := zap.NewDevelopment()
	s.Require().NoError(err)

	config, err := core.LoadConfigFromEnv()
	s.Require().NoError(err)

	s.service = NewDefaultService(
		testutils.NewStubRepoService(
			stubRepos(),
			stubRepoReleases(),
		),
		&testutils.StubHTTPClient{
			ContentsProvider: contentsProvider,
		},
		&config,
		logger,
	)
}

func (s *DefaultServiceTestSuite) TestListVersions() {
	versions, err := s.service.ListVersions(
		context.Background(),
		"newstack-cloud",
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

func (s *DefaultServiceTestSuite) TestGetPackageInfo() {
	signingKeysInfo, err := testutils.GetSigningKeysFromEnv()
	s.Require().NoError(err)

	packageInfo, err := s.service.GetPackageInfo(
		context.Background(),
		&PackageInfoParams{
			Organisation: "newstack-cloud",
			Plugin:       "example",
			Version:      "1.0.1",
			OS:           "linux",
			Arch:         "amd64",
		},
		"test-token",
	)
	s.Require().NoError(err)
	s.Assert().Equal(
		&types.PluginVersionPackage{
			SupportedProtocols: []string{"1.4", "2.1"},
			OS:                 "linux",
			Arch:               "amd64",
			Filename:           "bluelink-provider-example_1.0.1_linux_amd64.zip",
			// See the stubRepoReleases function for the URL in the source github releases.
			DownloadURL:         *testutils.GithubAssetURL(6),
			SHASumsURL:          packageInfoRegistrySHA256SumsURL(),
			SHASumsSignatureURL: *testutils.GithubAssetURL(8),
			SHASum:              "c635e6201021832cc1f4cfe5345",
			SigningKeys:         signingKeysInfo.Expected,
			Dependencies: map[string]string{
				"bluelink/aws": "^1.0.0",
			},
		},
		packageInfo,
	)
}

func stubRepos() []*github.Repository {
	return []*github.Repository{
		{
			Name:        github.Ptr("bluelink-provider-example"),
			FullName:    github.Ptr("newstack-cloud/bluelink-provider-example"),
			Description: github.Ptr("A plugin for Bluelink"),
			Private:     github.Ptr(true),
			Owner: &github.User{
				Login: github.Ptr("newstack-cloud"),
			},
		},
		{
			Name:        github.Ptr("bluelink-transformer-exampleTransform"),
			FullName:    github.Ptr("newstack-cloud/bluelink-transformer-exampleTransform"),
			Description: github.Ptr("A plugin for Bluelink"),
			Private:     github.Ptr(true),
			Owner: &github.User{
				Login: github.Ptr("newstack-cloud"),
			},
		},
	}
}

func stubRepoReleases() map[string][]*github.RepositoryRelease {
	return map[string][]*github.RepositoryRelease{
		"bluelink-provider-example": {
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
					// The following file gets its own separate URL to allow the contents retriever
					// to easily identify it to return the SHA256SUMS contents instead of the registry
					// info contents.
					{
						Name: github.Ptr("bluelink-provider-example_1.0.1_SHA256SUMS"),
						URL:  github.Ptr(packageInfoRegistrySHA256SumsURL()),
					},
					{
						Name: github.Ptr("bluelink-provider-example_1.0.1_SHA256SUMS.sig"),
						URL:  testutils.GithubAssetURL(8),
					},
				},
			},
		},
		"bluelink-transformer-exampleTransform": {
			{
				TagName: github.Ptr("v1.0.0"),
				Assets: []*github.ReleaseAsset{
					{
						Name: github.Ptr("bluelink-transformer-exampleTransform_1.0.0_darwin_amd64.zip"),
						URL:  testutils.GithubAssetURL(9),
					},
					{
						Name: github.Ptr("bluelink-transformer-exampleTransform_1.0.0_linux_amd64.zip"),
						URL:  testutils.GithubAssetURL(10),
					},
					{
						Name: github.Ptr("bluelink-transformer-exampleTransform_1.0.0_windows_amd64.zip"),
						URL:  testutils.GithubAssetURL(11),
					},
					{
						Name: github.Ptr("bluelink-transformer-exampleTransform_1.0.0_registry_info.json"),
						URL:  testutils.GithubAssetURL(12),
					},
				},
			},
			{
				TagName: github.Ptr("v1.1.0"),
				Assets: []*github.ReleaseAsset{
					{
						Name: github.Ptr("bluelink-transformer-exampleTransform_1.1.0_darwin_amd64.zip"),
						URL:  testutils.GithubAssetURL(13),
					},
					{
						Name: github.Ptr("bluelink-transformer-exampleTransform_1.1.0_linux_amd64.zip"),
						URL:  testutils.GithubAssetURL(14),
					},
					{
						Name: github.Ptr("bluelink-transformer-exampleTransform_1.1.0_windows_amd64.zip"),
						URL:  testutils.GithubAssetURL(15),
					},
					{
						Name: github.Ptr("bluelink-transformer-exampleTransform_1.1.0_registry_info.json"),
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
		"supportedProtocols": ["1.4", "2.1"],
		"dependencies": {
			"bluelink/aws": "^1.0.0"
		}
	}
	`)
}

func contentsProvider(url string) ([]byte, error) {
	if strings.HasPrefix(url, packageInfoRegistrySHA256SumsURL()) {
		return packageSHASumContents(), nil
	}

	return registryInfoContents(), nil
}

func packageInfoRegistrySHA256SumsURL() string {
	return "https://artifacts.example.com/bluelink-provider-example/1.0.1/bluelink-provider-example_1.0.1_SHA256SUMS"
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

func TestDefaultServiceTestSuite(t *testing.T) {
	suite.Run(t, new(DefaultServiceTestSuite))
}
