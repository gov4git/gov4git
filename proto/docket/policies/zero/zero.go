package zero

import (
	"context"

	"github.com/gov4git/gov4git/proto/docket/policy"
	"github.com/gov4git/gov4git/proto/docket/schema"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/notice"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/ns"
)

func init() {
	policy.Install(context.Background(), ZeroPolicyName, zeroPolicy{})
}

const ZeroPolicyName = schema.PolicyName("zero-policy")

type zeroPolicy struct{}

func (x zeroPolicy) Open(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion schema.Motion,
	policyNS ns.NS,
	args ...any,

) (policy.Report, notice.Notices) {

	return nil, notice.Noticef(ctx, "open motion #%v", motion.ID)
}

func (x zeroPolicy) Score(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion schema.Motion,
	policyNS ns.NS,
	args ...any,

) (schema.Score, notice.Notices) {

	return schema.Score{}, notice.Noticef(ctx, "score motion #%v", motion.ID)
}

func (x zeroPolicy) Update(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion schema.Motion,
	policyNS ns.NS,
	args ...any,

) (policy.Report, notice.Notices) {

	return nil, notice.Noticef(ctx, "update motion #%v", motion.ID)
}

func (x zeroPolicy) Close(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion schema.Motion,
	policyNS ns.NS,
	decision schema.Decision,
	args ...any,

) (policy.Report, notice.Notices) {

	return nil, notice.Noticef(ctx, "close motion #%v", motion.ID)
}

func (x zeroPolicy) Cancel(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion schema.Motion,
	policyNS ns.NS,
	args ...any,

) (policy.Report, notice.Notices) {

	return nil, notice.Noticef(ctx, "cancel motion #%v", motion.ID)
}

func (x zeroPolicy) Show(
	ctx context.Context,
	cloned gov.Cloned,
	motion schema.Motion,
	policyNS ns.NS,
	args ...any,

) form.Form {

	return nil
}

func (x zeroPolicy) AddRefTo(
	ctx context.Context,
	cloned gov.OwnerCloned,
	refType schema.RefType,
	from schema.Motion,
	to schema.Motion,
	fromPolicyNS ns.NS,
	toPolicyNS ns.NS,
	args ...any,

) (policy.Report, notice.Notices) {

	return nil, notice.Noticef(ctx, "add %v ref to motion #%v, from motion #%v", refType, to.ID, from.ID)
}

func (x zeroPolicy) AddRefFrom(
	ctx context.Context,
	cloned gov.OwnerCloned,
	refType schema.RefType,
	from schema.Motion,
	to schema.Motion,
	fromPolicyNS ns.NS,
	toPolicyNS ns.NS,
	args ...any,

) (policy.Report, notice.Notices) {

	return nil, notice.Noticef(ctx, "add %v ref from motion #%v, to motion #%v", refType, from.ID, to.ID)
}

func (x zeroPolicy) RemoveRefTo(
	ctx context.Context,
	cloned gov.OwnerCloned,
	refType schema.RefType,
	from schema.Motion,
	to schema.Motion,
	fromPolicyNS ns.NS,
	toPolicyNS ns.NS,
	args ...any,

) (policy.Report, notice.Notices) {

	return nil, notice.Noticef(ctx, "remove %v ref to motion #%v, from motion #%v", refType, to.ID, from.ID)
}

func (x zeroPolicy) RemoveRefFrom(
	ctx context.Context,
	cloned gov.OwnerCloned,
	refType schema.RefType,
	from schema.Motion,
	to schema.Motion,
	fromPolicyNS ns.NS,
	toPolicyNS ns.NS,
	args ...any,

) (policy.Report, notice.Notices) {

	return nil, notice.Noticef(ctx, "remove %v ref from motion #%v, to motion #%v", refType, from.ID, to.ID)
}

func (x zeroPolicy) Freeze(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion schema.Motion,
	policyNS ns.NS,
	args ...any,

) (policy.Report, notice.Notices) {

	return nil, notice.Noticef(ctx, "freeze motion #%v", motion.ID)
}

func (x zeroPolicy) Unfreeze(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion schema.Motion,
	policyNS ns.NS,
	args ...any,

) (policy.Report, notice.Notices) {

	return nil, notice.Noticef(ctx, "unfreeze motion #%v", motion.ID)
}
