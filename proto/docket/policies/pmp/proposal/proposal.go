package proposal

import (
	"context"

	"github.com/gov4git/gov4git/proto/docket/policy"
	"github.com/gov4git/gov4git/proto/docket/schema"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/ns"
)

func init() {
	policy.Install(context.Background(), proposalPolicy{})
}

const ProposalPolicyName = schema.PolicyName("pmp-proposal-policy")

type proposalPolicy struct{}

func (x proposalPolicy) Name() string {
	return ProposalPolicyName.String()
}

func (x proposalPolicy) Open(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion schema.Motion,
	instancePolicyNS ns.NS,

) {

}

func (x proposalPolicy) Score(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion schema.Motion,
	instancePolicyNS ns.NS,

) schema.Score {

	return schema.Score{}
}

func (x proposalPolicy) Close(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion schema.Motion,
	instancePolicyNS ns.NS,

) {

}

func (x proposalPolicy) Cancel(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion schema.Motion,
	instancePolicyNS ns.NS,

) {

}

func (x proposalPolicy) Show(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion schema.Motion,
	instancePolicyNS ns.NS,

) form.Map {

	return nil
}

func (x proposalPolicy) AddRefTo(
	ctx context.Context,
	cloned gov.OwnerCloned,
	refType schema.RefType,
	from schema.Motion,
	to schema.Motion,
	fromPolicyNS ns.NS,
	toPolicyNS ns.NS,
) {

}

func (x proposalPolicy) AddRefFrom(
	ctx context.Context,
	cloned gov.OwnerCloned,
	refType schema.RefType,
	from schema.Motion,
	to schema.Motion,
	fromPolicyNS ns.NS,
	toPolicyNS ns.NS,
) {

}

func (x proposalPolicy) RemoveRefTo(
	ctx context.Context,
	cloned gov.OwnerCloned,
	refType schema.RefType,
	from schema.Motion,
	to schema.Motion,
	fromPolicyNS ns.NS,
	toPolicyNS ns.NS,
) {

}

func (x proposalPolicy) RemoveRefFrom(
	ctx context.Context,
	cloned gov.OwnerCloned,
	refType schema.RefType,
	from schema.Motion,
	to schema.Motion,
	fromPolicyNS ns.NS,
	toPolicyNS ns.NS,
) {

}