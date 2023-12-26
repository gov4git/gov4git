package ballotapi

import (
	"context"

	"github.com/gov4git/gov4git/v2/proto/ballot/ballotio"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func Show(ctx context.Context, govAddr gov.Address, ballotName ballotproto.BallotName) ballotproto.AdTally {

	return Show_Local(ctx, gov.Clone(ctx, govAddr).Tree(), ballotName)
}

func Show_Local(
	ctx context.Context,
	govTree *git.Tree,
	ballotName ballotproto.BallotName,
) ballotproto.AdTally {

	ad, _ := ballotio.LoadStrategy(ctx, govTree, ballotName)
	var tally ballotproto.Tally
	must.Try(func() { tally = loadTally_Local(ctx, govTree, ballotName) })
	return ballotproto.AdTally{Ad: ad, Tally: tally}
}
