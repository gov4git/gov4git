package ops

import (
	"context"
	"fmt"
	"time"

	"github.com/gov4git/gov4git/proto/docket/policy"
	"github.com/gov4git/gov4git/proto/docket/schema"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func CloseMotion(
	ctx context.Context,
	addr gov.GovPrivateAddress,
	id schema.MotionID,

) git.ChangeNoResult {

	cloned := gov.CloneOrganizer(ctx, addr)
	return CloseMotion_StageOnly(ctx, addr, cloned, id)
}

func CloseMotion_StageOnly(
	ctx context.Context,
	addr gov.GovPrivateAddress,
	cloned id.OwnerCloned,
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
	pcy.Close(
		ctx,
		addr,
		cloned,
		motion,
		policy.MotionPolicyNS(id),
	)

	return git.NewChangeNoResult(fmt.Sprintf("Close motion %v", id), "docket_close_motion")
}
