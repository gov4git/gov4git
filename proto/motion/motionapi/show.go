package motionapi

import (
	"context"

	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
)

func ShowMotion(
	ctx context.Context,
	addr gov.Address,
	id motionproto.MotionID,
	args ...any,

) motionproto.MotionView {

	return ShowMotion_Local(ctx, gov.Clone(ctx, addr), id, args...)
}

func ShowMotion_Local(
	ctx context.Context,
	cloned gov.Cloned,
	id motionproto.MotionID,
	args ...any,
) motionproto.MotionView {

	t := cloned.Tree()
	m := motionproto.MotionKV.Get(ctx, motionproto.MotionNS, t, id)

	p := motionproto.Get(ctx, m.Policy)
	pv, pb := p.Show(ctx, cloned, m, motionproto.MotionPolicyNS(id), args...)

	return motionproto.MotionView{
		Motion:  m,
		Ballots: pb,
		Policy:  pv,
	}
}
