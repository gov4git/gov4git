package github

import (
	"context"

	"github.com/google/go-github/v54/github"
	"github.com/gov4git/lib4git/must"
	vendor "github.com/gov4git/vendor4git/github"
)

func Deploy(
	ctx context.Context,
	githubClient *github.Client, // permissions: read project issues, create/write govPrefix
	project GithubRepo,
	govPrefix GithubRepo,
) any {

	// create governance public and private repos
	v := vendor.NewGithubVendorWithClient(ctx, githubClient)

	govPublic := GithubRepo{Owner: govPrefix.Owner, Name: govPrefix.Name + ":gov.public"}
	govPublicURLs, err := v.CreateRepo(ctx, govPublic.Name, govPublic.Owner, false)
	must.NoError(ctx, err)

	govPrivate := GithubRepo{Owner: govPrefix.Owner, Name: govPrefix.Name + ":gov.private"}
	govPrivateURLs, err := v.CreateRepo(ctx, govPrivate.Name, govPrivate.Owner, false)
	must.NoError(ctx, err)

	// initialize governance identity
	XXX

	// install github automation in the public governance repo
	XXX

	// return config for gov4git admin
	panic("XXX")
}
