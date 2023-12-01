package ops

import (
	"context"
	"fmt"
	"time"

	"github.com/gov4git/gov4git/proto/docket/policy"
	"github.com/gov4git/gov4git/proto/docket/schema"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/history"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func CancelMotion(
	ctx context.Context,
	addr gov.OwnerAddress,
	id schema.MotionID,

) git.ChangeNoResult {

	cloned := gov.CloneOwner(ctx, addr)
	return CancelMotion_StageOnly(ctx, cloned, id)
}

func CancelMotion_StageOnly(
	ctx context.Context,
	cloned gov.OwnerCloned,
	id schema.MotionID,

) git.ChangeNoResult {

	t := cloned.Public.Tree()
	motion := schema.MotionKV.Get(ctx, schema.MotionNS, t, id)
	must.Assert(ctx, !motion.Closed, schema.ErrMotionAlreadyClosed)
	motion.Closed = true
	motion.ClosedAt = time.Now()
	schema.MotionKV.Set(ctx, schema.MotionNS, t, id, motion)

	// apply policy
	pcy := policy.Get(ctx, motion.Policy.String())
	notices := pcy.Cancel(
		ctx,
		cloned,
		motion,
		policy.MotionPolicyNS(id),
	)
	AppendMotionNotices_StageOnly(ctx, cloned.PublicClone(), id, notices)

	// log
	history.Log_StageOnly(ctx, cloned.PublicClone(), &history.Event{
		Op: &history.Op{
			Op:     "motion_cancel",
			Args:   history.M{"id": id},
			Result: history.M{"motion": motion},
		},
	})

	return git.NewChangeNoResult(fmt.Sprintf("Cancel motion %v", id), "docket_cancel_motion")
}
