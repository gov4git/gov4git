package collab

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
)

func LinkMotions(
	ctx context.Context,
	addr gov.GovAddress,
	fromID MotionID,
	toID MotionID,
	typ RefType,
) git.Change[form.Map, Ref] {

	cloned := gov.Clone(ctx, addr)
	chg := LinkMotions_StageOnly(ctx, cloned.Tree(), fromID, toID, typ)
	proto.Commit(ctx, cloned.Tree(), chg) // XXX: if no change
	cloned.Push(ctx)
	return chg
}

func LinkMotions_StageOnly(
	ctx context.Context,
	t *git.Tree,
	fromID MotionID,
	toID MotionID,
	typ RefType,
) git.Change[form.Map, Ref] {

	// read state
	from := motionKV.Get(ctx, motionNS, t, fromID)
	to := motionKV.Get(ctx, motionNS, t, toID)

	ref := Ref{From: fromID, To: toID, Type: typ}

	// update
	if !from.RefersTo(toID, typ) {
		from.RefTo = append(from.RefTo, &ref)
	}
	if !to.ReferredBy(fromID, typ) {
		to.RefBy = append(to.RefBy, &ref)
	}

	// write state
	motionKV.Set(ctx, motionNS, t, fromID, from)
	motionKV.Set(ctx, motionNS, t, toID, to)

	return git.NewChange(
		fmt.Sprintf("Motion %v refers to motion %v", fromID, toID),
		"collab_link_motions",
		form.Map{"from": fromID, "to": toID},
		ref,
		nil,
	)
}
