package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v55/github"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func Import(
	ctx context.Context,
	repo GithubRepo,
	githubClient *github.Client,
	govAddr gov.OrganizerAddress,
) git.Change[form.Map, form.Map] {

	var chg1 git.Change[map[string]form.Form, GithubIssueBallots]
	err1 := must.Try(func() {
		chg1 = ImportIssuesForPrioritization(ctx, repo, githubClient, govAddr)
	})

	var chg2 git.Change[map[string]form.Form, ProcessJoinRequestIssuesReport]
	err2 := must.Try(func() {
		chg2 = ProcessJoinRequestIssuesApprovedByMaintainer(ctx, repo, githubClient, govAddr)
	})

	return git.NewChange[form.Map, form.Map](
		fmt.Sprintf("Import from GitHub"),
		"github_import",
		form.Map{},
		form.Map{
			"issues_for_prioritization":       chg1.Result,
			"join_requests":                   chg2.Result,
			"issues_for_prioritization_error": err1,
			"join_requests_error":             err2,
		},
		nil,
	)
}
