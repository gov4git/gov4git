package proposal

import (
	"context"

	"github.com/gov4git/gov4git/proto/docket/policy"
	"github.com/gov4git/gov4git/proto/docket/schema"
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

func (x proposalPolicy) Open(ctx context.Context, instancePolicyNS ns.NS, motion schema.Motion) {

}

func (x proposalPolicy) Close(ctx context.Context, instancePolicyNS ns.NS, motion schema.Motion) {

}
