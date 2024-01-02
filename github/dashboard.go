package github

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v55/github"
	"github.com/gov4git/gov4git/v2"
	"github.com/gov4git/gov4git/v2/materials"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/metrics"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
	"github.com/gov4git/lib4git/ns"
)

// the dashboard is published and updated on the first issue that is labelled "gov4git:dashboard"

func PublishDashboard(
	ctx context.Context,
	repo Repo,
	ghc *github.Client,
	cloned gov.Cloned,
) {

	assetsAddr := git.Address{
		Repo:   cloned.Address().Repo,
		Branch: cloned.Address().Branch + ".web-assets",
	}

	assetsRepo, err := ParseGithubRepoURL(string(assetsAddr.Repo))
	must.NoError(ctx, err)

	assets := metrics.AssembleReport_Local(
		ctx,
		cloned,
		func(assetRepoPath string) (url string) {
			return uploadedAssetURL(assetsRepo, string(assetsAddr.Branch), assetRepoPath)
		},
		metrics.TimeDailyLowerBound,
		metrics.Today().AddDate(0, 0, 1),
	)

	uploadAssets(ctx, assetsAddr, assets.Assets)

	header := fmt.Sprintf(
		"## <a href=%q><img src=%q alt=\"This project is governed with Gov4Git.\" width=\"65\" /></a> %s\n"+
			"On `%s` by Gov4Git `%s`\n\n",
		materials.Gov4GitWebsiteURL,
		materials.Gov4GitAvatarURL,
		"Gov4Git community dashboard",
		time.Now().Format(time.RFC850),
		gov4git.GetVersionInfo().Version,
	)
	updateDashboard(ctx, ghc, repo, "Gov4Git community dashboard", header+assets.ReportMD)
}

func uploadedAssetURL(repo Repo, branch string, gitPath string) string {
	return fmt.Sprintf(
		"https://raw.githubusercontent.com/%s/%s/%s/%s",
		repo.Owner,
		repo.Name,
		branch,
		gitPath,
	)
}

func uploadAssets(
	ctx context.Context,
	addr git.Address, // addr must be a GitHub repo URL
	assets map[string][]byte, // git path (in assets repo branch) -> content; git path must have no leading slashes

) {

	cloned := git.CloneOne(ctx, addr)
	for path, content := range assets {
		git.BytesToFileStage(ctx, cloned.Tree(), ns.ParseFromGitPath(path), content)
	}
	git.Commit(ctx, cloned.Tree(), "upload assets")
	cloned.Push(ctx)
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
