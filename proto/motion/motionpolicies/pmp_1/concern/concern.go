package concern

import (
	"bytes"
	"context"
	"fmt"
	"slices"

	"github.com/gov4git/gov4git/v2/proto/account"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotapi"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotio"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/history/metric"
	"github.com/gov4git/gov4git/v2/proto/member"
	"github.com/gov4git/gov4git/v2/proto/motion"
	"github.com/gov4git/gov4git/v2/proto/motion/motionapi"
	"github.com/gov4git/gov4git/v2/proto/motion/motionpolicies/pmp_0"
	"github.com/gov4git/gov4git/v2/proto/motion/motionpolicies/pmp_1"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
	"github.com/gov4git/gov4git/v2/proto/notice"
	"github.com/gov4git/gov4git/v2/proto/purpose"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/must"
	"github.com/gov4git/lib4git/ns"
)

func init() {
	motionproto.Install(context.Background(), ConcernPolicyName, concernPolicy{})
}

const ConcernPolicyName = motion.PolicyName("pmp-concern-policy-v1")

type concernPolicy struct{}

func (x concernPolicy) Open(
	ctx context.Context,
	cloned gov.OwnerCloned,
	con motionproto.Motion,
	policyNS ns.NS,
	args ...any,

) (motionproto.Report, notice.Notices) {

	// initialize state
	state := NewConcernState(con.ID)
	SaveState_StageOnly(ctx, cloned.Public.Tree(), policyNS, state)

	// open a poll for the motion
	ballotapi.Open_StageOnly(
		ctx,
		ballotio.QVPolicyName,
		cloned,
		state.PriorityPoll,
		pmp_1.ConcernAccountID(con.ID),
		purpose.Concern,
		con.Policy,
		fmt.Sprintf("Prioritization poll for motion %v", con.ID),
		fmt.Sprintf("Up/down vote the priority for concern (issue) %v", con.ID),
		[]string{pmp_1.ConcernBallotChoice},
		member.Everybody,
	)

	// metrics
	metric.Log_StageOnly(ctx, cloned.PublicClone(), &metric.Event{
		Motion: &metric.MotionEvent{
			Open: &metric.MotionOpen{
				ID:     metric.MotionID(con.ID),
				Type:   "concern-v1",
				Policy: metric.MotionPolicy(con.Policy),
			},
		},
	})

	return nil, notice.Noticef(ctx, "Started managing this issue as Gov4Git concern `%v` with initial __priority score__ of `%0.6f`."+
		pmp_1.Welcome, con.ID, state.PriorityScore)
}

func (x concernPolicy) Score(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion motionproto.Motion,
	policyNS ns.NS,
	args ...any,

) (motionproto.Score, notice.Notices) {

	state := LoadState_Local(ctx, cloned.Public.Tree(), policyNS)
	return motionproto.Score{
		Attention: state.PriorityScore,
	}, nil
}

func (x concernPolicy) Update(
	ctx context.Context,
	cloned gov.OwnerCloned,
	con motionproto.Motion,
	policyNS ns.NS,
	args ...any,

) (motionproto.Report, notice.Notices) {

	// outputs
	notices := notice.Notices{}

	// inputs
	policyState := LoadPolicyState_Local(ctx, cloned)
	instanceState := LoadState_Local(ctx, cloned.Public.Tree(), policyNS)
	ads := ballotapi.Show_Local(ctx, cloned.Public.Tree(), instanceState.PriorityPoll)

	// update idealized quadratic funding deficit
	capFunding := capitalistFunding(&ads.Tally)
	iqFunding := idealizedQuadraticFunding(&ads.Tally)
	iqDeficit := iqFunding - capFunding
	instanceState.IQDeficit = iqDeficit

	// update priority score
	matchFunds := pmp_0.GetMatchFundBalance_Local(ctx, cloned.PublicClone())
	latestPriorityScore := capFunding + matchRatio(matchFunds, policyState.MatchDeficit)*iqDeficit
	if latestPriorityScore != instanceState.PriorityScore {
		notices = append(
			notices,
			notice.Noticef(ctx, "This issue's __priority score__ was updated to `%0.6f`.", latestPriorityScore)...,
		)
	}
	instanceState.PriorityScore = latestPriorityScore

	// update eligible proposals

	eligible := computeEligibleProposals(ctx, cloned.PublicClone(), con)
	if !slices.Equal[motionproto.Refs](eligible, instanceState.EligibleProposals) {
		// display updated list of eligible proposals
		if len(eligible) == 0 {
			notices = append(
				notices,
				notice.Noticef(ctx, "The set of eligible proposals addressing this issue is now empty.\n")...,
			)
		} else {
			var w bytes.Buffer
			for _, ref := range eligible {
				prop := motionapi.LookupMotion_Local(ctx, cloned.PublicClone(), ref.From)
				fmt.Fprintf(&w, "- %s, managed as Gov4Git motion `%v` with community attention of `%0.6f`\n",
					prop.TrackerURL, prop.ID, prop.Score.Attention)
			}
			notices = append(
				notices,
				notice.Noticef(ctx, "The set of eligible proposals addressing this issue changed to:\n"+w.String())...,
			)
		}
	}
	instanceState.EligibleProposals = eligible

	//

	SaveState_StageOnly(ctx, cloned.Public.Tree(), policyNS, instanceState)

	r0, n0 := x.updateFreeze(ctx, cloned, con, policyNS)
	return r0, append(notices, n0...)
}

