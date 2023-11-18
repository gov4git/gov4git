package ops

import (
	"context"

	"github.com/gov4git/gov4git/proto/docket/policy"
	"github.com/gov4git/gov4git/proto/docket/schema"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func FreezeMotion_StageOnly(ctx context.Context, cloned gov.OwnerCloned, id schema.MotionID) git.ChangeNoResult {

	t := cloned.Public.Tree()

	motion := schema.MotionKV.Get(ctx, schema.MotionNS, t, id)
	must.Assert(ctx, !motion.Frozen, schema.ErrMotionAlreadyFrozen)
	motion.Frozen = true
	chg := schema.MotionKV.Set(ctx, schema.MotionNS, t, id, motion)

	// apply policy
	pcy := policy.Get(ctx, motion.Policy.String())
	pcy.Freeze(
		ctx,
		cloned,
		motion,
		policy.MotionPolicyNS(id),
	)

	return chg
}

func UnfreezeMotion_StageOnly(ctx context.Context, cloned gov.OwnerCloned, id schema.MotionID) git.ChangeNoResult {

	t := cloned.Public.Tree()

	motion := schema.MotionKV.Get(ctx, schema.MotionNS, t, id)
	must.Assert(ctx, motion.Frozen, schema.ErrMotionNotFrozen)
	motion.Frozen = false
	chg := schema.MotionKV.Set(ctx, schema.MotionNS, t, id, motion)

	// apply policy
	pcy := policy.Get(ctx, motion.Policy.String())
	pcy.Unfreeze(
		ctx,
		cloned,
		motion,
		policy.MotionPolicyNS(id),
	)

	return chg
}
