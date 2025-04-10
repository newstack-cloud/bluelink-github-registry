package plugins

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-github/v70/github"
	"github.com/two-hundred/celerity-github-registry/internal/core"
	"github.com/two-hundred/celerity-github-registry/internal/httputils"
	"github.com/two-hundred/celerity-github-registry/internal/repos"
	"github.com/two-hundred/celerity-github-registry/internal/types"
	"github.com/two-hundred/celerity-github-registry/internal/utils"
	"go.uber.org/zap"
)

// PackageInfoParams holds the parameters required
// to retrieve package information for a plugin version.
type PackageInfoParams struct {
	Organisation string
	Plugin       string
	Version      string
	OS           string
	Arch         string
}

// Service provides an interface for a service
// that allows fetching plugin version information
// from private GitHub repositories.
type Service interface {
	// ListVersions lists the versions of a plugin
	// for a given organisation and plugin name
	// that details the supported protocols and platforms
	// for each version.
	ListVersions(
		ctx context.Context,
		organisation string,
		plugin string,
		token string,
	) (*types.PluginVersions, error)

	// GetPackageInfo retrieves the package information
	// for a given plugin version and platform, including
	// the URLs for the package and the SHA256 checksum.
	GetPackageInfo(
		ctx context.Context,
		params *PackageInfoParams,
		token string,
	) (*types.PluginVersionPackage, error)
}

type serviceImpl struct {
	repoService repos.Service
	httpClient  httputils.Client
	config      *core.Config
	logger      *zap.Logger
}

// NewDefaultService creates a new instance of the default
// implementation of a service to retrieve plugin version
// information to fulfil the requirements of the
// Celerity registry protocol.
func NewDefaultService(
	repoService repos.Service,
	httpClient httputils.Client,
	config *core.Config,
	logger *zap.Logger,
) Service {
	return &serviceImpl{
		repoService: repoService,
		config:      config,
		logger:      logger,
		httpClient:  httpClient,
	}
}

func (s *serviceImpl) ListVersions(
	ctx context.Context,
	organisation string,
	plugin string,
	token string,
) (*types.PluginVersions, error) {
	repository, err := s.getPluginRepo(
		ctx,
		organisation,
		plugin,
		token,
	)
	if err != nil {
		return nil, err
	}

	releases, err := s.listReleases(
		ctx,
		organisation,
		repository,
		token,
	)
	if err != nil {
		return nil, err
	}

	return utils.ExtractPluginVersions(
		ctx,
		repository,
		releases,
		s.httpClient,
		token,
	)
}

func (s *serviceImpl) GetPackageInfo(
	ctx context.Context,
	params *PackageInfoParams,
	token string,
) (*types.PluginVersionPackage, error) {
	repository, err := s.getPluginRepo(
		ctx,
		params.Organisation,
		params.Plugin,
		token,
	)
	if err != nil {
		return nil, err
	}

	release, _, err := s.repoService.GetReleaseByTag(
		ctx,
		params.Organisation,
		repository,
		fmt.Sprintf("v%s", params.Version),
		token,
	)
	if err != nil {
		return nil, err
	}

	return utils.ExtractPluginVersionPackage(
		ctx,
		&utils.ExtractPluginVersionPackageParams{
			Repository:            repository,
			Release:               release,
			Version:               params.Version,
			OS:                    params.OS,
			Arch:                  params.Arch,
			SigningKeysSerialised: s.config.PublicSigningKeysString,
		},
		s.httpClient,
		token,
	)
}

func (s *serviceImpl) getPluginRepo(
	ctx context.Context,
	organisation string,
	plugin string,
	token string,
) (string, error) {
	repos, err := s.listRepos(
		ctx,
		organisation,
		token,
	)
	if err != nil {
		return "", err
	}

	repo := utils.FindPluginRepo(
		repos,
		organisation,
		plugin,
	)
	if repo == nil {
		return "", ErrRepoNotFound
	}

	return repo.GetName(), nil
}

func (s *serviceImpl) listRepos(
	ctx context.Context,
	organisation string,
	token string,
) ([]*github.Repository, error) {
	hasMorePages := true
	allRepos := []*github.Repository{}
	for hasMorePages {
		repos, resp, err := s.repoService.ListByOrg(
			ctx,
			organisation,
			&github.RepositoryListByOrgOptions{
				ListOptions: github.ListOptions{
					Page:    1,
					PerPage: 30,
				},
			},
			token,
		)
		if err != nil {
			return nil, handleGitHubErrorResponse(resp, err)
		}
		allRepos = append(allRepos, repos...)
		hasMorePages = resp.NextPage > 0
	}
	return allRepos, nil
}

func (s *serviceImpl) listReleases(
	ctx context.Context,
	organisation string,
	repository string,
	token string,
) ([]*github.RepositoryRelease, error) {
	hasMorePages := true
	allReleases := []*github.RepositoryRelease{}
	for hasMorePages {
		releases, resp, err := s.repoService.ListReleases(
			ctx,
			organisation,
			repository,
			&github.ListOptions{
				Page:    1,
				PerPage: 30,
			},
			token,
		)
		if err != nil {
			return nil, handleGitHubErrorResponse(resp, err)
		}

		hasMorePages = resp.NextPage > 0
		allReleases = append(allReleases, releases...)
	}

	return allReleases, nil
}

func handleGitHubErrorResponse(resp *github.Response, err error) error {
	if resp == nil {
		return err
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return ErrUnauthorised
	}

	if resp.StatusCode == http.StatusForbidden {
		return ErrForbidden
	}

	return err
}
