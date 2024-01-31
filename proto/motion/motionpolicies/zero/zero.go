package zero

import (
	"context"

	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/motion"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
	"github.com/gov4git/gov4git/v2/proto/notice"
	"github.com/gov4git/lib4git/form"
)

func init() {
	motionproto.Install(context.Background(), ZeroPolicyName, zeroPolicy{})
}

const ZeroPolicyName = motion.PolicyName("zero-policy")

type zeroPolicy struct{}

func (x zeroPolicy) PostClone(
	ctx context.Context,
	cloned gov.OwnerCloned,
) {
}

func (x zeroPolicy) Open(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion motionproto.Motion,
	args ...any,

) (motionproto.Report, notice.Notices) {

	return nil, notice.Noticef(ctx, "open motion #%v", motion.ID)
}

func (x zeroPolicy) Score(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion motionproto.Motion,
	args ...any,

) (motionproto.Score, notice.Notices) {

	return motionproto.Score{}, notice.Noticef(ctx, "score motion #%v", motion.ID)
}

func (x zeroPolicy) Update(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion motionproto.Motion,
	args ...any,

) (motionproto.Report, notice.Notices) {

	return nil, notice.Noticef(ctx, "update motion #%v", motion.ID)
}

func (x zeroPolicy) Aggregate(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion motionproto.Motions,
) {
}

func (x zeroPolicy) Clear(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion motionproto.Motion,
	args ...any,

) (motionproto.Report, notice.Notices) {

	return nil, nil
}

func (x zeroPolicy) Close(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion motionproto.Motion,
	decision motionproto.Decision,
	args ...any,

) (motionproto.Report, notice.Notices) {

	return nil, notice.Noticef(ctx, "close motion #%v", motion.ID)
}

func (x zeroPolicy) Cancel(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion motionproto.Motion,
	args ...any,

) (motionproto.Report, notice.Notices) {

	return nil, notice.Noticef(ctx, "cancel motion #%v", motion.ID)
}

func (x zeroPolicy) Show(
	ctx context.Context,
	cloned gov.Cloned,
	motion motionproto.Motion,
	args ...any,

) (form.Form, motionproto.MotionBallots) {

	return nil, nil
}

func (x zeroPolicy) AddRefTo(
	ctx context.Context,
	cloned gov.OwnerCloned,
	refType motionproto.RefType,
	from motionproto.Motion,
	to motionproto.Motion,
	args ...any,

) (motionproto.Report, notice.Notices) {

	return nil, notice.Noticef(ctx, "add %v ref to motion #%v, from motion #%v", refType, to.ID, from.ID)
}

func (x zeroPolicy) AddRefFrom(
	ctx context.Context,
	cloned gov.OwnerCloned,
	refType motionproto.RefType,
	from motionproto.Motion,
	to motionproto.Motion,
	args ...any,

) (motionproto.Report, notice.Notices) {

	return nil, notice.Noticef(ctx, "add %v ref from motion #%v, to motion #%v", refType, from.ID, to.ID)
}

func (x zeroPolicy) RemoveRefTo(
	ctx context.Context,
	cloned gov.OwnerCloned,
	refType motionproto.RefType,
	from motionproto.Motion,
	to motionproto.Motion,
	args ...any,

) (motionproto.Report, notice.Notices) {

	return nil, notice.Noticef(ctx, "remove %v ref to motion #%v, from motion #%v", refType, to.ID, from.ID)
}

func (x zeroPolicy) RemoveRefFrom(
	ctx context.Context,
	cloned gov.OwnerCloned,
	refType motionproto.RefType,
	from motionproto.Motion,
	to motionproto.Motion,
	args ...any,

) (motionproto.Report, notice.Notices) {

	return nil, notice.Noticef(ctx, "remove %v ref from motion #%v, to motion #%v", refType, from.ID, to.ID)
}

func (x zeroPolicy) Freeze(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion motionproto.Motion,
	args ...any,

) (motionproto.Report, notice.Notices) {

	return nil, notice.Noticef(ctx, "freeze motion #%v", motion.ID)
}

func (x zeroPolicy) Unfreeze(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion motionproto.Motion,
	args ...any,

) (motionproto.Report, notice.Notices) {

	return nil, notice.Noticef(ctx, "unfreeze motion #%v", motion.ID)
}
