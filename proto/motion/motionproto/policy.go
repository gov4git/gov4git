package motionproto

import (
	"context"

	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/mod"
	"github.com/gov4git/gov4git/v2/proto/motion"
	"github.com/gov4git/gov4git/v2/proto/notice"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/must"
	"github.com/gov4git/lib4git/ns"
)

type Report = any

type PolicyDescriptor struct {
	Description       string `json:"description"`         // markdown description
	GithubLabel       string `json:"github_label"`        // label to apply on github to activate policy
	AppliesToConcern  bool   `json:"applies_to_concern"`  // can be applied to concerns
	AppliesToProposal bool   `json:"applies_to_proposal"` // can be applied to proposals
}

type Policy interface {

	// Descriptor returns a markdown description of the policy.
	Descriptor() PolicyDescriptor

	gov.PostCloner

	// pipeline

	// Update is invoked, on all motions that are not closed.
	Update(
		ctx context.Context,
		cloned gov.OwnerCloned,
		motion Motion,
		args ...any,
	) (Report, notice.Notices)

	// Aggregate is invoked after Update, over all motions that are not closed.
	Aggregate(
		ctx context.Context,
		cloned gov.OwnerCloned,
		motion Motions,
	)

	// Score is invoked after Aggregate, on all motions that are not closed.
	Score(
		ctx context.Context,
		cloned gov.OwnerCloned,
		motion Motion,
		args ...any,
	) (Score, notice.Notices)

	// Clear is invoked after Score, on all motions that are not archived.
	Clear(
		ctx context.Context,
		cloned gov.OwnerCloned,
		motion Motion,
		args ...any,
	) (Report, notice.Notices)

	// operations

	Open(
		ctx context.Context,
		cloned gov.OwnerCloned,
		motion Motion,
		args ...any,
	) (Report, notice.Notices)

	Close(
		ctx context.Context,
		cloned gov.OwnerCloned,
		motion Motion,
		decision Decision,
		args ...any,
	) (Report, notice.Notices)

	Cancel(
		ctx context.Context,
		cloned gov.OwnerCloned,
		motion Motion,
		args ...any,
	) (Report, notice.Notices)

	Show(
		ctx context.Context,
		cloned gov.Cloned,
		motion Motion,
		args ...any,
	) (form.Form, MotionBallots)

	// AddRefTo is invoked when a reference to this motion is added.
	// AddRefTo is invoked only when to and from motions are open.
	AddRefTo(
		ctx context.Context,
		cloned gov.OwnerCloned,
		refType RefType,
		from Motion,
		to Motion,
		args ...any,
	) (Report, notice.Notices)

	// AddRefFrom is invoked when a reference from this motion is added.
	// AddRefFrom is invoked only when to and from motions are open.
	AddRefFrom(
		ctx context.Context,
		cloned gov.OwnerCloned,
		refType RefType,
		from Motion,
		to Motion,
		args ...any,
	) (Report, notice.Notices)

	// RemoveRefTo is invoked when a reference to this motion is removed.
	// RemoveRefTo is invoked only when to and from motions are open.
	RemoveRefTo(
		ctx context.Context,
		cloned gov.OwnerCloned,
		refType RefType,
		from Motion,
		to Motion,
		args ...any,
	) (Report, notice.Notices)

	// RemoveRefFrom is invoked when a reference from this motion is removed.
	// RemoveRefFrom is invoked only when to and from motions are open.
	RemoveRefFrom(
		ctx context.Context,
		cloned gov.OwnerCloned,
		refType RefType,
		from Motion,
		to Motion,
		args ...any,
	) (Report, notice.Notices)

	// Freeze is invoked by motion.Freeze
	Freeze(
		ctx context.Context,
		cloned gov.OwnerCloned,
		motion Motion,
		args ...any,
	) (Report, notice.Notices)

	// Unfreeze is invoked by motion.Unfreeze
	Unfreeze(
		ctx context.Context,
		cloned gov.OwnerCloned,
		motion Motion,
		args ...any,
	) (Report, notice.Notices)
}

type MotionBallot struct {
	Label         string               `json:"ballot_label"`
	BallotID      ballotproto.BallotID `json:"ballot_id"`
	BallotChoices []string             `json:"ballot_choices"`
	BallotAd      ballotproto.Ad       `json:"ballot_ad"`
	BallotTally   ballotproto.Tally    `json:"ballot_tally"`
	BallotMargin  *ballotproto.Margin  `json:"ballot_margin,omitempty"`
}

type MotionBallots []MotionBallot

var policyRegistry = mod.NewModuleRegistry[motion.PolicyName, Policy]()

func Install(ctx context.Context, name motion.PolicyName, policy Policy) {
	policyRegistry.Set(ctx, name, policy)
	gov.InstallPostClone(ctx, "motion-policy-"+string(name), policy)
}

func TryGetPolicy(ctx context.Context, key motion.PolicyName) Policy {
	p, _ := must.Try1[Policy](
		func() Policy {
			return policyRegistry.Get(ctx, key)
		},
	)
	return p
}

func GetPolicy(ctx context.Context, key motion.PolicyName) Policy {
	p, err := must.Try1[Policy](
		func() Policy {
			return policyRegistry.Get(ctx, key)
		},
	)
	must.Assertf(ctx, err == nil, "motion policy not supported") // ERR
	return p
}

func InstalledPolicyKeys() []string {
	return namesToStrings(policyRegistry.ListKeys())
}

func InstalledPolicyDescriptors() map[string]PolicyDescriptor {
	x := map[string]PolicyDescriptor{}
	keys, policies := policyRegistry.List()
	for i, k := range keys {
		x[k.String()] = policies[i].Descriptor()
	}
	return x
}

func namesToStrings(ns []motion.PolicyName) []string {
	ss := make([]string, len(ns))
	for i := range ns {
		ss[i] = ns[i].String()
	}
	return ss
}

func GetMotionPolicy(ctx context.Context, m Motion) Policy {
	return GetMotionPolicyByName(ctx, m.Policy)
}

func GetMotionPolicyByName(ctx context.Context, pn motion.PolicyName) Policy {
	return GetPolicy(ctx, pn)
}

// MotionPolicyNS returns the private policy namespace for a given motion instance.
func MotionPolicyNS(id MotionID) ns.NS {
	return MotionKV.KeyNS(MotionNS, id).Append("policy")
}
