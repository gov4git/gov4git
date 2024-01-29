package proposal

import (
	"context"

	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/motion/motionpolicies/pmp_1/concern"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
)

func sumClaimedConcernEscrows(
	ctx context.Context,
	cloned gov.OwnerCloned,
	prop motionproto.Motion,
	eligible motionproto.Refs,

) float64 {

	escrow := 0.0
	for _, ref := range eligible {
		conState := concern.LoadMotionPolicyState_Local(
			ctx,
			cloned.PublicClone().Tree(),
			motionproto.MotionPolicyNS(ref.To),
		)
		escrow += conState.PriorityScore // equals concern escrow
	}

	return escrow
}
