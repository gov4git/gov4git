package github

import (
	"context"

	"github.com/google/go-github/v55/github"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/lib4git/base"
	"github.com/gov4git/lib4git/form"
)

func Import(
	ctx context.Context,
	repo Repo,
	githubClient *github.Client,
	govAddr gov.OrganizerAddress,
) form.Map {

	chg1 := ImportIssuesForPrioritization(ctx, repo, githubClient, govAddr)
	chg2 := ImportJoinsAndDirectives(ctx, repo, githubClient, govAddr)

	return form.Map{
		"issues_for_prioritization": chg1.Result,
		"join_requests":             chg2,
	}
}

func ImportJoinsAndDirectives(
	ctx context.Context,
	repo Repo,
	ghc *github.Client,
	govAddr gov.OrganizerAddress,
) form.Map {

	base.Infof("importing join requests and organizer directives ...")

	maintainers := FetchRepoMaintainers(ctx, repo, ghc)
	base.Infof("maintainers for %v are %v", repo, form.SprintJSON(maintainers))

	chg1 := ProcessJoinRequestIssues(ctx, repo, ghc, govAddr, maintainers)
	chg2 := ProcessDirectiveIssues(ctx, repo, ghc, govAddr, maintainers)
	return form.Map{
		"joins":      chg1.Result,
		"directives": chg2.Result,
	}
}
