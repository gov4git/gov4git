package ballot

import (
	"context"
	"path/filepath"

	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
	"github.com/gov4git/lib4git/ns"
)

func ListOpen(
	ctx context.Context,
	govAddr gov.CommunityAddress,
) []common.Advertisement {

	_, govTree := git.Clone(ctx, git.Address(govAddr))
	return ListOpenLocal(ctx, govTree)
}

func ListOpenLocal(
	ctx context.Context,
	govTree *git.Tree,
) []common.Advertisement {

	openNS := common.OpenBallotNS(ns.NS(""))

	files, err := git.ListFilesRecursively(govTree, openNS.Path())
	must.NoError(ctx, err)

	ads := []common.Advertisement{}
	for _, file := range files {
		if filepath.Base(file) != common.AdFilebase {
			continue
		}
		var ad common.Advertisement
		err := must.Try(func() { ad = git.FromFile[common.Advertisement](ctx, govTree, file) })
		if err != nil {
			continue
		}
		ads = append(ads, ad)
	}

	return ads
}
