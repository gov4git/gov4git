package motionapi

import (
	"context"
	"time"

	"github.com/gov4git/gov4git/v2/proto"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/history"
	"github.com/gov4git/gov4git/v2/proto/motion/motionpolicy"
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

) (motionpolicy.Report, notice.Notices) {

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

) (motionpolicy.Report, notice.Notices) {

	t := cloned.Public.Tree()
	motion := motionproto.MotionKV.Get(ctx, motionproto.MotionNS, t, id)
	must.Assert(ctx, !motion.Closed, motionproto.ErrMotionAlreadyClosed)

	// apply policy
	pcy := motionpolicy.Get(ctx, motion.Policy)
	report, notices := pcy.Close(
		ctx,
		cloned,
		motion,
		motionpolicy.MotionPolicyNS(id),
		decision,
		args...,
	)
	AppendMotionNotices_StageOnly(ctx, cloned.PublicClone(), id, notices)

	// commit closure
	motion.Closed = true
	motion.ClosedAt = time.Now()
	motionproto.MotionKV.Set(ctx, motionproto.MotionNS, t, id, motion)

	// log
	history.Log_StageOnly(ctx, cloned.PublicClone(), &history.Event{
		Op: &history.Op{
			Op:     "motion_close",
			Args:   history.M{"id": id},
			Result: history.M{"motion": motion},
		},
	})

	return report, notices
}
