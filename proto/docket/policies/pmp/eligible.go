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

) bool {

	propMot := ops.LookupMotion_Local(ctx, cloned, propID)

	if !propMot.IsProposal() {
		return false
	}
	return propMot.Score.Attention > 0
}
