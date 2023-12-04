package ballot

import (
	"context"

	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/ballot/load"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func Show(ctx context.Context, govAddr gov.Address, ballotName common.BallotName) common.AdTally {

	return Show_Local(ctx, gov.Clone(ctx, govAddr).Tree(), ballotName)
}

func Show_Local(
	ctx context.Context,
	govTree *git.Tree,
	ballotName common.BallotName,
) common.AdTally {

	ad, _ := load.LoadStrategy(ctx, govTree, ballotName)
	var tally common.Tally
	must.Try(func() { tally = LoadTally(ctx, govTree, ballotName) })
	return common.AdTally{Ad: ad, Tally: tally}
}
