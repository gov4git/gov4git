package concern

import (
	"bytes"
	"context"
	"fmt"
	"reflect"

	"github.com/gov4git/gov4git/v2/proto/account"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotapi"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/history/metric"
	"github.com/gov4git/gov4git/v2/proto/member"
	"github.com/gov4git/gov4git/v2/proto/motion/motionapi"
	"github.com/gov4git/gov4git/v2/proto/motion/motionpolicies/waimea"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
	"github.com/gov4git/gov4git/v2/proto/notice"
	"github.com/gov4git/gov4git/v2/proto/purpose"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func init() {
	motionproto.Install(context.Background(), waimea.ConcernPolicyName, ConcernPolicy)
}

var ConcernPolicy concernPolicy

type concernPolicy struct{}

func (x concernPolicy) Descriptor() motionproto.PolicyDescriptor {
	return motionproto.PolicyDescriptor{
		Description:       "Waimea Collective Governance Protocol",
		GithubLabel:       waimea.ConcernPolicyGithubLabel,
		AppliesToConcern:  true,
		AppliesToProposal: false,
	}
}

func (x concernPolicy) PostClone(
	ctx context.Context,
	cloned gov.OwnerCloned,
) {

	err := must.Try(
		func() {
			waimea.LoadConcernClassState_Local(ctx, cloned)
		},
	)
	if git.IsNotExist(err) {
		waimea.SaveConcernClassState_StageOnly(ctx, cloned, waimea.InitialPolicyState)
	}

	waimea.Boot_StageOnly(ctx, cloned.PublicClone())
}

func (x concernPolicy) Open(
	ctx context.Context,
	cloned gov.OwnerCloned,
	con motionproto.Motion,
	args ...any,

) (motionproto.Report, notice.Notices) {

	// initialize state
	policyState := waimea.LoadConcernClassState_Local(ctx, cloned)
	state := waimea.NewConcernState(con.ID, policyState.PriorityMatch)
	motionapi.SavePolicyState_StageOnly[*waimea.ConcernState](ctx, cloned.PublicClone(), con.ID, state)

	// open a priority poll for the motion
	ballotapi.Open_StageOnly(
		ctx,
		ConcernPriorityPollPolicyName,
		cloned,
		state.PriorityPoll,
		waimea.ConcernAccountID(con.ID),
		purpose.Concern,
		con.Policy,
		fmt.Sprintf("Prioritization poll for motion %v", con.ID),
		fmt.Sprintf("Vote for the priority of concern (issue) %v", con.ID),
		[]string{waimea.ConcernBallotChoice},
		member.Everybody,
	)

	// metrics
	metric.Log_StageOnly(ctx, cloned.PublicClone(), &metric.Event{
		Motion: &metric.MotionEvent{
			Open: &metric.MotionOpen{
				ID:     metric.MotionID(con.ID),
				Type:   "concern",
				Policy: metric.MotionPolicy(con.Policy),
			},
		},
	})

	return nil, notice.Noticef(ctx, "Started managing this issue as Gov4Git concern `%v` with initial __priority score__ of `%0.6f`."+
		waimea.Welcome, con.ID, state.PriorityScore)
}

func (x concernPolicy) Score(
	ctx context.Context,
	cloned gov.OwnerCloned,
	con motionproto.Motion,
	args ...any,

) (motionproto.Score, notice.Notices) {

	state := motionapi.LoadPolicyState_Local[*waimea.ConcernState](ctx, cloned.PublicClone(), con.ID)
	return motionproto.Score{
		Attention: state.PriorityScore,
	}, nil
}

func (x concernPolicy) Update(
	ctx context.Context,
	cloned gov.OwnerCloned,
	con motionproto.Motion,
	args ...any,

) (motionproto.Report, notice.Notices) {

	// outputs
	notices := notice.Notices{}

	// inputs
	policyState := waimea.LoadConcernClassState_Local(ctx, cloned)
	conStatePrev := motionapi.LoadPolicyState_Local[*waimea.ConcernState](ctx, cloned.PublicClone(), con.ID)
	conState := conStatePrev.Copy()
	ads := ballotapi.Show_Local(ctx, cloned.PublicClone(), conState.PriorityPoll)

	// update cost of priority
	conState.CostOfPriority = ads.Tally.Capitalization()

	// update priority score
	conState.PriorityMatch = policyState.PriorityMatch
	conState.PriorityScore = ads.Tally.Scores[waimea.ConcernBallotChoice]

	// update eligible proposals
	conState.EligibleProposals = computeEligibleProposals(ctx, cloned.PublicClone(), con)

	// notices
	if !reflect.DeepEqual(conState, conStatePrev) {

		notices = append(
			notices,
			notice.Noticef(ctx,
				"This issue's __priority score__ is now `%0.6f`.\n"+
					"The __cost of priority__ is `%0.6f`.\n"+
					"The __projected bounty__ is now `%0.6f`.",
				conState.PriorityScore, conState.CostOfPriority, conState.ProjectedBounty())...,
		)

		// display updated list of eligible proposals
		if len(conState.EligibleProposals) == 0 {
			notices = append(
				notices,
				notice.Noticef(ctx, "The set of eligible proposals claiming this issue is empty.\n")...,
			)
		} else {
			var w bytes.Buffer
			for _, ref := range conState.EligibleProposals {
				prop := motionapi.LookupMotion_Local(ctx, cloned.PublicClone(), ref.From)
				propState := motionapi.LoadPolicyState_Local[*waimea.ProposalState](ctx, cloned.PublicClone(), prop.ID)
				fmt.Fprintf(&w, "- %s, managed as Gov4Git proposal `%v` with approval score of `%0.6f`\n",
					prop.TrackerURL, prop.ID, propState.ApprovalScore)
			}
			notices = append(
				notices,
				notice.Noticef(ctx, "The set of eligible proposals claiming this issue changed:\n"+w.String())...,
			)
		}

	}

	//

	motionapi.SavePolicyState_StageOnly[*waimea.ConcernState](ctx, cloned.PublicClone(), con.ID, conState)

	r0, n0 := x.updateFreeze(ctx, cloned, con)
	return r0, append(notices, n0...)
}

func computeEligibleProposals(ctx context.Context, cloned gov.Cloned, con motionproto.Motion) motionproto.Refs {
	eligible := motionproto.Refs{}
	for _, ref := range con.RefBy {
		if waimea.AreEligible(ctx, cloned, con.ID, ref.From, ref.Type) {
			eligible = append(eligible, ref)
		}
	}
	eligible.Sort()
	return eligible
}

func (x concernPolicy) updateFreeze(
	ctx context.Context,
	cloned gov.OwnerCloned,
	con motionproto.Motion,
	args ...any,

) (motionproto.Report, notice.Notices) {

	if con.Closed {
		return nil, nil
	}

	conState := motionapi.LoadPolicyState_Local[*waimea.ConcernState](ctx, cloned.PublicClone(), con.ID)

	notices := notice.Notices{}
	if conState.EligibleProposals.Len() > 0 && !con.Frozen {
		motionapi.FreezeMotion_StageOnly(notice.Mute(ctx), cloned, con.ID)

		var w bytes.Buffer
		fmt.Fprintf(&w, "Freezing ❄️ this issue as there are eligible PRs addressing it:\n")
		for _, pr := range conState.EligibleProposals {
			pr := motionapi.LookupMotion_Local(ctx, cloned.PublicClone(), pr.From)
			fmt.Fprintf(&w, "- %s\n", pr.TrackerURL)
		}
		notices = append(notices, notice.Noticef(ctx, w.String())...)
	}
	if conState.EligibleProposals.Len() == 0 && con.Frozen {
		motionapi.UnfreezeMotion_StageOnly(notice.Mute(ctx), cloned, con.ID)
		notices = append(notices, notice.Noticef(ctx, "Unfreezing 🌤️ issue as there are no eligible PRs addressing it.")...)
	}

	return nil, notices
}

func (x concernPolicy) Aggregate(
	ctx context.Context,
	cloned gov.OwnerCloned,
	cons motionproto.Motions,

) {

	cons = motionproto.SelectOpenMotions(cons)

	// load all motion policy states
	conStates := make([]*waimea.ConcernState, len(cons))
	for i, con := range cons {
		conStates[i] = motionapi.LoadPolicyState_Local[*waimea.ConcernState](ctx, cloned.PublicClone(), con.ID)
	}

	// aggregate cost of priority
	totalCostOfPriority := 0.0
	for _, conState := range conStates {
		totalCostOfPriority += conState.CostOfPriority
	}

	// update policy state
	policyState := waimea.LoadConcernClassState_Local(ctx, cloned)
	policyState.TotalCostOfPriority = totalCostOfPriority

	waimea.SaveConcernClassState_StageOnly(ctx, cloned, policyState)
}

func (x concernPolicy) Clear(
	ctx context.Context,
	cloned gov.OwnerCloned,
	con motionproto.Motion,
	args ...any,

) (motionproto.Report, notice.Notices) {

	// clear engages only after close or cancel
	if !con.Closed {
		return nil, nil
	}

	return nil, nil
}

func (x concernPolicy) Close(
	ctx context.Context,
	cloned gov.OwnerCloned,
	con motionproto.Motion,
	decision motionproto.Decision,
	args ...any,
	// args[0]=toID account.AccountID
	// args[1]=prop schema.Motion

) (motionproto.Report, notice.Notices) {

	// NOTE: the first argument is not used by this policy
	must.Assertf(ctx, len(args) == 2, "issue closure requires two arguments, got %v", args)
	_, ok := args[0].(account.AccountID)
	must.Assertf(ctx, ok, "unrecognized account ID argument %v", args[0])
	prop, ok := args[1].(motionproto.Motion)
	must.Assertf(ctx, ok, "unrecognized proposal motion argument %v", args[1])

	conState := motionapi.LoadPolicyState_Local[*waimea.ConcernState](ctx, cloned.PublicClone(), con.ID)

	// cancel the poll for the motion, so that all charges are refunded
	priorityPollName := waimea.ConcernPollBallotName(con.ID)
	chg := ballotapi.Cancel_StageOnly(
		ctx,
		cloned,
		priorityPollName,
	)

	// metrics
	metric.Log_StageOnly(ctx, cloned.PublicClone(), &metric.Event{
		Motion: &metric.MotionEvent{
			Close: &metric.MotionClose{
				ID:       metric.MotionID(con.ID),
				Type:     "concern",
				Decision: decision.MetricDecision(),
				Policy:   metric.MotionPolicy(con.Policy),
				Receipts: chg.Result.RefundedHistoryReceipts(),
			},
		},
	})

	return &CloseReport{}, closeNotice(ctx, con, conState, chg.Result, prop)
}

func (x concernPolicy) Cancel(
	ctx context.Context,
	cloned gov.OwnerCloned,
	con motionproto.Motion,
	args ...any,

) (motionproto.Report, notice.Notices) {

	conState := motionapi.LoadPolicyState_Local[*waimea.ConcernState](ctx, cloned.PublicClone(), con.ID)

	// cancel the poll for the motion (returning credits to users)
	priorityPollName := waimea.ConcernPollBallotName(con.ID)
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
				Type:     "concern",
				Policy:   metric.MotionPolicy(con.Policy),
				Receipts: chg.Result.RefundedHistoryReceipts(),
			},
		},
	})

	return &CancelReport{
		PriorityPollOutcome: chg.Result,
	}, cancelNotice(ctx, con, conState, chg.Result)
}

