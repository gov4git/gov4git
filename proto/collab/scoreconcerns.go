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

func ScoreConcerns(ctx context.Context, addr gov.GovAddress) git.Change[form.Map, Concerns] {

	cloned := gov.Clone(ctx, addr)
	chg := ScoreConcerns_StageOnly(ctx, addr, cloned.Tree())
	proto.Commit(ctx, cloned.Tree(), chg)
	cloned.Push(ctx)
	return chg
}

func ScoreConcerns_StageOnly(ctx context.Context, govAddr gov.GovAddress, t *git.Tree) git.Change[form.Map, Concerns] {

	concerns := ListConcerns_Local(ctx, t)
	for i, concern := range concerns {
		switch {
		case concern.Priority.Fixed != nil:
			concerns[i].Score = *concern.Priority.Fixed
			concernKV.Set(ctx, concernNS, t, concerns[i].Name, concerns[i])
		case concern.Priority.Ballot != nil:
			ast := ballot.Show_Local(ctx, govAddr, t, *concern.Priority.Ballot)
			concerns[i].Score = ast.Tally.Scores[PriorityBallotChoice]
			concernKV.Set(ctx, concernNS, t, concerns[i].Name, concerns[i])
		}
	}

	concerns.Sort()

	return git.NewChange(
		fmt.Sprintf("Score all %d concerns", len(concerns)),
		"collab_score_concerns",
		form.Map{},
		concerns,
		form.Forms{},
	)
}
