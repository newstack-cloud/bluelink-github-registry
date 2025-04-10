package testutils

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/go-github/v70/github"
)

// StubRepoService is a stub implementation of the
// repos.Service interface for testing purposes.
type StubRepoService struct {
	repos []*github.Repository
	// A mapping of repository names to lists of releases.
	releases map[string][]*github.RepositoryRelease
	// A mapping of repo name and tag in the format `{repo}::{tag}`
	// to the release.
	releaseTagLookup map[string]*github.RepositoryRelease
}

// NewStubRepoService creates a new instance of the
// stub repositories service, this should be used
// instead of instanting the StubRepoService struct directly
// as will prepare internal maps to efficiently lookup items
// in the stub service.
func NewStubRepoService(
	repos []*github.Repository,
	releases map[string][]*github.RepositoryRelease,
) *StubRepoService {
	return &StubRepoService{
		releases:         releases,
		repos:            repos,
		releaseTagLookup: toTagLookup(releases),
	}
}

func (s *StubRepoService) ListByOrg(
	ctx context.Context,
	org string,
	opts *github.RepositoryListByOrgOptions,
	token string,
) ([]*github.Repository, *github.Response, error) {
	orgRepos := []*github.Repository{}

	for _, repo := range s.repos {
		if repo.GetOwner().GetLogin() == org {
			orgRepos = append(orgRepos, repo)
		}
	}

	return orgRepos, &github.Response{
		NextPage: 0,
	}, nil
}

func (s *StubRepoService) ListReleases(
	ctx context.Context,
	owner, repo string,
	opts *github.ListOptions,
	token string,
) ([]*github.RepositoryRelease, *github.Response, error) {
	if repoReleases, ok := s.releases[repo]; ok {
		return repoReleases, &github.Response{
			NextPage: 0,
		}, nil
	}

	return []*github.RepositoryRelease{}, &github.Response{
		NextPage: 0,
	}, nil
}

func (g *StubRepoService) GetReleaseByTag(
	ctx context.Context,
	owner, repo, tag string,
	token string,
) (*github.RepositoryRelease, *github.Response, error) {
	if releases, ok := g.releaseTagLookup[fmt.Sprintf("%s::%s", repo, tag)]; ok {
		return releases, nil, nil
	}

	return nil, &github.Response{
		Response: &http.Response{
			StatusCode: http.StatusNotFound,
		},
	}, errors.New("release not found")
}

func toTagLookup(
	releaseMap map[string][]*github.RepositoryRelease,
) map[string]*github.RepositoryRelease {
	tagLookup := map[string]*github.RepositoryRelease{}

	for repo, releases := range releaseMap {
		for _, release := range releases {
			tagKey := fmt.Sprintf("%s::%s", repo, release.GetTagName())
			tagLookup[tagKey] = release
		}
	}

	return tagLookup
}
