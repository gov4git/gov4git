package proposal

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
	"github.com/gov4git/gov4git/v2/proto/motion/motionpolicies/waimea/concern"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
	"github.com/gov4git/gov4git/v2/proto/notice"
	"github.com/gov4git/gov4git/v2/proto/purpose"
	"github.com/gov4git/lib4git/form"
)

func init() {
	motionproto.Install(context.Background(), waimea.ProposalPolicyName, proposalPolicy{})
}

type proposalPolicy struct{}

func (x proposalPolicy) Descriptor() motionproto.PolicyDescriptor {
	return motionproto.PolicyDescriptor{
		Description:       "Waimea Collective Governance Protocol",
		GithubLabel:       waimea.ConcernPolicyGithubLabel,
		AppliesToConcern:  false,
		AppliesToProposal: true,
	}
}

func (x proposalPolicy) PostClone(
	ctx context.Context,
	cloned gov.OwnerCloned,
) {

	concern.ConcernPolicy.PostClone(ctx, cloned)
}

func (x proposalPolicy) Open(
	ctx context.Context,
	cloned gov.OwnerCloned,
	prop motionproto.Motion,
	args ...any,

) (motionproto.Report, notice.Notices) {

	// initialize state
	policyState := waimea.LoadConcernClassState_Local(ctx, cloned)
	state := waimea.NewProposalState(prop.ID, policyState.ReviewMatch)
	motionapi.SavePolicyState_StageOnly[*waimea.ProposalState](ctx, cloned.PublicClone(), prop.ID, state)

	// create a bounty account for the proposal
	account.Create_StageOnly(
		ctx,
		cloned.PublicClone(),
		waimea.ProposalBountyAccountID(prop.ID),
		waimea.ProposalAccountID(prop.ID),
		fmt.Sprintf("bounty account for proposal %v", prop.ID),
	)

	// create a reward account for the proposal
	account.Create_StageOnly(
		ctx,
		cloned.PublicClone(),
		waimea.ProposalRewardAccountID(prop.ID),
		waimea.ProposalAccountID(prop.ID),
		fmt.Sprintf("reward account for proposal %v", prop.ID),
	)

	// open a poll for the motion
	ballotapi.Open_StageOnly(
		ctx,
		ProposalApprovalPollPolicyName,
		cloned,
		state.ApprovalPoll,
		waimea.ProposalAccountID(prop.ID),
		purpose.Proposal,
		prop.Policy,
		fmt.Sprintf("Approval poll for motion %v", prop.ID),
		fmt.Sprintf("Vote for the approval of proposal (pull request) %v", prop.ID),
		[]string{waimea.ProposalBallotChoice},
		member.Everybody,
	)
	zeroState := ApprovalPollState{
		MotionID:              prop.ID,
		Bounty:                0.0,
		InverseCostMultiplier: 1.0,
	}
	ballotapi.SavePolicyState_StageOnly[ApprovalPollState](
		ctx,
		cloned.PublicClone(),
		state.ApprovalPoll,
		zeroState,
	)

	// metrics
	metric.Log_StageOnly(ctx, cloned.PublicClone(), &metric.Event{
		Motion: &metric.MotionEvent{
			Open: &metric.MotionOpen{
				ID:     metric.MotionID(prop.ID),
				Type:   "proposal",
				Policy: metric.MotionPolicy(prop.Policy),
			},
		},
	})

	return nil, notice.Noticef(ctx,
		"Started managing this PR as Gov4Git proposal `%v` with initial __approval score__ of `%0.6f`."+
			waimea.Welcome, prop.ID, state.ApprovalScore)
}

func (x proposalPolicy) Score(
	ctx context.Context,
	cloned gov.OwnerCloned,
	prop motionproto.Motion,
	args ...any,

) (motionproto.Score, notice.Notices) {

	state := motionapi.LoadPolicyState_Local[*waimea.ProposalState](ctx, cloned.PublicClone(), prop.ID)
	return motionproto.Score{
		Attention: state.ApprovalScore,
	}, nil
}

