package github

import (
	"context"
	"testing"

	"github.com/google/go-github/v58/github"
	govgh "github.com/gov4git/gov4git/v2/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
)

func TestGithubMock(t *testing.T) {

	mockedHTTPClient := mock.NewMockedHTTPClient(
		mock.WithRequestMatch(
			mock.GetReposIssuesByOwnerByRepo,
			[]github.Issue{
				{ // issue without governance
					ID:     github.Int64(123),
					Number: github.Int(1),
					Title:  github.String("Issue 1"),
					URL:    github.String("https://test/issue/1"),
					Labels: []*github.Label{},
					Locked: github.Bool(false),
					//
					State: github.String("open"),
				},
				{ // open, non-frozen issue with governance
					ID:     github.Int64(456),
					Number: github.Int(2),
					Title:  github.String("Issue 2"),
					URL:    github.String("https://test/issue/2"),
					Labels: []*github.Label{{Name: github.String(govgh.PrioritizeIssueByGovernanceLabel)}},
					Locked: github.Bool(false),
					//
					State: github.String("open"),
				},
			},
			[]github.Issue{},
		),
	)

	c := github.NewClient(mockedHTTPClient)

	ctx := context.Background()

	issues1, _, repo1Err := c.Issues.ListByRepo(ctx, "owner1", "repo1", &github.IssueListByRepoOptions{})
	if len(issues1) != 2 || repo1Err != nil {
		t.Errorf("unexpected")
	}

	issues2, _, repo2Err := c.Issues.ListByRepo(ctx, "owner1", "repo2", &github.IssueListByRepoOptions{})
	if len(issues2) != 0 || repo2Err != nil {
		t.Errorf("unexpected")
	}
}
