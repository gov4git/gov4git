package proposal

import (
	"context"

	"github.com/gov4git/gov4git/proto/docket/policy"
	"github.com/gov4git/gov4git/proto/docket/schema"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/lib4git/ns"
)

func init() {
	policy.Install(context.Background(), proposalPolicy{})
}

const ProposalPolicyName = schema.PolicyName("proposal-policy")

type proposalPolicy struct{}

func (x proposalPolicy) Name() string {
	return ProposalPolicyName.String()
}

func (x proposalPolicy) Open(
	ctx context.Context,
	govAddr gov.GovOwnerAddress,
	govCloned gov.GovOwnerCloned,
	motion schema.Motion,
	instancePolicyNS ns.NS,

) {

}

func (x proposalPolicy) Score(
	ctx context.Context,
	govAddr gov.GovOwnerAddress,
	govCloned gov.GovOwnerCloned,
	motion schema.Motion,
	instancePolicyNS ns.NS,

) schema.Score {

	return schema.Score{}
}

func (x proposalPolicy) Close(
	ctx context.Context,
	govAddr gov.GovOwnerAddress,
	govCloned gov.GovOwnerCloned,
	motion schema.Motion,
	instancePolicyNS ns.NS,

) {

}
