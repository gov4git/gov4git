package motionapi

import (
	"context"

	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
	"github.com/gov4git/lib4git/git"
)

func ListMotions(
	ctx context.Context,
	addr gov.Address,

) motionproto.Motions {

	return ListMotions_Local(ctx, gov.Clone(ctx, addr).Tree())
}

func ListMotions_Local(
	ctx context.Context,
	t *git.Tree,

) motionproto.Motions {

	_, motions := motionproto.MotionKV.ListKeyValues(ctx, motionproto.MotionNS, t)
	motionproto.MotionsByID(motions).Sort()
	return motions
}

func ListMotionViews(
	ctx context.Context,
	addr gov.Address,

) motionproto.MotionViews {

	return ListMotionViews_Local(ctx, gov.Clone(ctx, addr))
}

func ListMotionViews_Local(
	ctx context.Context,
	cloned gov.Cloned,

) motionproto.MotionViews {

	t := cloned.Tree()
	ids := motionproto.MotionKV.ListKeys(ctx, motionproto.MotionNS, t)
	mvs := make(motionproto.MotionViews, 0, len(ids))
	for _, id := range ids {
		mv := ShowMotion_Local(ctx, cloned, id)
		if !mv.IsMissingPolicy() {
			mvs = append(mvs, mv)
		}
	}
	mvs.Sort()
	return mvs
}
