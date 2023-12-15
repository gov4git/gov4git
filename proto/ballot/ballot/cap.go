package ballot

import (
	"context"

	"github.com/gov4git/gov4git/v2/proto/ballot/common"
	"github.com/gov4git/lib4git/git"
)

func Capitalization_Local(
	ctx context.Context,
	govTree *git.Tree,
	ballotName common.BallotName,
) float64 {

	return Show_Local(ctx, govTree, ballotName).Tally.Capitalization()
}
