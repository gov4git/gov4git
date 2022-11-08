package id

import (
	"context"

	"github.com/gov4git/gov4git/lib/form"
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/lib/ns"
)

type PublicMod struct {
	ns.NS
	Public git.Address
}

func Public(repo git.Address) PublicMod {
	return PublicMod{NS: PublicNS, Public: repo}
}

func FetchPublicCredentials(ctx context.Context, m PublicMod) PublicCredentials {
	return GetPublicCredentials(ctx, m, git.CloneBranchTree(ctx, m.Public))
}

func GetPublicCredentials(ctx context.Context, m PublicMod, wt *git.Tree) PublicCredentials {
	return form.FromFile[PublicCredentials](ctx, wt.Filesystem, m.Sub(PublicCredentialsFilebase).Path())
}
