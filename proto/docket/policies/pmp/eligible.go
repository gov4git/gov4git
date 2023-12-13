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

	if refType != ResolvesRefType {
		return false
	}

	conMot := ops.LookupMotion_Local(ctx, cloned, conID)
	propMot := ops.LookupMotion_Local(ctx, cloned, propID)

	if !conMot.IsConcern() {
		return false
	}

	if !propMot.IsProposal() {
		return false
	}

	if conMot.Closed {
		return false
	}

	if propMot.Closed {
		return false
	}

	return propMot.Score.Attention > 0
}
