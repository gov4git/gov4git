package ballot

import (
	"context"
	"path/filepath"

	"github.com/gov4git/gov4git/mod/gov"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
	"github.com/gov4git/lib4git/ns"
)

func ListOpen[S Strategy](
	ctx context.Context,
	govAddr gov.CommunityAddress,
) []Advertisement {

	_, govTree := git.Clone(ctx, git.Address(govAddr))
	return ListOpenLocal[S](ctx, govTree)
}

func ListOpenLocal[S Strategy](
	ctx context.Context,
	govTree *git.Tree,
) []Advertisement {

	openNS := OpenBallotNS(ns.NS(""))

	files, err := git.ListFilesRecursively(govTree, openNS.Path())
	must.NoError(ctx, err)

	ads := []Advertisement{}
	for _, file := range files {
		if filepath.Base(file) != adFilebase {
			continue
		}
		var ad Advertisement
		err := must.Try(func() { ad = git.FromFile[Advertisement](ctx, govTree, file) })
		if err != nil {
			continue
		}
		ads = append(ads, ad)
	}

	return ads
}