func (x proposalPolicy) Update(
	ctx context.Context,
	cloned gov.OwnerCloned,
	prop motionproto.Motion,
	args ...any,

) (motionproto.Report, notice.Notices) {

	notices := notice.Notices{}

	// inputs
	policyState := waimea.LoadConcernClassState_Local(ctx, cloned)
	propStatePrev := motionapi.LoadPolicyState_Local[*waimea.ProposalState](ctx, cloned.PublicClone(), prop.ID)
	propState := propStatePrev.Copy()
	ads := ballotapi.Show_Local(ctx, cloned.PublicClone(), propState.ApprovalPoll)

	// update cost of review
	propState.CostOfReview = ads.Tally.Capitalization()

	// update approval score
	propState.ReviewMatch = policyState.ReviewMatch
	propState.ApprovalScore = ads.Tally.Scores[waimea.ProposalBallotChoice]

	// update eligible concerns
	propState.EligibleConcerns = computeEligibleConcerns(ctx, cloned.PublicClone(), prop)

	// update projected bounty
	projectedBounty := 0.0
	for _, ref := range propState.EligibleConcerns {
		conState := motionapi.LoadPolicyState_Local[*waimea.ConcernState](ctx, cloned.PublicClone(), ref.To)
		projectedBounty += conState.ProjectedBounty()
	}
	propState.ProjectedPriorityBounty = projectedBounty

	// notices

	if !reflect.DeepEqual(propState, propStatePrev) {

		notices = append(
			notices,
			notice.Noticef(ctx, "This PR's __approval score__ is now `%0.6f`.\n"+
				"The __cost of review__ is `%0.6f`.\n"+
				"The __projected bounty__ is now `%0.6f`.",
				propState.ApprovalScore, propState.CostOfReview, propState.ProjectedPriorityBounty)...,
		)

		// display list of eligible concerns
		if len(propState.EligibleConcerns) == 0 {
			notices = append(
				notices,
				notice.Noticef(ctx, "The set of eligible issues claimed by this PR is empty.\n")...,
			)
		} else {
			var w bytes.Buffer
			for _, ref := range propState.EligibleConcerns {
				con := motionapi.LookupMotion_Local(ctx, cloned.PublicClone(), ref.To)
				fmt.Fprintf(&w, "- %s, managed as Gov4Git concern `%v` with __priority score__ of `%0.6f`\n",
					con.TrackerURL, con.ID, con.Score.Attention)
			}
			notices = append(
				notices,
				notice.Noticef(ctx, "The set of eligible issues claimed by this PR changed:\n"+w.String())...,
			)
		}

	}

	//

	motionapi.SavePolicyState_StageOnly[*waimea.ProposalState](ctx, cloned.PublicClone(), prop.ID, propState)

	// update approval poll state
	currentState := ApprovalPollState{
		MotionID:              prop.ID,
		Bounty:                projectedBounty,
		InverseCostMultiplier: 1.0,
	}
	ballotapi.SavePolicyState_StageOnly[ApprovalPollState](
		ctx,
		cloned.PublicClone(),
		propState.ApprovalPoll,
		currentState,
	)

	return nil, notices
}

func (x proposalPolicy) Aggregate(
	ctx context.Context,
	cloned gov.OwnerCloned,
	props motionproto.Motions,

) {

	props = motionproto.SelectOpenMotions(props)

	// load all motion policy states
	propStates := make([]*waimea.ProposalState, len(props))
	for i, prop := range props {
		propStates[i] = motionapi.LoadPolicyState_Local[*waimea.ProposalState](ctx, cloned.PublicClone(), prop.ID)
	}

	// aggregate cost of review
	totalCostOfReview := 0.0
	for _, propState := range propStates {
		totalCostOfReview += propState.CostOfReview
	}

	// update policy state
	policyState := waimea.LoadConcernClassState_Local(ctx, cloned)
	policyState.TotalCostOfReview = totalCostOfReview

	waimea.SaveConcernClassState_StageOnly(ctx, cloned, policyState)

}

func (x proposalPolicy) Clear(
	ctx context.Context,
	cloned gov.OwnerCloned,
	prop motionproto.Motion,
	_ ...any,

) (motionproto.Report, notice.Notices) {

	return nil, nil
}

func (x proposalPolicy) Cancel(
	ctx context.Context,
	cloned gov.OwnerCloned,
	prop motionproto.Motion,
	args ...any,

) (motionproto.Report, notice.Notices) {

	propState := motionapi.LoadPolicyState_Local[*waimea.ProposalState](ctx, cloned.PublicClone(), prop.ID)

	// cancel the approval poll (and return credits to users)
	approvalPollName := waimea.ProposalApprovalPollName(prop.ID)
	chg := ballotapi.Cancel_StageOnly(
		ctx,
		cloned,
		approvalPollName,
	)

	// metrics
	metric.Log_StageOnly(ctx, cloned.PublicClone(), &metric.Event{
		Motion: &metric.MotionEvent{
			Cancel: &metric.MotionCancel{
				ID:       metric.MotionID(prop.ID),
				Type:     "proposal",
				Policy:   metric.MotionPolicy(prop.Policy),
				Receipts: chg.Result.RefundedHistoryReceipts(),
			},
		},
	})

	return &CancelReport{
		ApprovalPollOutcome: chg.Result,
	}, cancelNotice(ctx, prop, propState, chg.Result)
}

