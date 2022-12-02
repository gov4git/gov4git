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

func List(
	ctx context.Context,
	govAddr gov.GovAddress,
	closed bool,
) []common.Advertisement {

	return ListLocal(ctx, git.Clone(ctx, git.Address(govAddr)).Tree(), closed)
}

func ListLocal(
	ctx context.Context,
	govTree *git.Tree,
	closed bool,
) []common.Advertisement {

	var ballotsNS ns.NS
	if closed {
		ballotsNS = common.ClosedBallotNS(ns.NS{})
	} else {
		ballotsNS = common.OpenBallotNS(ns.NS{})
	}

	files, err := git.ListFilesRecursively(govTree, ballotsNS.Path())
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
