package ballotapi

import (
	"context"

	"github.com/gov4git/gov4git/v2/proto/ballot/ballotio"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func Show(
	ctx context.Context,
	addr gov.Address,
	ballotName ballotproto.BallotName,
) ballotproto.AdTally {

	return Show_Local(ctx, gov.Clone(ctx, addr).Tree(), ballotName)
}

func Show_Local(
	ctx context.Context,
	t *git.Tree,
	ballotName ballotproto.BallotName,
) ballotproto.AdTally {

	ad, _ := ballotio.LoadStrategy(ctx, t, ballotName)
	var tally ballotproto.Tally
	must.Try(func() { tally = loadTally_Local(ctx, t, ballotName) })
	return ballotproto.AdTally{Ad: ad, Tally: tally}
}
