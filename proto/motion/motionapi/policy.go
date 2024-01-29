package motionapi

import (
	"context"

	"github.com/gov4git/gov4git/v2/proto"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
)

func LoadPolicyState[PS form.Form](
	ctx context.Context,
	addr gov.Address,
	id motionproto.MotionID,

) PS {

	cloned := gov.Clone(ctx, addr)
	return LoadPolicyState_Local[PS](ctx, cloned, id)
}

func LoadPolicyState_Local[PS form.Form](
	ctx context.Context,
	cloned gov.Cloned,
	id motionproto.MotionID,

) PS {

	return git.FromFile[PS](ctx, cloned.Tree(), id.PolicyNS())
}

func SavePolicyState[PS form.Form](
	ctx context.Context,
	addr gov.Address,
	id motionproto.MotionID,
	policyState PS,

) {

	cloned := gov.Clone(ctx, addr)
	SavePolicyState_StageOnly[PS](ctx, cloned, id, policyState)
	proto.Commitf(ctx, cloned, "motion_save_policy_state", "update motion policy state")
	cloned.Push(ctx)
}

func SavePolicyState_StageOnly[PS form.Form](
	ctx context.Context,
	cloned gov.Cloned,
	id motionproto.MotionID,
	policyState PS,

) {

	git.ToFileStage[PS](ctx, cloned.Tree(), id.PolicyNS(), policyState)
}
