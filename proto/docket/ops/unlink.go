package ops

import (
	"context"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/docket/policy"
	"github.com/gov4git/gov4git/proto/docket/schema"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/notice"
)

func UnlinkMotions(
	ctx context.Context,
	addr gov.OwnerAddress,
	fromID schema.MotionID,
	toID schema.MotionID,
	typ schema.RefType,
	args ...any,

) (fromReport policy.Report, fromNotices notice.Notices, toReport policy.Report, toNotices notice.Notices) {

	cloned := gov.CloneOwner(ctx, addr)
	ft, fn, tr, tn := UnlinkMotions_StageOnly(ctx, cloned, fromID, toID, typ, args...)
	proto.Commitf(ctx, cloned.PublicClone(), "motion_unlink", "Unlink from motion %v to motion %v as %v", fromID, toID, typ)
	return ft, fn, tr, tn
}

func UnlinkMotions_StageOnly(
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

	unref := schema.Ref{From: fromID, To: toID, Type: typ}

	// update
	from.RemoveRef(unref)
	to.RemoveRef(unref)

	// write state
	schema.MotionKV.Set(ctx, schema.MotionNS, t, fromID, from)
	schema.MotionKV.Set(ctx, schema.MotionNS, t, toID, to)

	// apply policies
	fromPolicy := policy.Get(ctx, from.Policy)
	toPolicy := policy.Get(ctx, to.Policy)
	// RemoveRefs are called in the opposite order of AddRefs
	reportTo, noticesTo := toPolicy.RemoveRefTo(
		ctx,
		cloned,
		unref.Type,
		from,
		to,
		policy.MotionPolicyNS(fromID),
		policy.MotionPolicyNS(toID),
		args...,
	)
	AppendMotionNotices_StageOnly(ctx, cloned.PublicClone(), toID, noticesTo)
	reportFrom, noticesFrom := fromPolicy.RemoveRefFrom(
		ctx,
		cloned,
		unref.Type,
		from,
		to,
		policy.MotionPolicyNS(fromID),
		policy.MotionPolicyNS(toID),
		args...,
	)
	AppendMotionNotices_StageOnly(ctx, cloned.PublicClone(), fromID, noticesFrom)

	return reportFrom, noticesFrom, reportTo, noticesTo
}