func computeEligibleProposals(ctx context.Context, cloned gov.Cloned, con motionproto.Motion) motionproto.Refs {
	eligible := motionproto.Refs{}
	for _, ref := range con.RefBy {
		if pmp_1.IsConcernProposalEligible(ctx, cloned, con.ID, ref.From, ref.Type) {
			eligible = append(eligible, ref)
		}
	}
	eligible.Sort()
	return eligible
}

func (x concernPolicy) updateFreeze(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion motionproto.Motion,
	policyNS ns.NS,
	args ...any,

) (motionproto.Report, notice.Notices) {

	toState := LoadState_Local(ctx, cloned.Public.Tree(), policyNS)

	notices := notice.Notices{}
	if toState.EligibleProposals.Len() > 0 && !motion.Frozen {
		motionapi.FreezeMotion_StageOnly(notice.Mute(ctx), cloned, motion.ID)

		var w bytes.Buffer
		fmt.Fprintf(&w, "Freezing ❄️ this issue as there are eligible PRs addressing it:\n")
		for _, pr := range toState.EligibleProposals {
			pr := motionapi.LookupMotion_Local(ctx, cloned.PublicClone(), pr.From)
			fmt.Fprintf(&w, "- %s\n", pr.TrackerURL)
		}
		notices = append(notices, notice.Noticef(ctx, w.String())...)
	}
	if toState.EligibleProposals.Len() == 0 && motion.Frozen {
		motionapi.UnfreezeMotion_StageOnly(notice.Mute(ctx), cloned, motion.ID)
		notices = append(notices, notice.Noticef(ctx, "Unfreezing 🌤️ issue as there are no eligible PRs addressing it.")...)
	}

	return nil, notices
}

func (x concernPolicy) Aggregate(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion motionproto.Motions,
	instancePolicyNS []ns.NS,

) {

	// load all motion policy states
	concernPolicyStates := make([]*ConcernState, len(motion))
	for i := range motion {
		concernPolicyStates[i] = LoadState_Local(ctx, cloned.PublicClone().Tree(), instancePolicyNS[i])
	}

	// aggregate match deficit
	matchDeficit := 0.0
	for _, concernPolicyState := range concernPolicyStates {
		matchDeficit += concernPolicyState.IQDeficit
	}

	// update policy state
	policyState := LoadPolicyState_Local(ctx, cloned)
	policyState.MatchDeficit = matchDeficit
	SavePolicyState_StageOnly(ctx, cloned, policyState)
}

func (x concernPolicy) Close(
	ctx context.Context,
	cloned gov.OwnerCloned,
	con motionproto.Motion,
	policyNS ns.NS,
	decision motionproto.Decision,
	args ...any,
	// args[0]=toID account.AccountID
	// args[1]=prop schema.Motion

) (motionproto.Report, notice.Notices) {

	must.Assertf(ctx, len(args) == 2, "issue closure requires two arguments, got %v", args)
	toID, ok := args[0].(account.AccountID)
	must.Assertf(ctx, ok, "unrecognized account ID argument %v", args[0])
	prop, ok := args[1].(motionproto.Motion)
	must.Assertf(ctx, ok, "unrecognized proposal motion argument %v", args[1])

	// update the policy state before closing the motion
	x.Update(ctx, cloned, prop, policyNS)

	// close the poll for the motion
	priorityPollName := pmp_1.ConcernPollBallotName(con.ID)
	chg := ballotapi.Close_StageOnly(
		ctx,
		cloned,
		priorityPollName,
		toID,
	)

	// metrics
	metric.Log_StageOnly(ctx, cloned.PublicClone(), &metric.Event{
		Motion: &metric.MotionEvent{
			Close: &metric.MotionClose{
				ID:       metric.MotionID(con.ID),
				Type:     "concern-v1",
				Decision: decision.MetricDecision(),
				Policy:   metric.MotionPolicy(con.Policy),
				Receipts: nil, // rewards are accounted for by the proposal
			},
		},
	})

	return &CloseReport{}, closeNotice(ctx, con, chg.Result, prop)
}

func (x concernPolicy) Cancel(
	ctx context.Context,
	cloned gov.OwnerCloned,
	con motionproto.Motion,
	policyNS ns.NS,
	args ...any,

) (motionproto.Report, notice.Notices) {

	// cancel the poll for the motion (returning credits to users)
	priorityPollName := pmp_1.ConcernPollBallotName(con.ID)
	chg := ballotapi.Cancel_StageOnly(
		ctx,
		cloned,
		priorityPollName,
	)

	// metrics
	metric.Log_StageOnly(ctx, cloned.PublicClone(), &metric.Event{
		Motion: &metric.MotionEvent{
			Cancel: &metric.MotionCancel{
				ID:       metric.MotionID(con.ID),
				Type:     "concern-v1",
				Policy:   metric.MotionPolicy(con.Policy),
				Receipts: nil, // refunds are accounted for by the proposal
			},
		},
	})

	return &CancelReport{
		PriorityPollOutcome: chg.Result,
	}, cancelNotice(ctx, con, chg.Result)
}

