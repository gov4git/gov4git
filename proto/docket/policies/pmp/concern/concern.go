package concern

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/proto/account"
	"github.com/gov4git/gov4git/proto/ballot/ballot"
	"github.com/gov4git/gov4git/proto/ballot/qv"
	"github.com/gov4git/gov4git/proto/docket/ops"
	"github.com/gov4git/gov4git/proto/docket/policies/pmp"
	"github.com/gov4git/gov4git/proto/docket/policy"
	"github.com/gov4git/gov4git/proto/docket/schema"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/member"
	"github.com/gov4git/gov4git/proto/notice"
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
	args ...any,

) notice.Notices {

	// initialize state
	state := NewConcernState(motion.ID)
	SaveState_StageOnly(ctx, cloned.Public.Tree(), policyNS, state)

	// create an account for the concern
	account.Create_StageOnly(
		ctx,
		cloned.PublicClone(),
		schema.PolicyAccountID(motion.ID),
		schema.MotionOwnerID(motion.ID),
	)

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

	return notice.Noticef("Started managing this issue as Gov4Git concern #%v.", motion.ID)
}

func (x concernPolicy) Score(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion schema.Motion,
	policyNS ns.NS,
	args ...any,

) (schema.Score, notice.Notices) {

	state := LoadState_Local(ctx, cloned.Public.Tree(), policyNS)

	// compute motion score from the priority poll ballot
	ads := ballot.Show_Local(ctx, cloned.Public.Tree(), state.PriorityPoll)
	attention := ads.Tally.Attention()

	return schema.Score{
		Attention: attention,
	}, notice.Noticef("Updated prioritization tally to %v.", ads.Tally.Scores[schema.MotionPollBallotChoice])
}

func (x concernPolicy) Update(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion schema.Motion,
	policyNS ns.NS,
	args ...any,

) notice.Notices {

	// update eligible proposals

	state := LoadState_Local(ctx, cloned.Public.Tree(), policyNS)

	eligible := schema.Refs{}
	for ref := range state.EligibleProposals.RefSet() {
		if IsProposalEligible(ctx, cloned.PublicClone(), ref.From) {
			eligible = append(eligible, ref)
		}
	}
	eligible.Sort()
	state.EligibleProposals = eligible

	SaveState_StageOnly(ctx, cloned.Public.Tree(), policyNS, state)

	return x.updateFreeze(ctx, cloned, motion, policyNS)
}

func (x concernPolicy) Close(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion schema.Motion,
	policyNS ns.NS,
	args ...any,

) notice.Notices {

	// close the poll for the motion
	priorityPollName := pmp.ConcernPollBallotName(motion.ID)
	ballot.Close_StageOnly(
		ctx,
		cloned,
		priorityPollName,
		false,
	)

	return notice.Noticef("Closing managment of this issue, managed as Gov4Git concern #%v.", motion.ID)
}

func (x concernPolicy) Cancel(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion schema.Motion,
	policyNS ns.NS,
	args ...any,

) notice.Notices {

	// cancel the poll for the motion (returning credits to users)
	priorityPollName := pmp.ConcernPollBallotName(motion.ID)
	ballot.Close_StageOnly(
		ctx,
		cloned,
		priorityPollName,
		true,
	)

	return notice.Noticef("Cancelling management of this issue, managed as Gov4Git concern #%v.", motion.ID)
}

func (x concernPolicy) Show(
	ctx context.Context,
	cloned gov.Cloned,
	motion schema.Motion,
	policyNS ns.NS,
	args ...any,

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
	args ...any,

) notice.Notices {

	return nil
}

