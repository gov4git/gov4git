package github

import (
	"context"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/v55/github"
	"github.com/gov4git/lib4git/base"
	"github.com/gov4git/lib4git/must"
	"github.com/gov4git/lib4git/util"
)

func LoadIssues(
	ctx context.Context,
	ghc *github.Client, // if nil, a new client for repo will be created
	repo Repo,
	loadPR LoadPRFunc,

) (ImportedIssues, map[string]ImportedIssue) {

	if ghc == nil {
		ghc = GetGithubClient(ctx, repo)
	}

	issues := FetchIssues(ctx, repo, ghc)
	key := map[string]ImportedIssue{}
	order := ImportedIssues{}
	for _, issue := range issues {
		ghIssue := TransformIssue(ctx, ghc, repo, issue, loadPR)
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

func LabelsToStrings(labels []*github.Label) []string {
	var labelStrings []string
	for _, label := range labels {
		labelStrings = append(labelStrings, label.GetName())
	}
	slices.Sort(labelStrings)
	return labelStrings
}

func IsIssueForPrioritization(issue *github.Issue) bool {
	return util.IsIn(PrioritizeIssueByGovernanceLabel, LabelsToStrings(issue.Labels)...)
}

func IsIssueManaged(issue *github.Issue) bool {
	return util.IsIn(IssueIsManagedLabel, LabelsToStrings(issue.Labels)...)
}

type LoadPRFunc func(
	ctx context.Context,
	repo Repo,
	issue *github.Issue,
) bool

func TransformIssue(
	ctx context.Context,
	ghc *github.Client,
	repo Repo,
	issue *github.Issue,
	loadPR LoadPRFunc,

) ImportedIssue {

	author, _ := getIssueAuthorLogin(issue)
	var pr *github.PullRequest
	if issue.IsPullRequest() && loadPR(ctx, repo, issue) {
		var err error
		pr, _, err = ghc.PullRequests.Get(ctx, repo.Owner, repo.Name, issue.GetNumber())
		must.NoError(ctx, err)
	}
	return ImportedIssue{
		Managed:           IsIssueManaged(issue),
		ForPrioritization: IsIssueForPrioritization(issue),
		URL:               issue.GetURL(),
		Author:            author,
		Number:            int64(issue.GetNumber()),
		Title:             issue.GetTitle(),
		Body:              issue.GetBody(),
		Labels:            LabelsToStrings(issue.Labels),
		ClosedAt:          unwrapTimestamp(issue.ClosedAt),
		CreatedAt:         unwrapTimestamp(issue.CreatedAt),
		UpdatedAt:         unwrapTimestamp(issue.UpdatedAt),
		Refs:              parseIssueRefs(ctx, repo, issue),
		Locked:            issue.GetLocked(),
		Closed:            issue.GetState() == "closed",
		PullRequest:       issue.IsPullRequest(),
		Merged:            pr != nil && pr.GetMerged(),
	}
}

func unwrapTimestamp(ts *github.Timestamp) *time.Time {
	if ts == nil {
		return nil
	}
	return &ts.Time
}

// parseIssueRefs parses all references to issues or pull requests from the body of an issue.
// Reference directives are of the form: "addresses|resolves|etc. https://github.com/gov4git/testing.project/issues/2"
// References are extracted syntactically and are not guaranteed to correspond to real issues.
func parseIssueRefs(ctx context.Context, repo Repo, issue *github.Issue) []ImportedRef {

	refs := []ImportedRef{}
	matches := refRegexp.FindAllStringSubmatch(issue.GetBody(), -1)
	for _, m := range matches {
		n, err := strconv.Atoi(m[5])
		if err != nil {
			// an attacker could inject invalid github issue links
			base.Infof("reference %q has unparsable issue number %q", m[0], m[5])
			continue
		}
		ref := ImportedRef{To: int64(n), Type: strings.ToLower(m[1])}
		refs = append(refs, ref)
	}
	return refs
}

const refRegexpSrc = `([a-zA-Z0-9\-:_]+)\s+https://github\.com/([a-zA-Z0-9\-]+)/([a-zA-Z0-9\.\-]+)/(issues|pull)/(\d+)`

var refRegexp = regexp.MustCompile(refRegexpSrc)