type PolicyView struct {
	State          *waimea.ProposalState     `json:"state"`
	ApprovalPoll   ballotproto.AdTallyMargin `json:"approval_poll"`
	ApprovalMargin ballotproto.Margin        `json:"priority_margin"`
	BountyAccount  account.AccountID         `json:"bounty_account"`
	RewardAccount  account.AccountID         `json:"reward_account"`
}

func (x proposalPolicy) Show(
	ctx context.Context,
	cloned gov.Cloned,
	prop motionproto.Motion,
	args ...any,

) (form.Form, motionproto.MotionBallots) {

	// retrieve policy state
	propState := motionapi.LoadPolicyState_Local[*waimea.ProposalState](ctx, cloned, prop.ID)

	// retrieve approval poll
	approvalPoll := loadApprovalPoll(ctx, cloned, prop)

	return PolicyView{
			State:          propState,
			ApprovalPoll:   approvalPoll,
			ApprovalMargin: *ballotapi.GetMargin_Local(ctx, cloned, approvalPoll.Ad.ID),
			BountyAccount:  waimea.ProposalBountyAccountID(prop.ID),
			RewardAccount:  waimea.ProposalRewardAccountID(prop.ID),
		}, motionproto.MotionBallots{
			motionproto.MotionBallot{
				Label:         "approval_poll",
				BallotID:      propState.ApprovalPoll,
				BallotChoices: approvalPoll.Ad.Choices,
				BallotAd:      approvalPoll.Ad,
				BallotTally:   approvalPoll.Tally,
				BallotMargin:  approvalPoll.Margin,
			},
		}
}

func (x proposalPolicy) AddRefTo(
	ctx context.Context,
	cloned gov.OwnerCloned,
	refType motionproto.RefType,
	from motionproto.Motion,
	to motionproto.Motion,
	args ...any,

) (motionproto.Report, notice.Notices) {

	return nil, nil
}

func (x proposalPolicy) AddRefFrom(
	ctx context.Context,
	cloned gov.OwnerCloned,
	refType motionproto.RefType,
	from motionproto.Motion,
	to motionproto.Motion,
	args ...any,

) (motionproto.Report, notice.Notices) {

	if !to.IsConcern() {
		return nil, nil
	}

	if refType != waimea.ClaimsRefType {
		return nil, nil
	}

	return nil, notice.Noticef(ctx, "This PR referenced %v, managed as Gov4Git concern `%v`.", to.TrackerURL, to.ID)
}

func (x proposalPolicy) RemoveRefTo(
	ctx context.Context,
	cloned gov.OwnerCloned,
	refType motionproto.RefType,
	from motionproto.Motion,
	to motionproto.Motion,
	args ...any,

) (motionproto.Report, notice.Notices) {

	return nil, nil
}

func (x proposalPolicy) RemoveRefFrom(
	ctx context.Context,
	cloned gov.OwnerCloned,
	refType motionproto.RefType,
	from motionproto.Motion,
	to motionproto.Motion,
	args ...any,

) (motionproto.Report, notice.Notices) {

	if !to.IsConcern() {
		return nil, nil
	}

	if refType != waimea.ClaimsRefType {
		return nil, nil
	}

	return nil, notice.Noticef(ctx, "This PR no longer references %v, managed as Gov4Git concern `%v`.", to.TrackerURL, to.ID)
}

func (x proposalPolicy) Freeze(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion motionproto.Motion,
	args ...any,

) (motionproto.Report, notice.Notices) {

	return nil, notice.Noticef(ctx, "This PR, managed by Gov4Git proposal `%v`, has been frozen ❄️", motion.ID)
}

func (x proposalPolicy) Unfreeze(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion motionproto.Motion,
	args ...any,

) (motionproto.Report, notice.Notices) {

	return nil, notice.Noticef(ctx, "This PR, managed by Gov4Git proposal `%v`, has been unfrozen 🌤️", motion.ID)
}
