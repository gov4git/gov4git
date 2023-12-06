package policy

import (
	"context"

	"github.com/gov4git/gov4git/proto/docket/schema"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/mod"
	"github.com/gov4git/gov4git/proto/notice"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/ns"
)

type Policy interface {
	Open(
		ctx context.Context,
		cloned gov.OwnerCloned,
		motion schema.Motion,
		instancePolicyNS ns.NS,
		args ...any,
	) notice.Notices

	// Score is invoked only on open motions.
	Score(
		ctx context.Context,
		cloned gov.OwnerCloned,
		motion schema.Motion,
		instancePolicyNS ns.NS,
		args ...any,
	) (schema.Score, notice.Notices)

	// Update is invoked only on open motions, after rescoring all motions.
	Update(
		ctx context.Context,
		cloned gov.OwnerCloned,
		motion schema.Motion,
		instancePolicyNS ns.NS,
		args ...any,
	) notice.Notices

	Close(
		ctx context.Context,
		cloned gov.OwnerCloned,
		motion schema.Motion,
		instancePolicyNS ns.NS,
		args ...any,
	) notice.Notices

	Cancel(
		ctx context.Context,
		cloned gov.OwnerCloned,
		motion schema.Motion,
		instancePolicyNS ns.NS,
		args ...any,
	) notice.Notices

	Show(
		ctx context.Context,
		cloned gov.Cloned,
		motion schema.Motion,
		instancePolicyNS ns.NS,
		args ...any,
	) form.Form

	// AddRefTo is invoked only when to and from motions are open.
	AddRefTo(
		ctx context.Context,
		cloned gov.OwnerCloned,
		refType schema.RefType,
		from schema.Motion,
		to schema.Motion,
		fromPolicyNS ns.NS,
		toPolicyNS ns.NS,
		args ...any,
	) notice.Notices

	// AddRefFrom is invoked only when to and from motions are open.
	AddRefFrom(
		ctx context.Context,
		cloned gov.OwnerCloned,
		refType schema.RefType,
		from schema.Motion,
		to schema.Motion,
		fromPolicyNS ns.NS,
		toPolicyNS ns.NS,
		args ...any,
	) notice.Notices

	// RemoveRefTo is invoked only when to and from motions are open.
	RemoveRefTo(
		ctx context.Context,
		cloned gov.OwnerCloned,
		refType schema.RefType,
		from schema.Motion,
		to schema.Motion,
		fromPolicyNS ns.NS,
		toPolicyNS ns.NS,
		args ...any,
	) notice.Notices

	// RemoveRefFrom is invoked only when to and from motions are open.
	RemoveRefFrom(
		ctx context.Context,
		cloned gov.OwnerCloned,
		refType schema.RefType,
		from schema.Motion,
		to schema.Motion,
		fromPolicyNS ns.NS,
		toPolicyNS ns.NS,
		args ...any,
	) notice.Notices

	// Freeze is invoked by motion.Freeze
	Freeze(
		ctx context.Context,
		cloned gov.OwnerCloned,
		motion schema.Motion,
		instancePolicyNS ns.NS,
		args ...any,
	) notice.Notices

	// Unfreeze is invoked by motion.Unfreeze
	Unfreeze(
		ctx context.Context,
		cloned gov.OwnerCloned,
		motion schema.Motion,
		instancePolicyNS ns.NS,
		args ...any,
	) notice.Notices
}

var policyRegistry = mod.NewModuleRegistry[schema.PolicyName, Policy]()

func Install(ctx context.Context, name schema.PolicyName, policy Policy) {
	policyRegistry.Set(ctx, name, policy)
}

func Get(ctx context.Context, key schema.PolicyName) Policy {
	return policyRegistry.Get(ctx, key)
}

func InstalledMotionPolicies() []string {
	return namesToStrings(policyRegistry.List())
}

func namesToStrings(ns []schema.PolicyName) []string {
	ss := make([]string, len(ns))
	for i := range ns {
		ss[i] = ns[i].String()
	}
	return ss
}

func GetMotionPolicy(ctx context.Context, m schema.Motion) Policy {
	return Get(ctx, m.Policy)
}

// MotionPolicyNS returns the private policy namespace for a given motion instance.
func MotionPolicyNS(id schema.MotionID) ns.NS {
	return schema.MotionKV.KeyNS(schema.MotionNS, id).Append("policy")
}
