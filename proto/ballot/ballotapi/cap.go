package ballotapi

import (
	"context"

	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/gov"
)

func Capitalization_Local(
	ctx context.Context,
	cloned gov.Cloned,
	ballotName ballotproto.BallotID,
) float64 {

	return Show_Local(ctx, cloned, ballotName).Tally.Capitalization()
}
