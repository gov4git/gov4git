package ballotapi

import (
	"context"

	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/lib4git/git"
)

func Capitalization_Local(
	ctx context.Context,
	govTree *git.Tree,
	ballotName ballotproto.BallotName,
) float64 {

	return Show_Local(ctx, govTree, ballotName).Tally.Capitalization()
}
