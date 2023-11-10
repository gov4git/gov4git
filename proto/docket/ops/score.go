package ops

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/docket/policy"
	"github.com/gov4git/gov4git/proto/docket/schema"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
)

func ScoreMotions(
	ctx context.Context,
	addr gov.GovPrivateAddress,

) git.Change[form.Map, schema.Motions] {

	cloned := gov.CloneOrganizer(ctx, addr)
	chg := ScoreMotions_StageOnly(ctx, addr, cloned)
	return proto.CommitIfChanged(ctx, cloned.Public, chg)
}

func ScoreMotions_StageOnly(
	ctx context.Context,
	addr gov.GovPrivateAddress,
	cloned id.OwnerCloned,

) git.Change[form.Map, schema.Motions] {

	t := cloned.Public.Tree()
	motions := ListMotions_Local(ctx, t)
	for i, motion := range motions {
		p := policy.GetMotionPolicy(ctx, motion)
		motions[i].Score = p.Score(
			ctx,
			addr,
			cloned,
			motion,
			policy.MotionPolicyNS(motions[i].ID),
		)
		schema.MotionKV.Set(ctx, schema.MotionNS, t, motions[i].ID, motions[i])
	}

	motions.Sort()

	return git.NewChange(
		fmt.Sprintf("Score all %d motions", len(motions)),
		"docket_score_motions",
		form.Map{},
		motions,
		form.Forms{},
	)
}
