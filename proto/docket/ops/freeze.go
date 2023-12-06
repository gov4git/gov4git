package ops

import (
	"context"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/docket/policy"
	"github.com/gov4git/gov4git/proto/docket/schema"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/history"
	"github.com/gov4git/gov4git/proto/notice"
	"github.com/gov4git/lib4git/must"
)

func FreezeMotion(
	ctx context.Context,
	addr gov.OwnerAddress,
	id schema.MotionID,
	args ...any,

) (policy.Report, notice.Notices) {

	cloned := gov.CloneOwner(ctx, addr)
	report, notices := FreezeMotion_StageOnly(ctx, cloned, id, args...)
	proto.Commitf(ctx, cloned.PublicClone(), "motion_freeze", "Freeze motion %v", id)
	return report, notices
}

func FreezeMotion_StageOnly(
	ctx context.Context,
	cloned gov.OwnerCloned,
	id schema.MotionID,
	args ...any,

) (policy.Report, notice.Notices) {

	t := cloned.Public.Tree()

	motion := schema.MotionKV.Get(ctx, schema.MotionNS, t, id)
	must.Assert(ctx, !motion.Closed, schema.ErrMotionAlreadyClosed)
	must.Assert(ctx, !motion.Frozen, schema.ErrMotionAlreadyFrozen)
	motion.Frozen = true
	schema.MotionKV.Set(ctx, schema.MotionNS, t, id, motion)

	// apply policy
	pcy := policy.Get(ctx, motion.Policy)
	report, notices := pcy.Freeze(
		ctx,
		cloned,
		motion,
		policy.MotionPolicyNS(id),
		args...,
	)
	AppendMotionNotices_StageOnly(ctx, cloned.PublicClone(), id, notices)

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
	id schema.MotionID,
	args ...any,

) (policy.Report, notice.Notices) {

	cloned := gov.CloneOwner(ctx, addr)
	report, notices := UnfreezeMotion_StageOnly(ctx, cloned, id, args...)
	proto.Commitf(ctx, cloned.PublicClone(), "motion_unfreeze", "Unfreeze motion %v", id)
	return report, notices
}

func UnfreezeMotion_StageOnly(
	ctx context.Context,
	cloned gov.OwnerCloned,
	id schema.MotionID,
	args ...any,

) (policy.Report, notice.Notices) {

	t := cloned.Public.Tree()

	motion := schema.MotionKV.Get(ctx, schema.MotionNS, t, id)
	must.Assert(ctx, !motion.Closed, schema.ErrMotionAlreadyClosed)
	must.Assert(ctx, motion.Frozen, schema.ErrMotionNotFrozen)
	motion.Frozen = false
	schema.MotionKV.Set(ctx, schema.MotionNS, t, id, motion)

	// apply policy
	pcy := policy.Get(ctx, motion.Policy)
	report, notices := pcy.Unfreeze(
		ctx,
		cloned,
		motion,
		policy.MotionPolicyNS(id),
		args...,
	)
	AppendMotionNotices_StageOnly(ctx, cloned.PublicClone(), id, notices)

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
