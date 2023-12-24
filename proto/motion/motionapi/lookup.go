package motionapi

import (
	"context"

	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
)

func LookupMotion(
	ctx context.Context,
	addr gov.Address,
	id motionproto.MotionID,
	args ...any,

) motionproto.Motion {

	return LookupMotion_Local(ctx, gov.Clone(ctx, addr), id, args...)
}

func LookupMotion_Local(
	ctx context.Context,
	cloned gov.Cloned,
	id motionproto.MotionID,
	args ...any,

) motionproto.Motion {

	return motionproto.MotionKV.Get(ctx, motionproto.MotionNS, cloned.Tree(), id)
}
