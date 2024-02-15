package alcap

import (
	"context"

	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/motion/motionapi"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
)

func AreEligible(
	ctx context.Context,
	cloned gov.Cloned,
	conID motionproto.MotionID,
	propID motionproto.MotionID,
	refType motionproto.RefType,

) bool {

	if refType != ClaimsRefType {
		return false
	}

	con := motionapi.LookupMotion_Local(ctx, cloned, conID)
	prop := motionapi.LookupMotion_Local(ctx, cloned, propID)

	if !con.IsConcern() || con.Policy != ConcernPolicyName {
		return false
	}

	if !prop.IsProposal() || prop.Policy != ProposalPolicyName {
		return false
	}

	if con.Closed {
		return false
	}

	if prop.Closed {
		return false
	}

	propState := motionapi.LoadPolicyState_Local[*ProposalState](ctx, cloned, prop.ID)

	return propState.LatestApprovalScore > 0
}
