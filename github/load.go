package github

import (
	"context"
	"time"

	"github.com/google/go-github/v55/github"
	"github.com/gov4git/lib4git/must"
	"github.com/gov4git/lib4git/util"
)

func LoadIssues(
	ctx context.Context,
	repo Repo,
	githubClient *github.Client, // if nil, a new client for repo will be created
) (ImportedIssues, map[string]ImportedIssue) {

	if githubClient == nil {
		githubClient = GetGithubClient(ctx, repo)
	}

	issues := FetchIssues(ctx, repo, githubClient)
	key := map[string]ImportedIssue{}
	order := ImportedIssues{}
	for _, issue := range issues {
		ghIssue := TransformIssue(ctx, issue)
		key[ghIssue.Key()] = ghIssue
		order = append(order, ghIssue)
	}
	order.Sort()
	return order, key
}

func FetchIssues(ctx context.Context, repo Repo, ghc *github.Client) []*github.Issue {

	opt := &github.IssueListByRepoOptions{State: "all"}
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

func labelsToStrings(labels []*github.Label) []string {
	var labelStrings []string
	for _, label := range labels {
		labelStrings = append(labelStrings, label.GetName())
	}
	return labelStrings
}

func IsIssueForPrioritization(issue *github.Issue) bool {
	return util.IsIn(PrioritizeIssueByGovernanceLabel, labelsToStrings(issue.Labels)...)
}

func IsIssueGoverned(issue *github.Issue) bool {
	return util.IsIn(IssueIsGovernedLabel, labelsToStrings(issue.Labels)...)
}

func TransformIssue(ctx context.Context, issue *github.Issue) ImportedIssue {
	return ImportedIssue{
		IsGoverned:        IsIssueGoverned(issue),
		ForPrioritization: IsIssueForPrioritization(issue),
		URL:               issue.GetURL(),
		Number:            int64(issue.GetNumber()),
		Title:             issue.GetTitle(),
		Body:              issue.GetBody(),
		Labels:            labelsToStrings(issue.Labels),
		ClosedAt:          unwrapTimestamp(issue.ClosedAt),
		CreatedAt:         unwrapTimestamp(issue.CreatedAt),
		UpdatedAt:         unwrapTimestamp(issue.UpdatedAt),
		Locked:            issue.GetLocked(),
		Closed:            issue.GetState() == "closed",
		IsPullRequest:     issue.GetPullRequestLinks() != nil,
	}
}

func unwrapTimestamp(ts *github.Timestamp) *time.Time {
	if ts == nil {
		return nil
	}
	return &ts.Time
}
