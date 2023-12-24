package motionapi

import (
	"context"

	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func IsMotion(ctx context.Context, addr gov.Address, id motionproto.MotionID) bool {
	return IsMotion_Local(ctx, gov.Clone(ctx, addr).Tree(), id)
}

func IsMotion_Local(ctx context.Context, t *git.Tree, id motionproto.MotionID) bool {
	err := must.Try(func() { motionproto.MotionKV.Get(ctx, motionproto.MotionNS, t, id) })
	return err == nil
}
