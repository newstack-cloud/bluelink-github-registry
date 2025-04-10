package registry

import (
	"context"

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
	return expectedVersions, nil
}
