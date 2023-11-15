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

func UnlinkMotions(
	ctx context.Context,
	addr gov.OwnerAddress,
	fromID schema.MotionID,
	toID schema.MotionID,
	typ schema.RefType,
) git.Change[form.Map, schema.Ref] {

	cloned := gov.CloneOwner(ctx, addr)
	chg := UnlinkMotions_StageOnly(ctx, cloned, fromID, toID, typ)
	return proto.CommitIfChanged(ctx, cloned.Public, chg)
}

func UnlinkMotions_StageOnly(
	ctx context.Context,
	cloned gov.OwnerCloned,
	fromID schema.MotionID,
	toID schema.MotionID,
	typ schema.RefType,
) git.Change[form.Map, schema.Ref] {

	t := cloned.Public.Tree()

	// read state
	from := schema.MotionKV.Get(ctx, schema.MotionNS, t, fromID)
	to := schema.MotionKV.Get(ctx, schema.MotionNS, t, toID)

	unref := schema.Ref{From: fromID, To: toID, Type: typ}

	// update
	from.RemoveRef(unref)
	to.RemoveRef(unref)

	// write state
	schema.MotionKV.Set(ctx, schema.MotionNS, t, fromID, from)
	schema.MotionKV.Set(ctx, schema.MotionNS, t, toID, to)

	// apply policies
	fromPolicy := policy.Get(ctx, from.Policy.String())
	toPolicy := policy.Get(ctx, to.Policy.String())
	// RemoveRefs are called in the opposite order of AddRefs
	toPolicy.RemoveRefTo(
		ctx,
		cloned,
		unref.Type,
		from,
		to,
		policy.MotionPolicyNS(fromID),
		policy.MotionPolicyNS(toID),
	)
	fromPolicy.RemoveRefFrom(
		ctx,
		cloned,
		unref.Type,
		from,
		to,
		policy.MotionPolicyNS(fromID),
		policy.MotionPolicyNS(toID),
	)

	return git.NewChange(
		fmt.Sprintf("Remove reference from motion %v to motion %v of type %v", fromID, toID, typ),
		"collab_unlink_motions",
		form.Map{"from": fromID, "to": toID, "type": typ},
		unref,
		nil,
	)
}
