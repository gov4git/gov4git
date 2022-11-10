package ballot

import (
	"context"
	"path/filepath"

	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/lib/must"
	"github.com/gov4git/gov4git/lib/ns"
	"github.com/gov4git/gov4git/mod/gov"
)

func ListOpen[S Strategy](
	ctx context.Context,
	govAddr gov.CommunityAddress,
) []AdForm {

	_, govTree := git.CloneBranchTree(ctx, git.Address(govAddr))
	return ListOpenTree[S](ctx, govTree)
}

func ListOpenTree[S Strategy](
	ctx context.Context,
	govTree *git.Tree,
) []AdForm {

	openNS := OpenBallotNS[S](ns.NS(""))

	files, err := git.ListFilesRecursively(govTree, openNS.Path())
	must.NoError(ctx, err)

	ads := []AdForm{}
	for _, file := range files {
		if filepath.Base(file) != adFilebase {
			continue
		}
		var ad AdForm
		err := must.Try(func() { ad = git.FromFile[AdForm](ctx, govTree, file) })
		if err != nil {
			continue
		}
		ads = append(ads, ad)
	}

	return ads
}
