package id

import (
	"context"

	"github.com/gov4git/gov4git/lib/form"
	"github.com/gov4git/gov4git/lib/git"
)

func FetchPublicCredentials(ctx context.Context, publicAddr git.Address) PublicCredentials {
	_, tree := git.CloneBranchTree(ctx, publicAddr)
	return GetPublicCredentials(ctx, tree)
}

func GetPublicCredentials(ctx context.Context, wt *git.Tree) PublicCredentials {
	return form.FromFile[PublicCredentials](ctx, wt.Filesystem, PublicCredentialsNS.Path())
}
