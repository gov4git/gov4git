package github

import (
	"context"

	"github.com/google/go-github/v55/github"
	"github.com/gov4git/lib4git/must"
)

func fetchOpenIssues(ctx context.Context, repo GithubRepo, ghc *github.Client, labelled string) []*github.Issue {
	opt := &github.IssueListByRepoOptions{State: "open", Labels: []string{labelled}}
	var allIssues []*github.Issue
	for {
		issues, resp, err := ghc.Issues.ListByRepo(ctx, repo.Owner, repo.Name, opt)
		must.NoError(ctx, err)
		allIssues = append(allIssues, issues...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return allIssues
}
