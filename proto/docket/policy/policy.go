package policy

import (
	"context"

	"github.com/gov4git/gov4git/proto/docket/schema"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/mod"
	"github.com/gov4git/lib4git/ns"
)

type Policy interface {
	Name() string

	Open(
		ctx context.Context,
		govAddr gov.OwnerAddress,
		govCloned gov.OwnerCloned,
		motion schema.Motion,
		instancePolicyNS ns.NS,
	)

	Score(
		ctx context.Context,
		govAddr gov.OwnerAddress,
		govCloned gov.OwnerCloned,
		motion schema.Motion,
		instancePolicyNS ns.NS,
	) schema.Score

	Close(
		ctx context.Context,
		govAddr gov.OwnerAddress,
		govCloned gov.OwnerCloned,
		motion schema.Motion,
		instancePolicyNS ns.NS,
	)
}

var policyRegistry = mod.NewModuleRegistry[Policy]()

func Install(ctx context.Context, policy Policy) {
	policyRegistry.Set(ctx, policy.Name(), policy)
}

func Get(ctx context.Context, key string) Policy {
	return policyRegistry.Get(ctx, key)
}

func InstalledPolicies() []string {
	return policyRegistry.Keys()
}

func GetMotionPolicy(ctx context.Context, m schema.Motion) Policy {
	return Get(ctx, m.Policy.String())
}

// MotionPolicyNS returns the private policy namespace for a given motion instance.
func MotionPolicyNS(id schema.MotionID) ns.NS {
	return schema.MotionKV.KeyNS(schema.MotionNS, id).Append("policy")
}
