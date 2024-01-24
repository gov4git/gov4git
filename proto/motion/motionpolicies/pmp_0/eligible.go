package pmp_0

import (
	"context"

	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/motion/motionapi"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
)

func IsConcernProposalEligible(
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

	if !con.IsConcern() {
		return false
	}

	if !prop.IsProposal() {
		return false
	}

	if con.Closed {
		return false
	}

	if prop.Closed {
		return false
	}

	return prop.Score.Attention > 0
}
