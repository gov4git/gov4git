package motionapi

import (
	"context"

	"github.com/gov4git/gov4git/v2/proto/ballot/ballotapi"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/id"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
)

func TrackMotionBatch(
	ctx context.Context,
	addr gov.Address,
	voterAddr id.OwnerAddress,

) motionproto.MotionViews {

	voterOwner := id.CloneOwner(ctx, voterAddr)
	return TrackMotionBatch_Local(ctx, gov.Clone(ctx, addr), voterAddr, voterOwner)
}

func TrackMotionBatch_Local(
	ctx context.Context,
	cloned gov.Cloned,
	voterAddr id.OwnerAddress,
	voterOwner id.OwnerCloned,

) motionproto.MotionViews {

	t := cloned.Tree()
	ids := motionproto.MotionKV.ListKeys(ctx, motionproto.MotionNS, t)
	mvs := make(motionproto.MotionViews, len(ids))
	for i, id := range ids {
		mvs[i] = ShowMotion_Local(ctx, cloned, id)
		if len(mvs[i].Ballots) > 0 {
			vs := ballotapi.Track_StageOnly(
				ctx,
				voterAddr,
				voterOwner,
				cloned,
				mvs[i].Ballots[0].BallotID,
			)
			mvs[i].Voter = &vs
		}
	}
	mvs.Sort()
	return mvs
}
