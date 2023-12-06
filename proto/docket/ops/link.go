package ops

import (
	"context"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/docket/policy"
	"github.com/gov4git/gov4git/proto/docket/schema"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/notice"
)

func LinkMotions(
	ctx context.Context,
	addr gov.OwnerAddress,
	fromID schema.MotionID,
	toID schema.MotionID,
	typ schema.RefType,
	args ...any,

) (fromReport policy.Report, fromNotices notice.Notices, toReport policy.Report, toNotices notice.Notices) {

	cloned := gov.CloneOwner(ctx, addr)
	ft, fn, tr, tn := LinkMotions_StageOnly(ctx, cloned, fromID, toID, typ, args...)
	proto.Commitf(ctx, cloned.PublicClone(), "motion_link", "Link from motion %v to motion %v as %v", fromID, toID, typ)
	return ft, fn, tr, tn
}

func LinkMotions_StageOnly(
	ctx context.Context,
	cloned gov.OwnerCloned,
	fromID schema.MotionID,
	toID schema.MotionID,
	typ schema.RefType,
	args ...any,

) (fromReport policy.Report, fromNotices notice.Notices, toReport policy.Report, toNotices notice.Notices) {

	t := cloned.Public.Tree()

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

	// apply policies
	fromPolicy := policy.Get(ctx, from.Policy)
	toPolicy := policy.Get(ctx, to.Policy)
	// AddRefs are called in the opposite order of RemoveRefs
	reportFrom, noticesFrom := fromPolicy.AddRefFrom(
		ctx,
		cloned,
		ref.Type,
		from,
		to,
		policy.MotionPolicyNS(fromID),
		policy.MotionPolicyNS(toID),
		args...,
	)
	AppendMotionNotices_StageOnly(ctx, cloned.PublicClone(), fromID, noticesFrom)
	reportTo, noticesTo := toPolicy.AddRefTo(
		ctx,
		cloned,
		ref.Type,
		from,
		to,
		policy.MotionPolicyNS(fromID),
		policy.MotionPolicyNS(toID),
		args...,
	)
	AppendMotionNotices_StageOnly(ctx, cloned.PublicClone(), toID, noticesTo)

	return reportFrom, noticesFrom, reportTo, noticesTo
}
