package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v55/github"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
	"github.com/gov4git/lib4git/ns"
)

// the dashboard is published as a regularly-updating pinned issue in the project repo

func uploadAssets(
	ctx context.Context,
	addr git.Address, // addr must be a GitHub repo URL
	assets map[string][]byte, // filepath -> content; filepath must have no leading slashes

) (urls map[string]string) { // filepath -> url

	repo, err := ParseGithubRepoURL(string(addr.Repo))
	must.NoError(ctx, err)

	urls = map[string]string{}
	cloned := git.CloneOne(ctx, addr)

	for path, content := range assets {
		git.BytesToFileStage(ctx, cloned.Tree(), ns.ParseFromGitPath(path), content)
		urls[path] = fmt.Sprintf(
			"https://raw.githubusercontent.com/%s/%s/%s/%s",
			repo.Owner,
			repo.Name,
			addr.Branch,
			path,
		)
	}
	git.Commit(ctx, cloned.Tree(), "upload assets")
	cloned.Push(ctx)

	return urls
}

func updateDashboard(
	ctx context.Context,
	ghc *github.Client,
	repo Repo,
	title string,
	body string,

) {

	labels := []string{DashboardIssueLabel}

	// check if there is an existing dashboard issue
	opt := &github.IssueListByRepoOptions{
		State:  "open",
		Labels: labels,
	}
	issues, _, err := ghc.Issues.ListByRepo(ctx, repo.Owner, repo.Name, opt)
	must.NoError(ctx, err)

	// create a dashboard issue if there is none
	req := &github.IssueRequest{
		Title:  github.String(title),
		Body:   github.String(body),
		Labels: &labels,
	}
	if len(issues) == 0 {
		_, _, err := ghc.Issues.Create(ctx, repo.Owner, repo.Name, req)
		must.NoError(ctx, err)
	} else {
		_, _, err := ghc.Issues.Edit(ctx, repo.Owner, repo.Name, issues[0].GetNumber(), req)
		must.NoError(ctx, err)
	}
}
