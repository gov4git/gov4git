package policy

import (
	"context"

	"github.com/gov4git/gov4git/proto/docket/schema"
	"github.com/gov4git/gov4git/proto/mod"
	"github.com/gov4git/lib4git/ns"
)

type Policy interface {
	Name() string
	Open(ctx context.Context, instancePolicyNS ns.NS, motion schema.Motion)
	Close(ctx context.Context, instancePolicyNS ns.NS, motion schema.Motion)
}

var policyRegistry = mod.NewModuleRegistry[Policy]()

func Install(ctx context.Context, policy Policy) {
	policyRegistry.Set(ctx, policy.Name(), policy)
}

func Get(ctx context.Context, key string) Policy {
	return policyRegistry.Get(ctx, key)
}

func Policies() []string {
	return policyRegistry.Keys()
}
