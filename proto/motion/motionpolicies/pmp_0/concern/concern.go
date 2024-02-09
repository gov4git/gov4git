package concern

import (
	"bytes"
	"context"
	"fmt"
	"reflect"

	"github.com/gov4git/gov4git/v2/proto/account"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotapi"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotio"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/history/metric"
	"github.com/gov4git/gov4git/v2/proto/member"
	"github.com/gov4git/gov4git/v2/proto/motion/motionapi"
	"github.com/gov4git/gov4git/v2/proto/motion/motionpolicies/pmp_0"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
	"github.com/gov4git/gov4git/v2/proto/notice"
	"github.com/gov4git/gov4git/v2/proto/purpose"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/must"
)

func init() {
	motionproto.Install(context.Background(), pmp_0.ConcernPolicyName, concernPolicy{})
}

type concernPolicy struct{}

func (x concernPolicy) PostClone(
	ctx context.Context,
	cloned gov.OwnerCloned,
) {
}

func (x concernPolicy) Open(
	ctx context.Context,
	cloned gov.OwnerCloned,
	con motionproto.Motion,
	args ...any,

) (motionproto.Report, notice.Notices) {

	// initialize state
	state := pmp_0.NewConcernState(con.ID)
	motionapi.SavePolicyState_StageOnly[*pmp_0.ConcernState](ctx, cloned.PublicClone(), con.ID, state)

	// open a poll for the motion
	ballotapi.Open_StageOnly(
		ctx,
		ballotio.QVPolicyName,
		cloned,
		state.PriorityPoll,
		pmp_0.ConcernAccountID(con.ID),
		purpose.Concern,
		con.Policy,
		fmt.Sprintf("Prioritization poll for motion %v", con.ID),
		fmt.Sprintf("Up/down vote the priority for concern (issue) %v", con.ID),
		[]string{pmp_0.ConcernBallotChoice},
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
		pmp_0.Welcome, con.ID, state.LatestPriorityScore)
}

func (x concernPolicy) Score(
	ctx context.Context,
	cloned gov.OwnerCloned,
	con motionproto.Motion,
	args ...any,

) (motionproto.Score, notice.Notices) {

	state := motionapi.LoadPolicyState_Local[*pmp_0.ConcernState](ctx, cloned.PublicClone(), con.ID)

	// compute motion score from the priority poll ballot
	ads := ballotapi.Show_Local(ctx, cloned.PublicClone(), state.PriorityPoll)
	attention := ads.Tally.Attention()

	return motionproto.Score{
		Attention: attention,
	}, nil
}

func (x concernPolicy) Update(
	ctx context.Context,
	cloned gov.OwnerCloned,
	con motionproto.Motion,
	args ...any,

) (motionproto.Report, notice.Notices) {

	notices := notice.Notices{}
	conStatePrev := motionapi.LoadPolicyState_Local[*pmp_0.ConcernState](ctx, cloned.PublicClone(), con.ID)
	conState := conStatePrev.Copy()

	// update priority score

	ads := ballotapi.Show_Local(ctx, cloned.PublicClone(), conState.PriorityPoll)
	latestPriorityScore := ads.Tally.Scores[pmp_0.ConcernBallotChoice]
	costOfPriority := ads.Tally.Capitalization()
	conState.CostOfPriority = costOfPriority
	conState.LatestPriorityScore = latestPriorityScore

	// update eligible proposals

	conState.EligibleProposals = computeEligibleProposals(ctx, cloned.PublicClone(), con)

	// notices

	if !reflect.DeepEqual(conState, conStatePrev) {

		notices = append(
			notices,
			notice.Noticef(ctx, "This issue's __priority score__ is now `%0.6f`.\n"+
				"The __cost of priority__ is `%0.6f`.\n"+
				"The __projected bounty__ is now `%0.6f`.", latestPriorityScore, costOfPriority, costOfPriority)...,
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
				fmt.Fprintf(&w, "- %s, managed as Gov4Git motion `%v` with community attention of `%0.6f`\n",
					prop.TrackerURL, prop.ID, prop.Score.Attention)
			}
			notices = append(
				notices,
				notice.Noticef(ctx, "The set of eligible proposals claiming this issue changed:\n"+w.String())...,
			)
		}

	}

	//

	motionapi.SavePolicyState_StageOnly[*pmp_0.ConcernState](ctx, cloned.PublicClone(), con.ID, conState)

	r0, n0 := x.updateFreeze(ctx, cloned, con)
	return r0, append(notices, n0...)
}

