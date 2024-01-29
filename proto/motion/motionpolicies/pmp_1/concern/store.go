package concern

import (
	"context"

	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/ns"
)

// policy

var (
	PolicyNS = motionproto.PolicyNS(ConcernPolicyName)
)

func LoadPolicyState_Local(ctx context.Context, cloned gov.OwnerCloned) *PolicyState {

	return git.FromFile[*PolicyState](ctx, cloned.PublicClone().Tree(), PolicyNS.Append(StateFilebase))
}

func SavePolicyState_StageOnly(ctx context.Context, cloned gov.OwnerCloned, ps *PolicyState) {

	git.ToFileStage[*PolicyState](ctx, cloned.PublicClone().Tree(), PolicyNS.Append(StateFilebase), ps)
}

// instance

func SaveInstanceState_StageOnly(ctx context.Context, t *git.Tree, policyNS ns.NS, state *ConcernState) {
	git.ToFileStage[*ConcernState](ctx, t, policyNS.Append(StateFilebase), state)
}

func LoadInstanceState_Local(ctx context.Context, t *git.Tree, policyNS ns.NS) *ConcernState {
	state := git.FromFile[ConcernState](ctx, t, policyNS.Append(StateFilebase))
	return &state
}
