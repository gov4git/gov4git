package concern

import (
	"context"

	"github.com/gov4git/gov4git/proto/docket/ops"
	"github.com/gov4git/gov4git/proto/docket/schema"
	"github.com/gov4git/gov4git/proto/gov"
)

var AddressesRefType = schema.RefType("gov4git-addresses")

func IsProposalEligible(
	ctx context.Context,
	cloned gov.Cloned,
	proposalID schema.MotionID,

) bool {

	mv := ops.ShowMotion_Local(ctx, cloned, proposalID) //XXX: calls policy show

	if !mv.Motion.IsProposal() {
		return false
	}
	return mv.Motion.Score.Attention > 0
}
