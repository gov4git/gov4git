package motionapi

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/v2/proto"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/motion/motionpolicy"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
)

func ScoreMotions(
	ctx context.Context,
	addr gov.OwnerAddress,
	args ...any,

) git.Change[form.Map, motionproto.Motions] {

	cloned := gov.CloneOwner(ctx, addr)
	chg := ScoreMotions_StageOnly(ctx, cloned, args...)
	return proto.CommitIfChanged(ctx, cloned.Public, chg)
}

func ScoreMotions_StageOnly(
	ctx context.Context,
	cloned gov.OwnerCloned,
	args ...any,

) git.Change[form.Map, motionproto.Motions] {

	t := cloned.Public.Tree()
	motions := ListMotions_Local(ctx, t)
	for i, motion := range motions {
		// only score open motions
		if motion.Closed {
			continue
		}
		p := motionpolicy.GetMotionPolicy(ctx, motion)
		// NOTE: motion structure may change during scoring (if Score calls motion methods)
		score, notices := p.Score(
			ctx,
			cloned,
			motion,
			motionpolicy.MotionPolicyNS(motions[i].ID),
			args...,
		)
		AppendMotionNotices_StageOnly(ctx, cloned.PublicClone(), motions[i].ID, notices)

		// reload motion, update score and save
		m := motionproto.MotionKV.Get(ctx, motionproto.MotionNS, t, motions[i].ID)
		m.Score = score
		motionproto.MotionKV.Set(ctx, motionproto.MotionNS, t, motions[i].ID, m)
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
