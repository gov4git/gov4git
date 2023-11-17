package concern

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/proto/ballot/ballot"
	"github.com/gov4git/gov4git/proto/ballot/qv"
	"github.com/gov4git/gov4git/proto/docket/policies/pmp"
	"github.com/gov4git/gov4git/proto/docket/policy"
	"github.com/gov4git/gov4git/proto/docket/schema"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/member"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/ns"
)

func init() {
	policy.Install(context.Background(), concernPolicy{})
}

const ConcernPolicyName = schema.PolicyName("pmp-concern-policy")

type concernPolicy struct{}

func (x concernPolicy) Name() string {
	return ConcernPolicyName.String()
}

func (x concernPolicy) Open(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion schema.Motion,
	policyNS ns.NS,

) {

	// initialize state
	state := NewConcernState(motion.ID)
	SaveState_StageOnly(ctx, cloned.Public.Tree(), policyNS, state)

	// open a poll for the motion
	ballot.Open_StageOnly(
		ctx,
		qv.QV{},
		cloned,
		state.PriorityPoll,
		fmt.Sprintf("Prioritization poll for motion %v", motion.ID),
		fmt.Sprintf("Up/down vote the priority for concern (issue) %v", motion.ID),
		[]string{schema.MotionPollBallotChoice},
		member.Everybody,
	)
}

func (x concernPolicy) Score(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion schema.Motion,
	policyNS ns.NS,

) schema.Score {

	state := LoadState_Local(ctx, cloned.Public.Tree(), policyNS)

	// compute motion score from the priority poll ballot
	ads := ballot.Show_Local(ctx, cloned.Public.Tree(), state.PriorityPoll)
	attention := ads.Tally.Attention()

	return schema.Score{
		Attention: attention,
	}
}

func (x concernPolicy) Close(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion schema.Motion,
	policyNS ns.NS,

) {

	// close the poll for the motion
	priorityPollName := pmp.ConcernPollBallotName(motion.ID)
	ballot.Close_StageOnly(
		ctx,
		cloned,
		priorityPollName,
		false,
	)

}

func (x concernPolicy) Cancel(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion schema.Motion,
	policyNS ns.NS,

) {

	// cancel the poll for the motion (returning credits to users)
	priorityPollName := pmp.ConcernPollBallotName(motion.ID)
	ballot.Close_StageOnly(
		ctx,
		cloned,
		priorityPollName,
		true,
	)

}

func (x concernPolicy) Show(
	ctx context.Context,
	cloned gov.Cloned,
	motion schema.Motion,
	policyNS ns.NS,

) form.Map {

	// retrieve policy state
	policyState := LoadState_Local(ctx, cloned.Tree(), policyNS)

	// retrieve poll state
	priorityPollName := pmp.ConcernPollBallotName(motion.ID)
	pollState := ballot.Show_Local(ctx, cloned.Tree(), priorityPollName)

	return form.Map{
		"pmp_concern_policy_state":        policyState,
		"pmp_concern_priority_poll_state": pollState,
	}
}

func (x concernPolicy) AddRefTo(
	ctx context.Context,
	cloned gov.OwnerCloned,
	refType schema.RefType,
	from schema.Motion,
	to schema.Motion,
	fromPolicyNS ns.NS,
	toPolicyNS ns.NS,

) {
}

func (x concernPolicy) AddRefFrom(
	ctx context.Context,
	cloned gov.OwnerCloned,
	refType schema.RefType,
	from schema.Motion,
	to schema.Motion,
	fromPolicyNS ns.NS,
	toPolicyNS ns.NS,

) {

	// freeze poll if from-motion is a highly-ranked pr
	if !from.IsProposal() {
		return
	}
	if from.Score.Attention <= 0.0 { // XXX: add a global threshold for considering proposals as eligible solutions
		return
	}

	toState := LoadState_Local(ctx, cloned.Public.Tree(), toPolicyNS)
	ref := schema.Ref{Type: refType, From: from.ID, To: to.ID}

	// check if reference has already triggered a freeze
	if toState.EligibleProposals.Contains(ref) {
		return
	}

	// if not, record the reference as a trigger
	toState.EligibleProposals = append(toState.EligibleProposals, ref)
	SaveState_StageOnly(ctx, cloned.Public.Tree(), toPolicyNS, toState)

	// freeze priority poll, if not already frozen
	priorityPoll := pmp.ConcernPollBallotName(to.ID)
	if ballot.IsFrozen_Local(ctx, cloned.PublicClone(), priorityPoll) {
		return
	}
	ballot.Freeze_StageOnly(ctx, cloned, priorityPoll)
}

func (x concernPolicy) RemoveRefTo(
	ctx context.Context,
	cloned gov.OwnerCloned,
	refType schema.RefType,
	from schema.Motion,
	to schema.Motion,
	fromPolicyNS ns.NS,
	toPolicyNS ns.NS,
) {
}

func (x concernPolicy) RemoveRefFrom(
	ctx context.Context,
	cloned gov.OwnerCloned,
	refType schema.RefType,
	from schema.Motion,
	to schema.Motion,
	fromPolicyNS ns.NS,
	toPolicyNS ns.NS,

) {

	// unfreeze poll if from-motion had triggered a freeze, and there are no other triggers remaining

	toState := LoadState_Local(ctx, cloned.Public.Tree(), toPolicyNS)
	ref := schema.Ref{Type: refType, From: from.ID, To: to.ID}

	// check if reference is the only remaining trigger
	if toState.EligibleProposals.Len() != 1 || !schema.RefEqual(ref, toState.EligibleProposals[0]) {

		// remove it as the last remaining trigger
		toState.EligibleProposals = nil
		SaveState_StageOnly(ctx, cloned.Public.Tree(), toPolicyNS, toState)

		// unfreeze the priority poll ballot, if frozen
		priorityPoll := pmp.ConcernPollBallotName(to.ID)
		if !ballot.IsFrozen_Local(ctx, cloned.PublicClone(), priorityPoll) {
			return
		}
		ballot.Unfreeze_StageOnly(ctx, cloned, priorityPoll)

	} else if toState.EligibleProposals.Contains(ref) {

		// remove it from the list of triggers
		toState.EligibleProposals = toState.EligibleProposals.Remove(ref)
		SaveState_StageOnly(ctx, cloned.Public.Tree(), toPolicyNS, toState)
	}

}
