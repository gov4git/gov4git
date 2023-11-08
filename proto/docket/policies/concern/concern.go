package concern

import (
	"context"

	"github.com/gov4git/gov4git/proto/docket/policy"
	"github.com/gov4git/gov4git/proto/docket/schema"
	"github.com/gov4git/lib4git/ns"
)

func init() {
	policy.Install(context.Background(), concernPolicy{})
}

const ConcernPolicyName = schema.PolicyName("concern-policy")

type concernPolicy struct{}

func (x concernPolicy) Name() string {
	return ConcernPolicyName.String()
}

func (x concernPolicy) Open(ctx context.Context, instancePolicyNS ns.NS, motion schema.Motion) {

}

func (x concernPolicy) Close(ctx context.Context, instancePolicyNS ns.NS, motion schema.Motion) {

}
