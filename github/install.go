package github

import (
	"context"

	"github.com/google/go-github/v54/github"
)

func Install(
	ctx context.Context,
	githubClient *github.Client, // permissions: read project issues, create/write govPrefix
	project GithubRepo,
	govPrefix GithubRepo,
) {
	panic("XXX")
}
