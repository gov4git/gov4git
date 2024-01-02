package ballotapi

import (
	"context"

	"github.com/gov4git/gov4git/v2/proto/ballot/ballotio"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/gov"
)

func GetMargin(
	ctx context.Context,
	addr gov.Address,
	id ballotproto.BallotID,

) *ballotproto.Margin {

	cloned := gov.Clone(ctx, addr)
	return GetMargin_Local(ctx, cloned, id)
}

func GetMargin_Local(
	ctx context.Context,
	cloned gov.Cloned,
	id ballotproto.BallotID,

) *ballotproto.Margin {

	t := cloned.Tree()
	ad, policy := ballotio.LoadPolicy(ctx, t, id)
	tally := loadTally_Local(ctx, t, id)
	return policy.Margin(ctx, cloned, &ad, &tally)
}
