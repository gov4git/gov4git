package ops

import (
	"context"
	"fmt"
	"time"

	"github.com/gov4git/gov4git/proto/docket/policy"
	"github.com/gov4git/gov4git/proto/docket/schema"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func CloseMotion(ctx context.Context, addr gov.GovAddress, id schema.MotionID) git.ChangeNoResult {

	cloned := gov.Clone(ctx, addr)
	return CloseMotion_StageOnly(ctx, cloned.Tree(), id)
}

func CloseMotion_StageOnly(ctx context.Context, t *git.Tree, id schema.MotionID) git.ChangeNoResult {

	motion := schema.MotionKV.Get(ctx, schema.MotionNS, t, id)
	must.Assert(ctx, !motion.Closed, schema.ErrMotionAlreadyClosed)
	motion.Closed = true
	motion.ClosedAt = time.Now()
	schema.MotionKV.Set(ctx, schema.MotionNS, t, id, motion)

	// apply policy
	pcy := policy.Get(ctx, motion.Policy.String())
	pcy.Close(ctx, schema.MotionKV.KeyNS(schema.MotionNS, id).Append("policy"), motion)

	return git.NewChangeNoResult(fmt.Sprintf("Close motion %v", id), "docket_close_motion")
}
