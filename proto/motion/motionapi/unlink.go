package motionapi

import (
	"context"

	"github.com/gov4git/gov4git/v2/proto"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
	"github.com/gov4git/gov4git/v2/proto/notice"
)

func UnlinkMotions(
	ctx context.Context,
	addr gov.OwnerAddress,
	fromID motionproto.MotionID,
	toID motionproto.MotionID,
	typ motionproto.RefType,
	args ...any,

) (fromReport motionproto.Report, fromNotices notice.Notices, toReport motionproto.Report, toNotices notice.Notices) {

	cloned := gov.CloneOwner(ctx, addr)
	ft, fn, tr, tn := UnlinkMotions_StageOnly(ctx, cloned, fromID, toID, typ, args...)
	proto.Commitf(ctx, cloned.PublicClone(), "motion_unlink", "Unlink from motion %v to motion %v as %v", fromID, toID, typ)
	return ft, fn, tr, tn
}

func UnlinkMotions_StageOnly(
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

	unref := motionproto.Ref{From: fromID, To: toID, Type: typ}

	// update
	from.RemoveRef(unref)
	to.RemoveRef(unref)

	// write state
	motionproto.MotionKV.Set(ctx, motionproto.MotionNS, t, fromID, from)
	motionproto.MotionKV.Set(ctx, motionproto.MotionNS, t, toID, to)

	// apply policies
	fromPolicy := motionproto.Get(ctx, from.Policy)
	toPolicy := motionproto.Get(ctx, to.Policy)
	// RemoveRefs are called in the opposite order of AddRefs
	reportTo, noticesTo := toPolicy.RemoveRefTo(
		ctx,
		cloned,
		unref.Type,
		from,
		to,
		args...,
	)
	AppendMotionNotices_StageOnly(ctx, cloned.PublicClone(), toID, noticesTo)
	reportFrom, noticesFrom := fromPolicy.RemoveRefFrom(
		ctx,
		cloned,
		unref.Type,
		from,
		to,
		args...,
	)
	AppendMotionNotices_StageOnly(ctx, cloned.PublicClone(), fromID, noticesFrom)

	return reportFrom, noticesFrom, reportTo, noticesTo
}
