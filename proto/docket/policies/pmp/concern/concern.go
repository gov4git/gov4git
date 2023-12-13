package concern

import (
	"bytes"
	"context"
	"fmt"
	"slices"

	"github.com/gov4git/gov4git/proto/account"
	"github.com/gov4git/gov4git/proto/ballot/ballot"
	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/ballot/load"
	"github.com/gov4git/gov4git/proto/docket/ops"
	"github.com/gov4git/gov4git/proto/docket/policies/pmp"
	"github.com/gov4git/gov4git/proto/docket/policy"
	"github.com/gov4git/gov4git/proto/docket/schema"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/member"
	"github.com/gov4git/gov4git/proto/notice"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/must"
	"github.com/gov4git/lib4git/ns"
)

func init() {
	policy.Install(context.Background(), ConcernPolicyName, concernPolicy{})
}

const ConcernPolicyName = schema.PolicyName("pmp-concern-policy")

type concernPolicy struct{}

func (x concernPolicy) Open(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion schema.Motion,
	policyNS ns.NS,
	args ...any,

) (policy.Report, notice.Notices) {

	// initialize state
	state := NewConcernState(motion.ID)
	SaveState_StageOnly(ctx, cloned.Public.Tree(), policyNS, state)

	// open a poll for the motion
	ballot.Open_StageOnly(
		ctx,
		load.QVStrategyName,
		cloned,
		state.PriorityPoll,
		fmt.Sprintf("Prioritization poll for motion %v", motion.ID),
		fmt.Sprintf("Up/down vote the priority for concern (issue) %v", motion.ID),
		[]string{pmp.ConcernBallotChoice},
		member.Everybody,
	)

	return nil, notice.Noticef(ctx, "Started managing this issue as Gov4Git concern `%v` with initial __priority score__ of `%v`."+
		pmp.Welcome, motion.ID, state.LatestPriorityScore)
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
	}, nil
}

func (x concernPolicy) Update(
	ctx context.Context,
	cloned gov.OwnerCloned,
	con schema.Motion,
	policyNS ns.NS,
	args ...any,

) (policy.Report, notice.Notices) {

	notices := notice.Notices{}
	state := LoadState_Local(ctx, cloned.Public.Tree(), policyNS)

	// update priority score

	ads := ballot.Show_Local(ctx, cloned.Public.Tree(), state.PriorityPoll)
	latestPriorityScore := ads.Tally.Scores[pmp.ConcernBallotChoice]
	if latestPriorityScore != state.LatestPriorityScore {
		notices = append(
			notices,
			notice.Noticef(ctx, "This issue's __priority score__ was updated to `%0.6f`.", latestPriorityScore)...,
		)
	}
	state.LatestPriorityScore = latestPriorityScore

	// update eligible proposals

	eligible := schema.Refs{}
	for _, ref := range con.RefBy {
		if pmp.IsConcernProposalEligible(ctx, cloned.PublicClone(), con.ID, ref.From, ref.Type) {
			eligible = append(eligible, ref)
		}
	}
	eligible.Sort()
	if !slices.Equal[schema.Refs](eligible, state.EligibleProposals) {
		// display updated list of eligible proposals
		var w bytes.Buffer
		for _, ref := range eligible {
			propMot := ops.LookupMotion_Local(ctx, cloned.PublicClone(), ref.From)
			fmt.Fprintf(&w, "- %s, managed as Gov4Git motion `%v`\n", propMot.TrackerURL, propMot.ID)
		}
		notices = append(
			notices,
			notice.Noticef(ctx, "The set of eligible proposals addressing this issue changed:\n"+w.String())...,
		)
	}
	state.EligibleProposals = eligible

	//

	SaveState_StageOnly(ctx, cloned.Public.Tree(), policyNS, state)

	r0, n0 := x.updateFreeze(ctx, cloned, con, policyNS)
	return r0, append(notices, n0...)
}

func (x concernPolicy) updateFreeze(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion schema.Motion,
	policyNS ns.NS,
	args ...any,

) (policy.Report, notice.Notices) {

	toState := LoadState_Local(ctx, cloned.Public.Tree(), policyNS)

	notices := notice.Notices{}
	if toState.EligibleProposals.Len() > 0 && !motion.Frozen {
		ops.FreezeMotion_StageOnly(notice.Mute(ctx), cloned, motion.ID)

		var w bytes.Buffer
		fmt.Fprintf(&w, "Freezing ❄️ this issue as there are eligible PRs addressing it:\n")
		for _, pr := range toState.EligibleProposals {
			pr := ops.LookupMotion_Local(ctx, cloned.PublicClone(), pr.From)
			fmt.Fprintf(&w, "- %s\n", pr.TrackerURL)
		}
		notices = append(notices, notice.Noticef(ctx, w.String())...)
	}
	if toState.EligibleProposals.Len() == 0 && motion.Frozen {
		ops.UnfreezeMotion_StageOnly(notice.Mute(ctx), cloned, motion.ID)
		notices = append(notices, notice.Noticef(ctx, "Unfreezing 🌤️ issue as there are no eligible PRs are addressing it.")...)
	}

	return nil, notices
}

