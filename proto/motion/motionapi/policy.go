package motionapi

import (
	"context"

	"github.com/gov4git/gov4git/v2/proto"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/motion"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
)

// instance state

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

	return git.FromFile[PS](ctx, cloned.Tree(), id.PolicyNS(motionproto.PolicyStateFilebase))
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

	git.ToFileStage[PS](ctx, cloned.Tree(), id.PolicyNS(motionproto.PolicyStateFilebase), policyState)
}

// class state

func LoadClassState[PS form.Form](
	ctx context.Context,
	addr gov.OwnerAddress,
	policyName motion.PolicyName,

) PS {

	cloned := gov.CloneOwner(ctx, addr)
	return LoadClassState_Local[PS](ctx, cloned, policyName)
}

func LoadClassState_Local[PS form.Form](
	ctx context.Context,
	cloned gov.OwnerCloned,
	policyName motion.PolicyName,

) PS {

	return git.FromFile[PS](
		ctx,
		cloned.Public.Tree(),
		motionproto.PolicyNS(policyName).Append(motionproto.PolicyStateFilebase),
	)
}

func SaveClassState[PS form.Form](
	ctx context.Context,
	addr gov.OwnerAddress,
	policyName motion.PolicyName,
	policyState PS,

) {

	cloned := gov.CloneOwner(ctx, addr)
	SaveClassState_StageOnly[PS](ctx, cloned, policyName, policyState)
	proto.Commitf(ctx, cloned.PublicClone(), "motion_save_policy_class_state", "update motion policy class state")
	cloned.PublicClone().Push(ctx)
}

func SaveClassState_StageOnly[PS form.Form](
	ctx context.Context,
	cloned gov.OwnerCloned,
	policyName motion.PolicyName,
	policyState PS,

) {

	git.ToFileStage[PS](
		ctx,
		cloned.PublicClone().Tree(),
		motionproto.PolicyNS(policyName).Append(motionproto.PolicyStateFilebase),
		policyState,
	)
}

func SupportedPolicies(ctx context.Context) map[string]motionproto.PolicyDescriptor {

	return motionproto.InstalledPolicyDescriptors()
}
