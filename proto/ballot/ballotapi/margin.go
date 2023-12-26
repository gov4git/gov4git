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
	name ballotproto.BallotName,

) *ballotproto.Margin {

	cloned := gov.Clone(ctx, addr)
	return GetMargin_Local(ctx, cloned, name)
}

func GetMargin_Local(
	ctx context.Context,
	cloned gov.Cloned,
	name ballotproto.BallotName,

) *ballotproto.Margin {

	t := cloned.Tree()
	ad, strategy := ballotio.LoadStrategy(ctx, t, name)
	tally := loadTally_Local(ctx, t, name)
	return strategy.Margin(ctx, cloned, &ad, &tally)
}
