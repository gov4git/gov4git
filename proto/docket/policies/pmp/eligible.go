package pmp

import (
	"context"

	"github.com/gov4git/gov4git/proto/docket/ops"
	"github.com/gov4git/gov4git/proto/docket/schema"
	"github.com/gov4git/gov4git/proto/gov"
)

func IsConcernProposalEligible(
	ctx context.Context,
	cloned gov.Cloned,
	conID schema.MotionID,
	propID schema.MotionID,
	refType schema.RefType,

) bool {

	if refType != ClaimsRefType {
		return false
	}

	con := ops.LookupMotion_Local(ctx, cloned, conID)
	prop := ops.LookupMotion_Local(ctx, cloned, propID)

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
