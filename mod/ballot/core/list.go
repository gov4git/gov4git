package core

import (
	"context"
	"path/filepath"

	"github.com/gov4git/gov4git/mod/ballot/proto"
	"github.com/gov4git/gov4git/mod/gov"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
	"github.com/gov4git/lib4git/ns"
)

func ListOpen(
	ctx context.Context,
	govAddr gov.CommunityAddress,
) []proto.Advertisement {

	_, govTree := git.Clone(ctx, git.Address(govAddr))
	return ListOpenLocal(ctx, govTree)
}

func ListOpenLocal(
	ctx context.Context,
	govTree *git.Tree,
) []proto.Advertisement {

	openNS := proto.OpenBallotNS(ns.NS(""))

	files, err := git.ListFilesRecursively(govTree, openNS.Path())
	must.NoError(ctx, err)

	ads := []proto.Advertisement{}
	for _, file := range files {
		if filepath.Base(file) != proto.AdFilebase {
			continue
		}
		var ad proto.Advertisement
		err := must.Try(func() { ad = git.FromFile[proto.Advertisement](ctx, govTree, file) })
		if err != nil {
			continue
		}
		ads = append(ads, ad)
	}

	return ads
}
