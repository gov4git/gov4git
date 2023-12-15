package ops

import (
	"context"
	"time"

	"github.com/gov4git/gov4git/v2/proto"
	"github.com/gov4git/gov4git/v2/proto/docket/policy"
	"github.com/gov4git/gov4git/v2/proto/docket/schema"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/history"
	"github.com/gov4git/gov4git/v2/proto/notice"
	"github.com/gov4git/lib4git/must"
)

func CloseMotion(
	ctx context.Context,
	addr gov.OwnerAddress,
	id schema.MotionID,
	decision schema.Decision,
	args ...any,

) (policy.Report, notice.Notices) {

	cloned := gov.CloneOwner(ctx, addr)
	report, notices := CloseMotion_StageOnly(ctx, cloned, id, decision, args...)
	proto.Commitf(ctx, cloned.Public, "motion_close", "Close motion %v", id)
	return report, notices
}

func CloseMotion_StageOnly(
	ctx context.Context,
	cloned gov.OwnerCloned,
	id schema.MotionID,
	decision schema.Decision,
	args ...any,

) (policy.Report, notice.Notices) {

	t := cloned.Public.Tree()
	motion := schema.MotionKV.Get(ctx, schema.MotionNS, t, id)
	must.Assert(ctx, !motion.Closed, schema.ErrMotionAlreadyClosed)

	// apply policy
	pcy := policy.Get(ctx, motion.Policy)
	report, notices := pcy.Close(
		ctx,
		cloned,
		motion,
		policy.MotionPolicyNS(id),
		decision,
		args...,
	)
	AppendMotionNotices_StageOnly(ctx, cloned.PublicClone(), id, notices)

	// commit closure
	motion.Closed = true
	motion.ClosedAt = time.Now()
	schema.MotionKV.Set(ctx, schema.MotionNS, t, id, motion)

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
