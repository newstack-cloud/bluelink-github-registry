package registry

import (
	"context"

	"github.com/two-hundred/celerity-github-registry/internal/plugins"
	"github.com/two-hundred/celerity-github-registry/internal/testutils"
	"github.com/two-hundred/celerity-github-registry/internal/types"
)

type stubPluginService struct{}

var (
	expectedVersions = &types.PluginVersions{
		Versions: []*types.PluginVersion{
			{
				Version:            "3.0.1",
				SupportedProtocols: []string{"1.5", "2.1"},
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
				Version:            "3.1.0",
				SupportedProtocols: []string{"1.5", "2.1"},
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
)

func (s *stubPluginService) ListVersions(
	ctx context.Context,
	organisation string,
	plugin string,
	token string,
) (*types.PluginVersions, error) {
	if plugin != "aws" {
		return nil, plugins.ErrRepoNotFound
	}
	return expectedVersions, nil
}

var (
	expectedVersionPackage = &types.PluginVersionPackage{
		SupportedProtocols:  []string{"1.5", "2.1"},
		OS:                  "linux",
		Arch:                "amd64",
		Filename:            "celerity-provider-aws_3.0.1_linux_amd64.zip",
		DownloadURL:         *testutils.GithubAssetURL(1),
		SHASumsURL:          *testutils.GithubAssetURL(2),
		SHASumsSignatureURL: *testutils.GithubAssetURL(3),
		SHASum:              "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		SigningKeys: &types.PublicGPGSigningKeys{
			GPG: []*types.PublicGPGSigningKey{
				{
					HexKeyID:  "ABCDEF1234567890",
					PublicKey: "example-public-key",
				},
			},
		},
	}
)

func (s *stubPluginService) GetPackageInfo(
	ctx context.Context,
	params *plugins.PackageInfoParams,
	token string,
) (*types.PluginVersionPackage, error) {
	if params.Plugin != "aws" {
		return nil, plugins.ErrRepoNotFound
	}
	return expectedVersionPackage, nil
}
