package githubService

import (
	"context"
	"github.com/google/go-github/v43/github"
	"gofuzzyclone/internal/logger"
	"golang.org/x/oauth2"
)

// GithubService is a service for interacting with github
// It is a wrapper around github.Client with a useful logger
type GithubService struct {
	client *github.Client
	logger *logger.Logger
	ctx    context.Context
}

// GetOrgRepos returns a list of repositories for the given organization
func (g *GithubService) GetOrgRepos(org string) ([]*github.Repository, error) {
	// https://github.com/vrenjith/github-pr-manager/blob/5120424b4a9f4ac4675080eebce8582c5905e626/github.go#L20
	max := 1000
	var allRepos []*github.Repository
	opt := &github.RepositoryListByOrgOptions{
		Sort: "updated",
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}
	for {
		repos, resp, err := g.client.Repositories.ListByOrg(context.Background(), org, opt)
		if len(allRepos) > max {
			allRepos = allRepos[:max]
			break
		}
		if err != nil {
			return nil, err
		}
		if len(repos) == 0 {
			break
		}

		opt.Page = resp.NextPage
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
	}
	return allRepos, nil
}

// GetPersonalRepos returns all the repos for the user
func (g *GithubService) GetPersonalRepos(owner string) ([]*github.Repository, error) {
	max := 1000
	var allRepos []*github.Repository
	opt := &github.RepositoryListOptions{
		Sort: "updated",
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}
	for {
		repos, resp, err := g.client.Repositories.List(context.Background(), owner, opt)
		if len(allRepos) > max {
			allRepos = allRepos[:max]
			break
		}
		if err != nil {
			return nil, err
		}
		if len(repos) == 0 {
			break
		}

		opt.Page = resp.NextPage
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
	}
	return allRepos, nil
}

// GetUser returns the user information for the given username
func (g *GithubService) GetUser(user string) (u *github.User, err error) {
	u, _, err = g.client.Users.Get(g.ctx, user)
	return
}

// ValidateToken validates a token
func (g *GithubService) ValidateToken(ghToken string) (err error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: ghToken})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	_, _, err = client.Repositories.List(ctx, "github", nil)
	return
}

// New creates a new GithubService instance
func New(ctx context.Context, ghToken string) (svc *GithubService) {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: ghToken})
	tc := oauth2.NewClient(ctx, ts)
	svc = &GithubService{
		client: github.NewClient(tc),
		logger: logger.New(),
		ctx:    ctx,
	}
	return
}