type PolicyView struct {
	State          *waimea.ConcernState      `json:"state"`
	PriorityPoll   ballotproto.AdTallyMargin `json:"priority_poll"`
	PriorityMargin ballotproto.Margin        `json:"priority_margin"`
}

func (x concernPolicy) Show(
	ctx context.Context,
	cloned gov.Cloned,
	con motionproto.Motion,
	args ...any,

) (form.Form, motionproto.MotionBallots) {

	// retrieve policy state
	conState := motionapi.LoadPolicyState_Local[*waimea.ConcernState](ctx, cloned, con.ID)

	// retrieve poll state
	priorityPollName := waimea.ConcernPollBallotName(con.ID)
	priorityPoll := ballotapi.Show_Local(ctx, cloned, priorityPollName)

	return PolicyView{
			State:          conState,
			PriorityPoll:   priorityPoll,
			PriorityMargin: *ballotapi.GetMargin_Local(ctx, cloned, priorityPollName),
		}, motionproto.MotionBallots{
			motionproto.MotionBallot{
				Label:         "priority_poll",
				BallotID:      conState.PriorityPoll,
				BallotChoices: priorityPoll.Ad.Choices,
				BallotAd:      priorityPoll.Ad,
				BallotTally:   priorityPoll.Tally,
				BallotMargin:  priorityPoll.Margin,
			},
		}
}