func (x concernPolicy) Close(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion schema.Motion,
	policyNS ns.NS,
	decision schema.Decision,
	args ...any,
	// args[0]=toID account.AccountID
	// args[1]=prop schema.Motion

) (policy.Report, notice.Notices) {

	must.Assertf(ctx, len(args) == 2, "issue closure requires two arguments, got %v", args)
	toID, ok := args[0].(account.AccountID)
	must.Assertf(ctx, ok, "unrecognized account ID argument %v", args[0])
	prop, ok := args[1].(schema.Motion)
	must.Assertf(ctx, ok, "unrecognized proposal motion argument %v", args[1])

	// close the poll for the motion
	priorityPollName := pmp.ConcernPollBallotName(motion.ID)
	chg := ballot.Close_StageOnly(
		ctx,
		cloned,
		priorityPollName,
		toID,
	)

	return &CloseReport{}, closeNotice(ctx, motion, chg.Result, prop)
}

func (x concernPolicy) Cancel(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion schema.Motion,
	policyNS ns.NS,
	args ...any,

) (policy.Report, notice.Notices) {

	// cancel the poll for the motion (returning credits to users)
	priorityPollName := pmp.ConcernPollBallotName(motion.ID)
	chg := ballot.Cancel_StageOnly(
		ctx,
		cloned,
		priorityPollName,
	)

	return &CancelReport{
		PriorityPollOutcome: chg.Result,
	}, cancelNotice(ctx, motion, chg.Result)
}

type PolicyView struct {
	State        *ConcernState  `json:"state"`
	PriorityPoll common.AdTally `json:"priority_poll"`
}

func (x concernPolicy) Show(
	ctx context.Context,
	cloned gov.Cloned,
	motion schema.Motion,
	policyNS ns.NS,
	args ...any,

) form.Form {

	// retrieve policy state
	policyState := LoadState_Local(ctx, cloned.Tree(), policyNS)

	// retrieve poll state
	priorityPollName := pmp.ConcernPollBallotName(motion.ID)
	pollState := ballot.Show_Local(ctx, cloned.Tree(), priorityPollName)

	return PolicyView{
		State:        policyState,
		PriorityPoll: pollState,
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

) (policy.Report, notice.Notices) {

	if !from.IsProposal() {
		return nil, nil
	}

	if refType != pmp.ResolvesRefType {
		return nil, nil
	}

	return nil, notice.Noticef(ctx, "This issue was referenced by %v, managed as Gov4Git proposal `%v`.", from.TrackerURL, from.ID)
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

) (policy.Report, notice.Notices) {

	return nil, nil
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

) (policy.Report, notice.Notices) {

	if !from.IsProposal() {
		return nil, nil
	}

	if refType != pmp.ResolvesRefType {
		return nil, nil
	}

	return nil, notice.Noticef(ctx, "This issue is no longer referenced by %v, managed as Gov4Git proposal `%v`.", from.TrackerURL, from.ID)
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

) (policy.Report, notice.Notices) {

	return nil, nil
}

func (x concernPolicy) Freeze(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion schema.Motion,
	policyNS ns.NS,
	args ...any,

) (policy.Report, notice.Notices) {

	// freeze priority poll, if not already frozen
	priorityPoll := pmp.ConcernPollBallotName(motion.ID)
	if ballot.IsFrozen_Local(ctx, cloned.PublicClone(), priorityPoll) {
		return nil, nil
	}
	ballot.Freeze_StageOnly(ctx, cloned, priorityPoll)

	return nil, notice.Noticef(ctx, "This issue, managed by Gov4Git concern `%v`, has been frozen ❄️", motion.ID)
}

func (x concernPolicy) Unfreeze(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion schema.Motion,
	policyNS ns.NS,
	args ...any,

) (policy.Report, notice.Notices) {

	// unfreeze the priority poll ballot, if frozen
	priorityPoll := pmp.ConcernPollBallotName(motion.ID)
	if !ballot.IsFrozen_Local(ctx, cloned.PublicClone(), priorityPoll) {
		return nil, nil
	}
	ballot.Unfreeze_StageOnly(ctx, cloned, priorityPoll)

	return nil, notice.Noticef(ctx, "This issue, managed by Gov4Git concern `%v`, has been unfrozen 🌤️", motion.ID)
}

// motion.Un/Freeze --calls--> policy Un/Freeze --calls--> ballot Un/Freeze
