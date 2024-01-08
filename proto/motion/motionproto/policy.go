package motionproto

import (
	"context"

	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/mod"
	"github.com/gov4git/gov4git/v2/proto/motion"
	"github.com/gov4git/gov4git/v2/proto/notice"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/ns"
)

type Report = any

type Policy interface {
	Open(
		ctx context.Context,
		cloned gov.OwnerCloned,
		motion Motion,
		instancePolicyNS ns.NS,
		args ...any,
	) (Report, notice.Notices)

	// Score is invoked only on open motions.
	Score(
		ctx context.Context,
		cloned gov.OwnerCloned,
		motion Motion,
		instancePolicyNS ns.NS,
		args ...any,
	) (Score, notice.Notices)

	// Update is invoked only on open motions, after rescoring all motions.
	Update(
		ctx context.Context,
		cloned gov.OwnerCloned,
		motion Motion,
		instancePolicyNS ns.NS,
		args ...any,
	) (Report, notice.Notices)

	Close(
		ctx context.Context,
		cloned gov.OwnerCloned,
		motion Motion,
		instancePolicyNS ns.NS,
		decision Decision,
		args ...any,
	) (Report, notice.Notices)

	Cancel(
		ctx context.Context,
		cloned gov.OwnerCloned,
		motion Motion,
		instancePolicyNS ns.NS,
		args ...any,
	) (Report, notice.Notices)

	Show(
		ctx context.Context,
		cloned gov.Cloned,
		motion Motion,
		instancePolicyNS ns.NS,
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
		fromPolicyNS ns.NS,
		toPolicyNS ns.NS,
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
		fromPolicyNS ns.NS,
		toPolicyNS ns.NS,
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
		fromPolicyNS ns.NS,
		toPolicyNS ns.NS,
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
		fromPolicyNS ns.NS,
		toPolicyNS ns.NS,
		args ...any,
	) (Report, notice.Notices)

	// Freeze is invoked by motion.Freeze
	Freeze(
		ctx context.Context,
		cloned gov.OwnerCloned,
		motion Motion,
		instancePolicyNS ns.NS,
		args ...any,
	) (Report, notice.Notices)

	// Unfreeze is invoked by motion.Unfreeze
	Unfreeze(
		ctx context.Context,
		cloned gov.OwnerCloned,
		motion Motion,
		instancePolicyNS ns.NS,
		args ...any,
	) (Report, notice.Notices)
}

type MotionBallot struct {
	Label         string               `json:"ballot_label"`
	BallotID      ballotproto.BallotID `json:"ballot_id"`
	BallotChoices []string             `json:"ballot_choices"`
	BallotAd      ballotproto.Ad       `json:"ballot_ad"`
	BallotTally   ballotproto.Tally    `json:"ballot_tally"`
}

type MotionBallots []MotionBallot

var policyRegistry = mod.NewModuleRegistry[motion.PolicyName, Policy]()

func Install(ctx context.Context, name motion.PolicyName, policy Policy) {
	policyRegistry.Set(ctx, name, policy)
}

func Get(ctx context.Context, key motion.PolicyName) Policy {
	return policyRegistry.Get(ctx, key)
}

func InstalledMotionPolicies() []string {
	return namesToStrings(policyRegistry.List())
}

func namesToStrings(ns []motion.PolicyName) []string {
	ss := make([]string, len(ns))
	for i := range ns {
		ss[i] = ns[i].String()
	}
	return ss
}

func GetMotionPolicy(ctx context.Context, m Motion) Policy {
	return Get(ctx, m.Policy)
}

// MotionPolicyNS returns the private policy namespace for a given motion instance.
func MotionPolicyNS(id MotionID) ns.NS {
	return MotionKV.KeyNS(MotionNS, id).Append("policy")
}
