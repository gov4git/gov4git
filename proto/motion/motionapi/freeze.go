package motionapi

import (
	"context"

	"github.com/gov4git/gov4git/v2/proto"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/history"
	"github.com/gov4git/gov4git/v2/proto/motion/motionpolicy"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
	"github.com/gov4git/gov4git/v2/proto/notice"
	"github.com/gov4git/lib4git/must"
)

func FreezeMotion(
	ctx context.Context,
	addr gov.OwnerAddress,
	id motionproto.MotionID,
	args ...any,

) (motionpolicy.Report, notice.Notices) {

	cloned := gov.CloneOwner(ctx, addr)
	report, notices := FreezeMotion_StageOnly(ctx, cloned, id, args...)
	proto.Commitf(ctx, cloned.PublicClone(), "motion_freeze", "Freeze motion %v", id)
	return report, notices
}

func FreezeMotion_StageOnly(
	ctx context.Context,
	cloned gov.OwnerCloned,
	id motionproto.MotionID,
	args ...any,

) (motionpolicy.Report, notice.Notices) {

	t := cloned.Public.Tree()

	motion := motionproto.MotionKV.Get(ctx, motionproto.MotionNS, t, id)
	must.Assert(ctx, !motion.Closed, motionproto.ErrMotionAlreadyClosed)
	must.Assert(ctx, !motion.Frozen, motionproto.ErrMotionAlreadyFrozen)

	// apply policy
	pcy := motionpolicy.Get(ctx, motion.Policy)
	report, notices := pcy.Freeze(
		ctx,
		cloned,
		motion,
		motionpolicy.MotionPolicyNS(id),
		args...,
	)
	AppendMotionNotices_StageOnly(ctx, cloned.PublicClone(), id, notices)

	// commit freeze
	motion.Frozen = true
	motionproto.MotionKV.Set(ctx, motionproto.MotionNS, t, id, motion)

	// log
	history.Log_StageOnly(ctx, cloned.PublicClone(), &history.Event{
		Op: &history.Op{
			Op:     "motion_freeze",
			Args:   history.M{"id": id},
			Result: history.M{"motion": motion},
		},
	})

	return report, notices
}

func UnfreezeMotion(
	ctx context.Context,
	addr gov.OwnerAddress,
	id motionproto.MotionID,
	args ...any,

) (motionpolicy.Report, notice.Notices) {

	cloned := gov.CloneOwner(ctx, addr)
	report, notices := UnfreezeMotion_StageOnly(ctx, cloned, id, args...)
	proto.Commitf(ctx, cloned.PublicClone(), "motion_unfreeze", "Unfreeze motion %v", id)
	return report, notices
}

func UnfreezeMotion_StageOnly(
	ctx context.Context,
	cloned gov.OwnerCloned,
	id motionproto.MotionID,
	args ...any,

) (motionpolicy.Report, notice.Notices) {

	t := cloned.Public.Tree()

	motion := motionproto.MotionKV.Get(ctx, motionproto.MotionNS, t, id)
	must.Assert(ctx, !motion.Closed, motionproto.ErrMotionAlreadyClosed)
	must.Assert(ctx, motion.Frozen, motionproto.ErrMotionNotFrozen)

	// apply policy
	pcy := motionpolicy.Get(ctx, motion.Policy)
	report, notices := pcy.Unfreeze(
		ctx,
		cloned,
		motion,
		motionpolicy.MotionPolicyNS(id),
		args...,
	)
	AppendMotionNotices_StageOnly(ctx, cloned.PublicClone(), id, notices)

	// commit unfreeze
	motion.Frozen = false
	motionproto.MotionKV.Set(ctx, motionproto.MotionNS, t, id, motion)

	// log
	history.Log_StageOnly(ctx, cloned.PublicClone(), &history.Event{
		Op: &history.Op{
			Op:     "motion_unfreeze",
			Args:   history.M{"id": id},
			Result: history.M{"motion": motion},
		},
	})

	return report, notices
}