func (x concernPolicy) AddRefTo(
	ctx context.Context,
	cloned gov.OwnerCloned,
	refType motionproto.RefType,
	from motionproto.Motion,
	to motionproto.Motion,
	args ...any,

) (motionproto.Report, notice.Notices) {

	if !from.IsProposal() {
		return nil, nil
	}

	if refType != waimea.ClaimsRefType {
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
	args ...any,

) (motionproto.Report, notice.Notices) {

	if !from.IsProposal() {
		return nil, nil
	}

	if refType != waimea.ClaimsRefType {
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
	args ...any,

) (motionproto.Report, notice.Notices) {

	return nil, nil
}

func (x concernPolicy) Freeze(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion motionproto.Motion,
	args ...any,

) (motionproto.Report, notice.Notices) {

	// freeze priority poll, if not already frozen
	priorityPoll := waimea.ConcernPollBallotName(motion.ID)
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
	args ...any,

) (motionproto.Report, notice.Notices) {

	// unfreeze the priority poll ballot, if frozen
	priorityPoll := waimea.ConcernPollBallotName(motion.ID)
	if !ballotapi.IsFrozen_Local(ctx, cloned.PublicClone(), priorityPoll) {
		return nil, nil
	}
	ballotapi.Unfreeze_StageOnly(ctx, cloned, priorityPoll)

	return nil, notice.Noticef(ctx, "This issue, managed by Gov4Git concern `%v`, has been unfrozen 🌤️", motion.ID)
}
