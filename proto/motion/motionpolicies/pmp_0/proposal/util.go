package proposal

import (
	"context"

	"github.com/gov4git/gov4git/v2/proto/motion/motionpolicies/pmp_0"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/ns"
)

func SaveState_StageOnly(ctx context.Context, t *git.Tree, policyNS ns.NS, state *pmp_0.ProposalState) {
	git.ToFileStage[*pmp_0.ProposalState](ctx, t, policyNS.Append(pmp_0.StateFilebase), state)
}

func LoadState_Local(ctx context.Context, t *git.Tree, policyNS ns.NS) *pmp_0.ProposalState {
	state := git.FromFile[pmp_0.ProposalState](ctx, t, policyNS.Append(pmp_0.StateFilebase))
	return &state
}
