package ops

import (
	"context"

	"github.com/gov4git/gov4git/proto/docket/schema"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func FreezeMotion_StageOnly(ctx context.Context, t *git.Tree, id schema.MotionID) git.ChangeNoResult {

	motion := schema.MotionKV.Get(ctx, schema.MotionNS, t, id)
	must.Assert(ctx, !motion.Frozen, schema.ErrMotionAlreadyFrozen)
	motion.Frozen = true
	return schema.MotionKV.Set(ctx, schema.MotionNS, t, id, motion)
}

func UnfreezeMotion_StageOnly(ctx context.Context, t *git.Tree, id schema.MotionID) git.ChangeNoResult {

	motion := schema.MotionKV.Get(ctx, schema.MotionNS, t, id)
	must.Assert(ctx, motion.Frozen, schema.ErrMotionNotFrozen)
	motion.Frozen = false
	return schema.MotionKV.Set(ctx, schema.MotionNS, t, id, motion)
}
