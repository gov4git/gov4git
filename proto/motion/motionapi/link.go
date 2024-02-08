package motionapi

import (
	"context"

	"github.com/gov4git/gov4git/v2/proto"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
	"github.com/gov4git/gov4git/v2/proto/notice"
)

func LinkMotions(
	ctx context.Context,
	addr gov.OwnerAddress,
	fromID motionproto.MotionID,
	toID motionproto.MotionID,
	typ motionproto.RefType,
	args ...any,

) (fromReport motionproto.Report, fromNotices notice.Notices, toReport motionproto.Report, toNotices notice.Notices) {

	cloned := gov.CloneOwner(ctx, addr)
	ft, fn, tr, tn := LinkMotions_StageOnly(ctx, cloned, fromID, toID, typ, args...)
	proto.Commitf(ctx, cloned.PublicClone(), "motion_link", "Link from motion %v to motion %v as %v", fromID, toID, typ)
	return ft, fn, tr, tn
}

func LinkMotions_StageOnly(
	ctx context.Context,
	cloned gov.OwnerCloned,
	fromID motionproto.MotionID,
	toID motionproto.MotionID,
	typ motionproto.RefType,
	args ...any,

) (fromReport motionproto.Report, fromNotices notice.Notices, toReport motionproto.Report, toNotices notice.Notices) {

	t := cloned.Public.Tree()

	// read state
	from := motionproto.MotionKV.Get(ctx, motionproto.MotionNS, t, fromID)
	to := motionproto.MotionKV.Get(ctx, motionproto.MotionNS, t, toID)

	ref := motionproto.Ref{From: fromID, To: toID, Type: typ}

	// update
	from.AddRefTo(ref)
	to.AddRefBy(ref)

	// write state
	motionproto.MotionKV.Set(ctx, motionproto.MotionNS, t, fromID, from)
	motionproto.MotionKV.Set(ctx, motionproto.MotionNS, t, toID, to)

	// apply policies
	fromPolicy := motionproto.GetPolicy(ctx, from.Policy)
	toPolicy := motionproto.GetPolicy(ctx, to.Policy)

	// AddRefs are called in the opposite order of RemoveRefs
	reportFrom, noticesFrom := fromPolicy.AddRefFrom(
		ctx,
		cloned,
		ref.Type,
		from,
		to,
		args...,
	)
	AppendMotionNotices_StageOnly(ctx, cloned.PublicClone(), fromID, noticesFrom)
	reportTo, noticesTo := toPolicy.AddRefTo(
		ctx,
		cloned,
		ref.Type,
		from,
		to,
		args...,
	)
	AppendMotionNotices_StageOnly(ctx, cloned.PublicClone(), toID, noticesTo)

	// update policy states
	_, fromUpdateNotices := fromPolicy.Update(ctx, cloned, from)
	_, toUpdateNotices := toPolicy.Update(ctx, cloned, to)

	return reportFrom,
		append(noticesFrom, fromUpdateNotices...),
		reportTo,
		append(noticesTo, toUpdateNotices...)
}
