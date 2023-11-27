package ops

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/docket/policy"
	"github.com/gov4git/gov4git/proto/docket/schema"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
)

func ScoreMotions(
	ctx context.Context,
	addr gov.OwnerAddress,

) git.Change[form.Map, schema.Motions] {

	cloned := gov.CloneOwner(ctx, addr)
	chg := ScoreMotions_StageOnly(ctx, cloned)
	return proto.CommitIfChanged(ctx, cloned.Public, chg)
}

func ScoreMotions_StageOnly(
	ctx context.Context,
	cloned gov.OwnerCloned,

) git.Change[form.Map, schema.Motions] {

	t := cloned.Public.Tree()
	motions := ListMotions_Local(ctx, t)
	for i, motion := range motions {
		// only score open motions
		if motion.Closed {
			continue
		}
		p := policy.GetMotionPolicy(ctx, motion)
		// NOTE: motion structure may change during scoring (if Score calls motion methods)
		score, notices := p.Score(
			ctx,
			cloned,
			motion,
			policy.MotionPolicyNS(motions[i].ID),
		)
		AppendMotionNotices_StageOnly(ctx, cloned.PublicClone(), motions[i].ID, notices)

		// reload motion, update score and save
		m := schema.MotionKV.Get(ctx, schema.MotionNS, t, motions[i].ID)
		m.Score = score
		schema.MotionKV.Set(ctx, schema.MotionNS, t, motions[i].ID, m)
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
