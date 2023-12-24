package concern

import (
	"context"

	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/ns"
)

func SaveState_StageOnly(ctx context.Context, t *git.Tree, policyNS ns.NS, state *ConcernState) {
	git.ToFileStage[*ConcernState](ctx, t, policyNS.Append(StateFilebase), state)
}

func LoadState_Local(ctx context.Context, t *git.Tree, policyNS ns.NS) *ConcernState {
	state := git.FromFile[ConcernState](ctx, t, policyNS.Append(StateFilebase))
	return &state
}
