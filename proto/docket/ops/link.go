package ops

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/docket/schema"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
)

func LinkMotions(
	ctx context.Context,
	addr gov.GovAddress,
	fromID schema.MotionID,
	toID schema.MotionID,
	typ schema.RefType,
) git.Change[form.Map, schema.Ref] {

	cloned := gov.Clone(ctx, addr)
	chg := LinkMotions_StageOnly(ctx, cloned.Tree(), fromID, toID, typ)
	return proto.CommitIfChanged(ctx, cloned, chg)
}

func LinkMotions_StageOnly(
	ctx context.Context,
	t *git.Tree,
	fromID schema.MotionID,
	toID schema.MotionID,
	typ schema.RefType,
) git.Change[form.Map, schema.Ref] {

	// read state
	from := schema.MotionKV.Get(ctx, schema.MotionNS, t, fromID)
	to := schema.MotionKV.Get(ctx, schema.MotionNS, t, toID)

	ref := schema.Ref{From: fromID, To: toID, Type: typ}

	// update
	from.AddRefTo(ref)
	to.AddRefBy(ref)

	// write state
	schema.MotionKV.Set(ctx, schema.MotionNS, t, fromID, from)
	schema.MotionKV.Set(ctx, schema.MotionNS, t, toID, to)

	return git.NewChange(
		fmt.Sprintf("Add reference from motion %v to motion %v of type %v", fromID, toID, typ),
		"collab_link_motions",
		form.Map{"from": fromID, "to": toID, "type": typ},
		ref,
		nil,
	)
}
