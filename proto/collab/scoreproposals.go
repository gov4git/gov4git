package collab

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/ballot/ballot"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
)

func ScoreProposals(ctx context.Context, addr gov.GovAddress) git.Change[form.Map, Proposals] {

	cloned := gov.Clone(ctx, addr)
	chg := ScoreProposals_StageOnly(ctx, addr, cloned.Tree())
	proto.Commit(ctx, cloned.Tree(), chg)
	cloned.Push(ctx)
	return chg
}

func ScoreProposals_StageOnly(ctx context.Context, govAddr gov.GovAddress, t *git.Tree) git.Change[form.Map, Proposals] {

	proposals := ListProposals_Local(ctx, t)
	for i, proposal := range proposals {
		switch {
		case proposal.Priority.Fixed != nil:
			proposals[i].Score = *proposal.Priority.Fixed
			proposalKV.Set(ctx, proposalNS, t, proposals[i].Name, proposals[i])
		case proposal.Priority.Ballot != nil:
			ast := ballot.Show_Local(ctx, govAddr, t, *proposal.Priority.Ballot)
			proposals[i].Score = ast.Tally.Scores[PriorityBallotChoice]
			proposalKV.Set(ctx, proposalNS, t, proposals[i].Name, proposals[i])
		}
	}

	proposals.Sort()

	return git.NewChange(
		fmt.Sprintf("Score all %d proposals", len(proposals)),
		"collab_score_proposals",
		form.Map{},
		proposals,
		form.Forms{},
	)
}
