package zero

import (
	"context"

	"github.com/gov4git/gov4git/proto/docket/policy"
	"github.com/gov4git/gov4git/proto/docket/schema"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/ns"
)

func init() {
	policy.Install(context.Background(), zeroPolicy{})
}

const ZeroPolicyName = schema.PolicyName("zero-policy")

type zeroPolicy struct{}

func (x zeroPolicy) Name() string {
	return ZeroPolicyName.String()
}

func (x zeroPolicy) Open(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion schema.Motion,
	policyNS ns.NS,

) {
}

func (x zeroPolicy) Score(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion schema.Motion,
	policyNS ns.NS,

) schema.Score {

	return schema.Score{}
}

func (x zeroPolicy) Update(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion schema.Motion,
	policyNS ns.NS,

) {
}

func (x zeroPolicy) Close(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion schema.Motion,
	policyNS ns.NS,

) {
}

func (x zeroPolicy) Cancel(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion schema.Motion,
	policyNS ns.NS,

) {
}

func (x zeroPolicy) Show(
	ctx context.Context,
	cloned gov.Cloned,
	motion schema.Motion,
	policyNS ns.NS,

) form.Map {

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
) {
}

func (x zeroPolicy) AddRefFrom(
	ctx context.Context,
	cloned gov.OwnerCloned,
	refType schema.RefType,
	from schema.Motion,
	to schema.Motion,
	fromPolicyNS ns.NS,
	toPolicyNS ns.NS,
) {
}

func (x zeroPolicy) RemoveRefTo(
	ctx context.Context,
	cloned gov.OwnerCloned,
	refType schema.RefType,
	from schema.Motion,
	to schema.Motion,
	fromPolicyNS ns.NS,
	toPolicyNS ns.NS,
) {
}

func (x zeroPolicy) RemoveRefFrom(
	ctx context.Context,
	cloned gov.OwnerCloned,
	refType schema.RefType,
	from schema.Motion,
	to schema.Motion,
	fromPolicyNS ns.NS,
	toPolicyNS ns.NS,
) {
}

func (x zeroPolicy) Freeze(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion schema.Motion,
	policyNS ns.NS,

) {
}

func (x zeroPolicy) Unfreeze(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion schema.Motion,
	policyNS ns.NS,

) {
}