type PolicyView struct {
	State          *ConcernState       `json:"state"`
	PriorityPoll   ballotproto.AdTally `json:"priority_poll"`
	PriorityMargin ballotproto.Margin  `json:"priority_margin"`
}

func (x concernPolicy) Show(
	ctx context.Context,
	cloned gov.Cloned,
	motion motionproto.Motion,
	policyNS ns.NS,
	args ...any,

) (form.Form, motionproto.MotionBallots) {

	// retrieve policy state
	policyState := LoadState_Local(ctx, cloned.Tree(), policyNS)

	// retrieve poll state
	priorityPollName := pmp_1.ConcernPollBallotName(motion.ID)
	pollState := ballotapi.Show_Local(ctx, cloned.Tree(), priorityPollName)

	return PolicyView{
			State:          policyState,
			PriorityPoll:   pollState,
			PriorityMargin: *ballotapi.GetMargin_Local(ctx, cloned, priorityPollName),
		}, motionproto.MotionBallots{
			motionproto.MotionBallot{
				Label:         "priority_poll",
				BallotID:      policyState.PriorityPoll,
				BallotChoices: pollState.Ad.Choices,
				BallotAd:      pollState.Ad,
				BallotTally:   pollState.Tally,
			},
		}
}

func (x concernPolicy) AddRefTo(
	ctx context.Context,
	cloned gov.OwnerCloned,
	refType motionproto.RefType,
	from motionproto.Motion,
	to motionproto.Motion,
	fromPolicyNS ns.NS,
	toPolicyNS ns.NS,
	args ...any,

) (motionproto.Report, notice.Notices) {

	if !from.IsProposal() {
		return nil, nil
	}

	if refType != pmp_1.ClaimsRefType {
		return nil, nil
	}

	return nil, notice.Noticef(ctx, "This issue was referenced by %v, managed as Gov4Git proposal `%v`.", from.TrackerURL, from.ID)
}

func (x concernPolicy) AddRefFrom(
	ctx context.Context,
	cloned gov.OwnerCloned,
	refType motionproto.RefType,
	from motionproto.Motion,
	to motionproto.Motion,
	fromPolicyNS ns.NS,
	toPolicyNS ns.NS,
	args ...any,

) (motionproto.Report, notice.Notices) {

	return nil, nil
}

func (x concernPolicy) RemoveRefTo(
	ctx context.Context,
	cloned gov.OwnerCloned,
	refType motionproto.RefType,
	from motionproto.Motion,
	to motionproto.Motion,
	fromPolicyNS ns.NS,
	toPolicyNS ns.NS,
	args ...any,

) (motionproto.Report, notice.Notices) {

	if !from.IsProposal() {
		return nil, nil
	}

	if refType != pmp_1.ClaimsRefType {
		return nil, nil
	}

	return nil, notice.Noticef(ctx, "This issue is no longer referenced by %v, managed as Gov4Git proposal `%v`.", from.TrackerURL, from.ID)
}

func (x concernPolicy) RemoveRefFrom(
	ctx context.Context,
	cloned gov.OwnerCloned,
	refType motionproto.RefType,
	from motionproto.Motion,
	to motionproto.Motion,
	fromPolicyNS ns.NS,
	toPolicyNS ns.NS,
	args ...any,

) (motionproto.Report, notice.Notices) {

	return nil, nil
}

func (x concernPolicy) Freeze(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion motionproto.Motion,
	policyNS ns.NS,
	args ...any,

) (motionproto.Report, notice.Notices) {

	// freeze priority poll, if not already frozen
	priorityPoll := pmp_1.ConcernPollBallotName(motion.ID)
	if ballotapi.IsFrozen_Local(ctx, cloned.PublicClone(), priorityPoll) {
		return nil, nil
	}
	ballotapi.Freeze_StageOnly(ctx, cloned, priorityPoll)

	return nil, notice.Noticef(ctx, "This issue, managed by Gov4Git concern `%v`, has been frozen ❄️", motion.ID)
}

func (x concernPolicy) Unfreeze(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion motionproto.Motion,
	policyNS ns.NS,
	args ...any,

) (motionproto.Report, notice.Notices) {

	// unfreeze the priority poll ballot, if frozen
	priorityPoll := pmp_1.ConcernPollBallotName(motion.ID)
	if !ballotapi.IsFrozen_Local(ctx, cloned.PublicClone(), priorityPoll) {
		return nil, nil
	}
	ballotapi.Unfreeze_StageOnly(ctx, cloned, priorityPoll)

	return nil, notice.Noticef(ctx, "This issue, managed by Gov4Git concern `%v`, has been unfrozen 🌤️", motion.ID)
}

// motion.Un/Freeze --calls--> policy Un/Freeze --calls--> ballot Un/Freeze
