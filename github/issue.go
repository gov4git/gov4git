package github

import (
	"context"
	"strconv"
	"time"

	"github.com/google/go-github/v54/github"
	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
	"github.com/gov4git/lib4git/util"
)

func FetchIssues(ctx context.Context, repo GithubRepo) []*github.Issue {

	c := GetGithubClient(ctx, repo)

	opt := &github.IssueListByRepoOptions{State: "all"}
	var allIssues []*github.Issue
	for {
		issues, resp, err := c.Issues.ListByRepo(ctx, repo.Owner, repo.Name, opt)
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

func IssueUsesGovernance(issue *github.Issue) bool {
	return util.IsIn(PrioritizeIssueByGovernanceLabel, labelsToStrings(issue.Labels)...)
}

func TransformIssue(ctx context.Context, issue *github.Issue) GithubBallotIssue {
	return GithubBallotIssue{
		URL:       issue.GetURL(),
		Number:    int64(issue.GetNumber()),
		Title:     issue.GetTitle(),
		Body:      issue.GetBody(),
		Labels:    labelsToStrings(issue.Labels),
		ClosedAt:  unwrapTimestamp(issue.ClosedAt),
		CreatedAt: unwrapTimestamp(issue.CreatedAt),
		UpdatedAt: unwrapTimestamp(issue.UpdatedAt),
		Locked:    issue.GetLocked(),
		Closed:    issue.GetState() == "closed",
	}
}

func unwrapTimestamp(ts *github.Timestamp) *time.Time {
	if ts == nil {
		return nil
	}
	return &ts.Time
}

func LoadBallotIssues(ctx context.Context, repo GithubRepo) GithubBallotIssues {
	issues := FetchIssues(ctx, repo)
	var concerns GithubBallotIssues
	for _, issue := range issues {
		if IssueUsesGovernance(issue) {
			concerns = append(concerns, TransformIssue(ctx, issue))
		}
	}
	return concerns
}

func ImportIssuesUsingPriorityBallots(ctx context.Context, repo GithubRepo, govAddr gov.GovAddress) GithubBallotIssues {
	return ImportIssuesUsingPriorityBallots_Local(ctx, repo, git.CloneOne(ctx, git.Address(govAddr)).Tree())
}

func ImportIssuesUsingPriorityBallots_Local(ctx context.Context, repo GithubRepo, govTree *git.Tree) GithubBallotIssues {

	// ghIssues := LoadBallotIssues(ctx, repo)
	// govBallots := filterIssueBallots(ballot.ListLocal(ctx, govTree))

	// keyIssues := make(map[string]GithubBallotIssue)
	// for _, issue := range ghIssues {
	// 	XXX
	// }

	// keyBallots := make(map[string]common.Advertisement)
	// XXX

}

func filterIssueBallots(ads []common.Advertisement) []common.Advertisement {
	var filtered []common.Advertisement
	for _, ad := range ads {
		if len(ad.Name) == 2 && ad.Name[0] == "issue" {
			if _, err := strconv.Atoi(ad.Name[1]); err == nil {
				filtered = append(filtered, ad)
			}
		}
	}
	return filtered
}
