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
	"github.com/gov4git/gov4git/v2/proto/motion"
	"github.com/gov4git/gov4git/v2/proto/motion/motionapi"
	"github.com/gov4git/gov4git/v2/proto/motion/motionpolicies/pmp_1"
	"github.com/gov4git/gov4git/v2/proto/motion/motionpolicies/pmp_1/concern"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
	"github.com/gov4git/gov4git/v2/proto/notice"
	"github.com/gov4git/gov4git/v2/proto/purpose"
	"github.com/gov4git/lib4git/form"
)

func init() {
	motionproto.Install(context.Background(), ProposalPolicyName, proposalPolicy{})
}

const (
	ProposalPolicyName             motion.PolicyName      = "pmp-proposal-v1"
	ProposalApprovalPollPolicyName ballotproto.PolicyName = "pmp-proposal-approval-v1"
)

type proposalPolicy struct{}

func (x proposalPolicy) PostClone(
	ctx context.Context,
	cloned gov.OwnerCloned,
) {
}

func (x proposalPolicy) Open(
	ctx context.Context,
	cloned gov.OwnerCloned,
	prop motionproto.Motion,
	args ...any,

) (motionproto.Report, notice.Notices) {

	// initialize state
	state := NewProposalState(prop.ID)
	motionapi.SavePolicyState_StageOnly[*ProposalState](ctx, cloned.PublicClone(), prop.ID, state)

	// create a bounty account for the proposal
	account.Create_StageOnly(
		ctx,
		cloned.PublicClone(),
		pmp_1.ProposalBountyAccountID(prop.ID),
		pmp_1.ProposalAccountID(prop.ID),
		fmt.Sprintf("bounty account for proposal %v", prop.ID),
	)

	// create a reward account for the proposal
	account.Create_StageOnly(
		ctx,
		cloned.PublicClone(),
		pmp_1.ProposalRewardAccountID(prop.ID),
		pmp_1.ProposalAccountID(prop.ID),
		fmt.Sprintf("reward account for proposal %v", prop.ID),
	)

	// open a poll for the motion
	ballotapi.Open_StageOnly(
		ctx,
		ProposalApprovalPollPolicyName,
		cloned,
		state.ApprovalPoll,
		pmp_1.ProposalAccountID(prop.ID),
		purpose.Proposal,
		prop.Policy,
		fmt.Sprintf("Approval referendum for motion %v", prop.ID),
		fmt.Sprintf("Up/down vote the approval vote for proposal (pull request) %v", prop.ID),
		[]string{pmp_1.ProposalBallotChoice},
		member.Everybody,
	)
	zeroState := ScoreKernelState{
		MotionID:              prop.ID,
		Bounty:                0.0,
		InverseCostMultiplier: 1.0,
	}
	ballotapi.SavePolicyState_StageOnly[ScoreKernelState](
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
				Type:   "proposal-v1",
				Policy: metric.MotionPolicy(prop.Policy),
			},
		},
	})

	return nil, notice.Noticef(ctx,
		"Started managing this PR as Gov4Git proposal `%v` with initial __approval score__ of `%0.6f`."+
			pmp_1.Welcome, prop.ID, state.LatestApprovalScore)
}

