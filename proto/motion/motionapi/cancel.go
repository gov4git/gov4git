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

func CancelMotion(
	ctx context.Context,
	addr gov.OwnerAddress,
	id motionproto.MotionID,
	args ...any,

) (motionpolicy.Report, notice.Notices) {

	cloned := gov.CloneOwner(ctx, addr)
	report, notices := CancelMotion_StageOnly(ctx, cloned, id, args...)
	proto.Commitf(ctx, cloned.Public, "motion_cancel", "Cancel motion %v", id)
	return report, notices
}

func CancelMotion_StageOnly(
	ctx context.Context,
	cloned gov.OwnerCloned,
	id motionproto.MotionID,
	args ...any,

) (motionpolicy.Report, notice.Notices) {

	t := cloned.Public.Tree()
	motion := motionproto.MotionKV.Get(ctx, motionproto.MotionNS, t, id)
	must.Assertf(ctx, !motion.Closed, "motion %v already closed", motion.ID)
	must.Assertf(ctx, !motion.Cancelled, "motion %v already cancelled", motion.ID)

	// apply policy
	pcy := motionpolicy.Get(ctx, motion.Policy)
	report, notices := pcy.Cancel(
		ctx,
		cloned,
		motion,
		motionpolicy.MotionPolicyNS(id),
		args...,
	)
	AppendMotionNotices_StageOnly(ctx, cloned.PublicClone(), id, notices)

	// commit cancellation
	motion.Closed = true
	motion.Cancelled = true
	motion.ClosedAt = time.Now()
	motionproto.MotionKV.Set(ctx, motionproto.MotionNS, t, id, motion)

	// log
	history.Log_StageOnly(ctx, cloned.PublicClone(), &history.Event{
		Op: &history.Op{
			Op:     "motion_cancel",
			Args:   history.M{"id": id},
			Result: history.M{"motion": motion},
		},
	})

	return report, notices
}
