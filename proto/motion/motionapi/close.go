package motionapi

import (
	"context"
	"time"

	"github.com/gov4git/gov4git/v2/proto"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/history/trace"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
	"github.com/gov4git/gov4git/v2/proto/notice"
	"github.com/gov4git/lib4git/must"
)

func CloseMotion(
	ctx context.Context,
	addr gov.OwnerAddress,
	id motionproto.MotionID,
	decision motionproto.Decision,
	args ...any,

) (motionproto.Report, notice.Notices) {

	cloned := gov.CloneOwner(ctx, addr)
	report, notices := CloseMotion_StageOnly(ctx, cloned, id, decision, args...)
	proto.Commitf(ctx, cloned.Public, "motion_close", "Close motion %v", id)
	return report, notices
}

func CloseMotion_StageOnly(
	ctx context.Context,
	cloned gov.OwnerCloned,
	id motionproto.MotionID,
	decision motionproto.Decision,
	args ...any,

) (motionproto.Report, notice.Notices) {

	t := cloned.Public.Tree()
	motion := motionproto.MotionKV.Get(ctx, motionproto.MotionNS, t, id)
	must.Assert(ctx, !motion.Closed, motionproto.ErrMotionAlreadyClosed)

	// apply policy
	pcy := motionproto.Get(ctx, motion.Policy)
	report, notices := pcy.Close(
		ctx,
		cloned,
		motion,
		motionproto.MotionPolicyNS(id),
		decision,
		args...,
	)
	AppendMotionNotices_StageOnly(ctx, cloned.PublicClone(), id, notices)

	// commit closure
	motion.Closed = true
	motion.ClosedAt = time.Now()
	motionproto.MotionKV.Set(ctx, motionproto.MotionNS, t, id, motion)

	// log
	trace.Log_StageOnly(ctx, cloned.PublicClone(), &trace.Event{
		Op:     "motion_close",
		Args:   trace.M{"id": id},
		Result: trace.M{"motion": motion},
	})

	return report, notices
}
