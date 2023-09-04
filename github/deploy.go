package github

import (
	"context"

	"github.com/google/go-github/v54/github"
	"github.com/gov4git/gov4git/gov4git/api"
	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/lib4git/base"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
	vendor "github.com/gov4git/vendor4git/github"
	"golang.org/x/oauth2"
)

func Deploy(
	ctx context.Context,
	token string, // permissions: read project issues, create/write govPrefix
	project GithubRepo,
	govPrefix GithubRepo,
) api.Config {

	// create authenticated GitHub client
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	ghClient := github.NewClient(tc)

	// create governance public and private repos
	v := vendor.NewGithubVendorWithClient(ctx, ghClient)

	govPublic := GithubRepo{Owner: govPrefix.Owner, Name: govPrefix.Name + ":gov.public"}
	base.Infof("creating GitHub repository %v", govPublic)
	govPublicURLs, err := v.CreateRepo(ctx, govPublic.Name, govPublic.Owner, false)
	must.NoError(ctx, err)

	govPrivate := GithubRepo{Owner: govPrefix.Owner, Name: govPrefix.Name + ":gov.private"}
	base.Infof("creating GitHub repository %v", govPrivate)
	govPrivateURLs, err := v.CreateRepo(ctx, govPrivate.Name, govPrivate.Owner, false)
	must.NoError(ctx, err)

	govOwnerAddr := id.OwnerAddress{
		Public: id.PublicAddress{
			Repo:   git.URL(govPublicURLs.HTTPSURL),
			Branch: git.MainBranch,
		},
		Private: id.PrivateAddress{
			Repo:   git.URL(govPrivateURLs.HTTPSURL),
			Branch: git.MainBranch,
		},
	}

	// attach access token authentication to context for git use
	git.SetAuth(ctx, govOwnerAddr.Public.Repo, git.MakeTokenAuth(ctx, token))
	git.SetAuth(ctx, govOwnerAddr.Private.Repo, git.MakeTokenAuth(ctx, token))

	// initialize governance identity
	base.Infof("initializing governance for %v", project)
	id.Init(ctx, govOwnerAddr)

	// install github automation in the public governance repo
	base.Infof("deploying GitHub actions for governance in %v, targetting %v", govPublic, project)
	// XXX

	// return config for gov4git admin
	return api.Config{
		Auth: XXX,
		//
		GovPublicURL:     XXX,
		GovPublicBranch:  XXX,
		GovPrivateURL:    XXX,
		GovPrivateBranch: XXX,
		//
		MemberPublicURL:     XXX,
		MemberPublicBranch:  XXX,
		MemberPrivateURL:    XXX,
		MemberPrivateBranch: XXX,
		//
		CacheDir:        XXX,
		CacheTTLSeconds: XXX,
	}
}