func (x proposalPolicy) Score(
	ctx context.Context,
	cloned gov.OwnerCloned,
	prop motionproto.Motion,
	args ...any,

) (motionproto.Score, notice.Notices) {

	state := motionapi.LoadPolicyState_Local[*ProposalState](ctx, cloned.PublicClone(), prop.ID)
	return motionproto.Score{
		Attention: state.LatestApprovalScore,
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

	conPolicyState := concern.LoadClassState_Local(ctx, cloned)
	propStatePrev := motionapi.LoadPolicyState_Local[*ProposalState](ctx, cloned.PublicClone(), prop.ID)
	propState := propStatePrev.Copy()
	ads := ballotapi.Show_Local(ctx, cloned.Public.Tree(), propState.ApprovalPoll)

	// update approval score

	latestApprovalScore := ads.Tally.Scores[pmp_1.ProposalBallotChoice]
	costOfReview := ads.Tally.Capitalization()
	propState.LatestApprovalScore = latestApprovalScore

	// update eligible concerns

	propState.EligibleConcerns = calcEligibleConcerns(ctx, cloned.PublicClone(), prop)

	// update cost multiplier

	projectedBounty := 0.0
	for _, ref := range propState.EligibleConcerns {
		conState := motionapi.LoadPolicyState_Local[*concern.ConcernState](ctx, cloned.PublicClone(), ref.To)
		projectedBounty += conState.ProjectedBounty()
	}

	inverseCostMultiplier := (4 * conPolicyState.WithheldEscrowFraction * projectedBounty) / (1 + float64(ads.Tally.NumVoters()))
	propState.InverseCostMultiplier = inverseCostMultiplier

	// notices

	if !reflect.DeepEqual(propState, propStatePrev) {

		notices = append(
			notices,
			notice.Noticef(ctx, "This PR's __approval score__ is now `%0.6f`.\n"+
				"The __cost of review__ is `%0.6f`.\n"+
				"The __projected bounty__ is now `%0.6f`.", latestApprovalScore, costOfReview, projectedBounty)...,
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
				fmt.Fprintf(&w, "- %s, managed as Gov4Git motion `%v` with community attention of `%0.6f`\n",
					con.TrackerURL, con.ID, con.Score.Attention)
			}
			notices = append(
				notices,
				notice.Noticef(ctx, "The set of eligible issues claimed by this PR changed:\n"+w.String())...,
			)
		}

	}

	//

	motionapi.SavePolicyState_StageOnly[*ProposalState](ctx, cloned.PublicClone(), prop.ID, propState)

	// update ScoreKernelState
	currentState := ScoreKernelState{
		MotionID:              prop.ID,
		Bounty:                projectedBounty,
		InverseCostMultiplier: inverseCostMultiplier,
	}
	ballotapi.SavePolicyState_StageOnly[ScoreKernelState](
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
	motion motionproto.Motions,
) {
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

	// cancel the referendum for the motion (and return credits to users)
	referendumName := pmp_1.ProposalApprovalPollName(prop.ID)
	chg := ballotapi.Cancel_StageOnly(
		ctx,
		cloned,
		referendumName,
	)

	// metrics
	metric.Log_StageOnly(ctx, cloned.PublicClone(), &metric.Event{
		Motion: &metric.MotionEvent{
			Cancel: &metric.MotionCancel{
				ID:       metric.MotionID(prop.ID),
				Type:     "proposals",
				Policy:   metric.MotionPolicy(prop.Policy),
				Receipts: chg.Result.RefundedHistoryReceipts(),
			},
		},
	})

	return &CancelReport{
		ApprovalPollOutcome: chg.Result,
	}, cancelNotice(ctx, prop, chg.Result)
}

type PolicyView struct {
	State          *ProposalState      `json:"state"`
	ApprovalPoll   ballotproto.AdTally `json:"approval_poll"`
	ApprovalMargin ballotproto.Margin  `json:"priority_margin"`
	BountyAccount  account.AccountID   `json:"bounty_account"`
	RewardAccount  account.AccountID   `json:"reward_account"`
}

func (x proposalPolicy) Show(
	ctx context.Context,
	cloned gov.Cloned,
	prop motionproto.Motion,
	args ...any,

) (form.Form, motionproto.MotionBallots) {

	// retrieve policy state
	policyState := motionapi.LoadPolicyState_Local[*ProposalState](ctx, cloned, prop.ID)

	// retrieve approval poll
	approvalPoll := loadPropApprovalPollTally(ctx, cloned, prop)

	return PolicyView{
			State:          policyState,
			ApprovalPoll:   approvalPoll,
			ApprovalMargin: *ballotapi.GetMargin_Local(ctx, cloned, approvalPoll.Ad.ID),
			BountyAccount:  pmp_1.ProposalBountyAccountID(prop.ID),
			RewardAccount:  pmp_1.ProposalRewardAccountID(prop.ID),
		}, motionproto.MotionBallots{
			motionproto.MotionBallot{
				Label:         "approval_poll",
				BallotID:      policyState.ApprovalPoll,
				BallotChoices: approvalPoll.Ad.Choices,
				BallotAd:      approvalPoll.Ad,
				BallotTally:   approvalPoll.Tally,
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

	if refType != pmp_1.ClaimsRefType {
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

	if refType != pmp_1.ClaimsRefType {
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
