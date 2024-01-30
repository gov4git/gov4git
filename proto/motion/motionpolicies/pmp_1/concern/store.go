package concern

import (
	"context"

	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/motion/motionapi"
)

func LoadClassState_Local(ctx context.Context, cloned gov.OwnerCloned) *PolicyState {
	return motionapi.LoadClassState_Local[*PolicyState](ctx, cloned, ConcernPolicyName)
}

func SaveClassState_StageOnly(ctx context.Context, cloned gov.OwnerCloned, ps *PolicyState) {
	motionapi.SaveClassState_StageOnly[*PolicyState](ctx, cloned, ConcernPolicyName, ps)
}
