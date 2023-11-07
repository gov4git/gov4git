package ops

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/ballot/ballot"
	"github.com/gov4git/gov4git/proto/ballot/qv"
	"github.com/gov4git/gov4git/proto/docket/schema"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/member"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func FixMotionScore(
	ctx context.Context,
	addr gov.GovAddress,
	id schema.MotionID,
	score float64,
) git.Change[form.Map, schema.Motion] {

	cloned := gov.Clone(ctx, addr)
	chg := FixMotionScore_StageOnly(ctx, addr, cloned, id, score)
	return proto.CommitIfChanged(ctx, cloned, chg)
}

func FixMotionScore_StageOnly(
	ctx context.Context,
	govAddr gov.GovAddress,
	govCloned git.Cloned,
	id schema.MotionID,
	score float64,
) git.Change[form.Map, schema.Motion] {

	motion := schema.MotionKV.Get(ctx, schema.MotionNS, govCloned.Tree(), id)
	motion.Scoring.Fixed = form.Float64(score)
	schema.MotionKV.Set(ctx, schema.MotionNS, govCloned.Tree(), id, motion)

	return git.NewChange(
		fmt.Sprintf("Fix motion %v score to %v", id, score),
		"collab_fix_motion_score",
		form.Map{"motion": motion, "score": score},
		motion,
		nil,
	)
}

func ScoreMotionByPoll(
	ctx context.Context,
	addr gov.GovAddress,
	id schema.MotionID,
) git.Change[form.Map, schema.Motion] {

	cloned := gov.Clone(ctx, addr)
	chg := ScoreMotionByPoll_StageOnly(ctx, addr, cloned, id)
	return proto.CommitIfChanged(ctx, cloned, chg)
}

func ScoreMotionByPoll_StageOnly(
	ctx context.Context,
	govAddr gov.GovAddress,
	govCloned git.Cloned,
	id schema.MotionID,
) git.Change[form.Map, schema.Motion] {

	motion := schema.MotionKV.Get(ctx, schema.MotionNS, govCloned.Tree(), id)

	must.Assertf(ctx, motion.Scoring.Poll == nil, "motion %v is already associated with poll %v", id, motion.Scoring.Poll.OSPath())

	ballotName := schema.MotionPollBallotName(id)
	chg := ballot.Open_StageOnly(
		ctx,
		qv.QV{},
		govAddr,
		govCloned,
		ballotName,
		fmt.Sprintf("Priority poll for motion %v", id),            // title
		fmt.Sprintf("Up/down vote the priority of motion %v", id), // description
		[]string{schema.MotionPollBallotChoice},                   // choices
		member.Everybody,
	)
	motion.Scoring.Poll = &ballotName

	schema.MotionKV.Set(ctx, schema.MotionNS, govCloned.Tree(), id, motion)

	return git.NewChange(
		fmt.Sprintf("Prioritize motion %v by ballot %v", id, ballotName),
		"collab_score_motion_by_poll",
		form.Map{"motion": motion},
		motion,
		form.Forms{chg},
	)
}
