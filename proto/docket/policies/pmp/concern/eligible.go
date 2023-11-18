package concern

import (
	"context"

	"github.com/gov4git/gov4git/proto/docket/ops"
	"github.com/gov4git/gov4git/proto/docket/schema"
	"github.com/gov4git/gov4git/proto/gov"
)

func IsProposalEligible(
	ctx context.Context,
	cloned gov.Cloned,
	proposalID schema.MotionID,

) bool {

	mv := ops.ShowMotion_Local(ctx, cloned, proposalID)
	return mv.Motion.Score.Attention > 0 //XXX: use a global threshold param
}
