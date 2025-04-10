package repos

import (
	"context"

	"github.com/google/go-github/v70/github"
)

// Service is the interface for interacting with
// GitHub repositories to get information for plugin
// versions.
type Service interface {
	// ListByOrg lists the repositories for an organization.
	//
	// GitHub API docs: https://docs.github.com/rest/repos/repos#list-organization-repositories
	//
	//meta:operation GET /orgs/{org}/repos
	ListByOrg(
		ctx context.Context,
		org string,
		opts *github.RepositoryListByOrgOptions,
		token string,
	) ([]*github.Repository, *github.Response, error)

	// ListReleases lists the releases for a repository.
	//
	// GitHub API docs: https://docs.github.com/rest/releases/releases#list-releases
	//
	//meta:operation GET /repos/{owner}/{repo}/releases
	ListReleases(
		ctx context.Context,
		owner, repo string,
		opts *github.ListOptions,
		token string,
	) ([]*github.RepositoryRelease, *github.Response, error)

	// GetReleaseByTag fetches a release with the specified tag.
	//
	// GitHub API docs: https://docs.github.com/rest/releases/releases#get-a-release-by-tag-name
	//
	//meta:operation GET /repos/{owner}/{repo}/releases/tags/{tag}
	GetReleaseByTag(
		ctx context.Context,
		owner, repo, tag string,
		token string,
	) (*github.RepositoryRelease, *github.Response, error)
}

type githubService struct{}

// NewGitHubService creates a new instance of the GitHub
// service for interacting with GitHub repositories.
func NewGitHubService() Service {
	return &githubService{}
}

func (g *githubService) ListByOrg(
	ctx context.Context,
	org string,
	opts *github.RepositoryListByOrgOptions,
	token string,
) ([]*github.Repository, *github.Response, error) {
	client := github.NewClient(nil).WithAuthToken(token)
	return client.Repositories.ListByOrg(ctx, org, opts)
}

func (g *githubService) ListReleases(
	ctx context.Context,
	owner, repo string,
	opts *github.ListOptions,
	token string,
) ([]*github.RepositoryRelease, *github.Response, error) {
	client := github.NewClient(nil).WithAuthToken(token)
	return client.Repositories.ListReleases(ctx, owner, repo, opts)
}

func (g *githubService) GetReleaseByTag(
	ctx context.Context,
	owner, repo, tag string,
	token string,
) (*github.RepositoryRelease, *github.Response, error) {
	client := github.NewClient(nil).WithAuthToken(token)
	return client.Repositories.GetReleaseByTag(ctx, owner, repo, tag)
}