func (x concernPolicy) AddRefFrom(
	ctx context.Context,
	cloned gov.OwnerCloned,
	refType schema.RefType,
	from schema.Motion,
	to schema.Motion,
	fromPolicyNS ns.NS,
	toPolicyNS ns.NS,
	args ...any,

) notice.Notices {

	if refType != AddressesRefType {
		return nil
	}

	if !IsProposalEligible(ctx, cloned.PublicClone(), from.ID) {
		return nil
	}

	toState := LoadState_Local(ctx, cloned.Public.Tree(), toPolicyNS)
	ref := schema.Ref{Type: refType, From: from.ID, To: to.ID}

	if toState.EligibleProposals.Contains(ref) {
		return nil
	}

	toState.EligibleProposals = append(toState.EligibleProposals, ref)
	SaveState_StageOnly(ctx, cloned.Public.Tree(), toPolicyNS, toState)

	return notice.Noticef("This issue has been referenced by an eligible PR, managed as Gov4Git proposal #%v.", from.ID)
}

func (x concernPolicy) RemoveRefTo(
	ctx context.Context,
	cloned gov.OwnerCloned,
	refType schema.RefType,
	from schema.Motion,
	to schema.Motion,
	fromPolicyNS ns.NS,
	toPolicyNS ns.NS,
	args ...any,

) notice.Notices {

	return nil
}

func (x concernPolicy) RemoveRefFrom(
	ctx context.Context,
	cloned gov.OwnerCloned,
	refType schema.RefType,
	from schema.Motion,
	to schema.Motion,
	fromPolicyNS ns.NS,
	toPolicyNS ns.NS,
	args ...any,

) notice.Notices {

	toState := LoadState_Local(ctx, cloned.Public.Tree(), toPolicyNS)
	ref := schema.Ref{Type: refType, From: from.ID, To: to.ID}

	if !toState.EligibleProposals.Contains(ref) {
		return nil
	}

	toState.EligibleProposals = toState.EligibleProposals.Remove(ref)
	SaveState_StageOnly(ctx, cloned.Public.Tree(), toPolicyNS, toState)

	return notice.Noticef("This issue is no longer referenced by the PR, managed as Gov4Git concern #%v.", from.ID)
}

func (x concernPolicy) updateFreeze(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion schema.Motion,
	policyNS ns.NS,
	args ...any,

) notice.Notices {

	toState := LoadState_Local(ctx, cloned.Public.Tree(), policyNS)

	notices := notice.Notices{}
	if toState.EligibleProposals.Len() > 0 && !motion.Frozen {
		ops.FreezeMotion_StageOnly(ctx, cloned, motion.ID)
		notices = append(notices, notice.Noticef("Freezing issue as there are eligible PRs addressing it.")...)
	}
	if toState.EligibleProposals.Len() == 0 && motion.Frozen {
		ops.UnfreezeMotion_StageOnly(ctx, cloned, motion.ID)
		notices = append(notices, notice.Noticef("Unfreezing issue as there are no eligible PRs are addressing it.")...)
	}

	return notices
}

func (x concernPolicy) Freeze(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion schema.Motion,
	policyNS ns.NS,
	args ...any,

) notice.Notices {

	// freeze priority poll, if not already frozen
	priorityPoll := pmp.ConcernPollBallotName(motion.ID)
	if ballot.IsFrozen_Local(ctx, cloned.PublicClone(), priorityPoll) {
		return nil
	}
	ballot.Freeze_StageOnly(ctx, cloned, priorityPoll)

	return notice.Noticef("This issue, managed by Gov4Git concern #%v, has been frozen.", motion.ID)
}

func (x concernPolicy) Unfreeze(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion schema.Motion,
	policyNS ns.NS,
	args ...any,

) notice.Notices {

	// unfreeze the priority poll ballot, if frozen
	priorityPoll := pmp.ConcernPollBallotName(motion.ID)
	if !ballot.IsFrozen_Local(ctx, cloned.PublicClone(), priorityPoll) {
		return nil
	}
	ballot.Unfreeze_StageOnly(ctx, cloned, priorityPoll)

	return notice.Noticef("This issue, managed by Gov4Git concern #%v, has been unfrozen.", motion.ID)
}

// motion.Un/Freeze --calls--> policy Un/Freeze --calls--> ballot Un/Freeze