func computeEligibleProposals(ctx context.Context, cloned gov.Cloned, con motionproto.Motion) motionproto.Refs {
	eligible := motionproto.Refs{}
	for _, ref := range con.RefBy {
		if pmp_0.IsConcernProposalEligible(ctx, cloned, con.ID, ref.From, ref.Type) {
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

	toState := motionapi.LoadPolicyState_Local[*pmp_0.ConcernState](ctx, cloned.PublicClone(), con.ID)

	notices := notice.Notices{}
	if toState.EligibleProposals.Len() > 0 && !con.Frozen {
		motionapi.FreezeMotion_StageOnly(notice.Mute(ctx), cloned, con.ID)

		var w bytes.Buffer
		fmt.Fprintf(&w, "Freezing ❄️ this issue as there are eligible PRs addressing it:\n")
		for _, pr := range toState.EligibleProposals {
			pr := motionapi.LookupMotion_Local(ctx, cloned.PublicClone(), pr.From)
			fmt.Fprintf(&w, "- %s\n", pr.TrackerURL)
		}
		notices = append(notices, notice.Noticef(ctx, w.String())...)
	}
	if toState.EligibleProposals.Len() == 0 && con.Frozen {
		motionapi.UnfreezeMotion_StageOnly(notice.Mute(ctx), cloned, con.ID)
		notices = append(notices, notice.Noticef(ctx, "Unfreezing 🌤️ issue as there are no eligible PRs addressing it.")...)
	}

	return nil, notices
}

func (x concernPolicy) Aggregate(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion motionproto.Motions,
) {
}

func (x concernPolicy) Clear(
	ctx context.Context,
	cloned gov.OwnerCloned,
	con motionproto.Motion,
	args ...any,

) (motionproto.Report, notice.Notices) {
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

	must.Assertf(ctx, len(args) == 2, "issue closure requires two arguments, got %v", args)
	toID, ok := args[0].(account.AccountID)
	must.Assertf(ctx, ok, "unrecognized account ID argument %v", args[0])
	prop, ok := args[1].(motionproto.Motion)
	must.Assertf(ctx, ok, "unrecognized proposal motion argument %v", args[1])

	// close the poll for the motion
	priorityPollName := pmp_0.ConcernPollBallotName(con.ID)
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
				Type:     "concern",
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
	args ...any,

) (motionproto.Report, notice.Notices) {

	// cancel the poll for the motion (returning credits to users)
	priorityPollName := pmp_0.ConcernPollBallotName(con.ID)
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
				Receipts: nil, // refunds are accounted for by the proposal
			},
		},
	})

	return &CancelReport{
		PriorityPollOutcome: chg.Result,
	}, cancelNotice(ctx, con, chg.Result)
}

type PolicyView struct {
	State          *pmp_0.ConcernState       `json:"state"`
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
	policyState := motionapi.LoadPolicyState_Local[*pmp_0.ConcernState](ctx, cloned, con.ID)

	// retrieve poll state
	priorityPollName := pmp_0.ConcernPollBallotName(con.ID)
	priorityPoll := ballotapi.Show_Local(ctx, cloned, priorityPollName)

	return PolicyView{
			State:          policyState,
			PriorityPoll:   priorityPoll,
			PriorityMargin: *ballotapi.GetMargin_Local(ctx, cloned, priorityPollName),
		}, motionproto.MotionBallots{
			motionproto.MotionBallot{
				Label:         "priority_poll",
				BallotID:      policyState.PriorityPoll,
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

	if refType != pmp_0.ClaimsRefType {
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

	if refType != pmp_0.ClaimsRefType {
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
	priorityPoll := pmp_0.ConcernPollBallotName(motion.ID)
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
	priorityPoll := pmp_0.ConcernPollBallotName(motion.ID)
	if !ballotapi.IsFrozen_Local(ctx, cloned.PublicClone(), priorityPoll) {
		return nil, nil
	}
	ballotapi.Unfreeze_StageOnly(ctx, cloned, priorityPoll)

	return nil, notice.Noticef(ctx, "This issue, managed by Gov4Git concern `%v`, has been unfrozen 🌤️", motion.ID)
}

// motion.Un/Freeze --calls--> policy Un/Freeze --calls--> ballot Un/Freeze
