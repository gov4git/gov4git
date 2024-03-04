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
	"github.com/gov4git/gov4git/v2/proto/motion/motionpolicies/pmp_0"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
	"github.com/gov4git/gov4git/v2/proto/notice"
	"github.com/gov4git/gov4git/v2/proto/purpose"
	"github.com/gov4git/lib4git/form"
)

func init() {
	motionproto.Install(context.Background(), pmp_0.ProposalPolicyName, proposalPolicy{})
}

const (
	ProposalApprovalPollPolicyName ballotproto.PolicyName = "pmp-proposal-approval"
)

type proposalPolicy struct{}

func (x proposalPolicy) Descriptor() motionproto.PolicyDescriptor {
	return motionproto.PolicyDescriptor{
		Description:       "[Plural Management Protocol](https://papers.ssrn.com/sol3/papers.cfm?abstract_id=4688040) v0",
		GithubLabel:       pmp_0.ProposalPolicyGithubLabel,
		AppliesToConcern:  false,
		AppliesToProposal: true,
	}
}

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
	state := pmp_0.NewProposalState(prop.ID)
	motionapi.SavePolicyState_StageOnly[*pmp_0.ProposalState](ctx, cloned.PublicClone(), prop.ID, state)

	// create a bounty account for the proposal
	account.Create_StageOnly(
		ctx,
		cloned.PublicClone(),
		pmp_0.ProposalBountyAccountID(prop.ID),
		pmp_0.ProposalAccountID(prop.ID),
		fmt.Sprintf("bounty account for proposal %v", prop.ID),
	)

	// create a reward account for the proposal
	account.Create_StageOnly(
		ctx,
		cloned.PublicClone(),
		pmp_0.ProposalRewardAccountID(prop.ID),
		pmp_0.ProposalAccountID(prop.ID),
		fmt.Sprintf("reward account for proposal %v", prop.ID),
	)

	// open a poll for the motion
	ballotapi.Open_StageOnly(
		ctx,
		ProposalApprovalPollPolicyName,
		cloned,
		state.ApprovalPoll,
		pmp_0.ProposalAccountID(prop.ID),
		purpose.Proposal,
		prop.Policy,
		fmt.Sprintf("Approval referendum for motion %v", prop.ID),
		fmt.Sprintf("Up/down vote the approval vote for proposal (pull request) %v", prop.ID),
		[]string{pmp_0.ProposalBallotChoice},
		member.Everybody,
	)
	zeroState := ScoreKernelState{
		Bounty: 0.0,
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
				Type:   "proposal",
				Policy: metric.MotionPolicy(prop.Policy),
			},
		},
	})

	return nil, notice.Noticef(ctx,
		"Started managing this PR, using the Plural Management Protocol v0, as Gov4Git proposal `%v` with initial __approval score__ of `%0.6f`."+
			pmp_0.Welcome, prop.ID, state.LatestApprovalScore)
}

func (x proposalPolicy) Score(
	ctx context.Context,
	cloned gov.OwnerCloned,
	prop motionproto.Motion,
	args ...any,

) (motionproto.Score, notice.Notices) {

	state := motionapi.LoadPolicyState_Local[*pmp_0.ProposalState](ctx, cloned.PublicClone(), prop.ID)

	// compute score
	ads := ballotapi.Show_Local(ctx, cloned.PublicClone(), state.ApprovalPoll)
	attention := ads.Tally.Attention()

	return motionproto.Score{
		Attention: attention,
	}, nil
}

func (x proposalPolicy) Update(
	ctx context.Context,
	cloned gov.OwnerCloned,
	prop motionproto.Motion,
	args ...any,

) (motionproto.Report, notice.Notices) {

	notices := notice.Notices{}
	propStatePrev := motionapi.LoadPolicyState_Local[*pmp_0.ProposalState](ctx, cloned.PublicClone(), prop.ID)
	propState := propStatePrev.Copy()

	// update approval score

	ads := ballotapi.Show_Local(ctx, cloned.PublicClone(), propState.ApprovalPoll)
	latestApprovalScore := ads.Tally.Scores[pmp_0.ProposalBallotChoice]
	costOfReview := ads.Tally.Capitalization()
	propState.LatestApprovalScore = latestApprovalScore

	// update eligible concerns

	propState.EligibleConcerns = computeEligibleConcerns(ctx, cloned.PublicClone(), prop)

	projectedBounty := 0.0
	for _, ref := range propState.EligibleConcerns {
		conState := motionapi.LoadPolicyState_Local[*pmp_0.ConcernState](ctx, cloned.PublicClone(), ref.To)
		projectedBounty += conState.ProjectedBounty()
	}

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
				notice.Noticef(ctx, "The set of eligible issues claimed by this PR is now empty.\n")...,
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

	motionapi.SavePolicyState_StageOnly[*pmp_0.ProposalState](ctx, cloned.PublicClone(), prop.ID, propState)

	// update ScoreKernelState
	currentState := ScoreKernelState{
		Bounty: calcBounty(
			ctx,
			cloned,
			prop,
			propState,
		),
	}
	ballotapi.SavePolicyState_StageOnly[ScoreKernelState](
		ctx,
		cloned.PublicClone(),
		propState.ApprovalPoll,
		currentState,
	)

	return nil, notices
}

func calcBounty(
	ctx context.Context,
	cloned gov.OwnerCloned,
	prop motionproto.Motion,
	state *pmp_0.ProposalState,

) float64 {

	bounty := 0.0
	for _, ref := range state.EligibleConcerns {
		adt := ballotapi.Show_Local(ctx, cloned.PublicClone(), pmp_0.ConcernPollBallotName(ref.To))
		bounty += adt.Tally.Capitalization()
	}
	return bounty
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
	args ...any,

) (motionproto.Report, notice.Notices) {
	return nil, nil
}

func (x proposalPolicy) Close(
	ctx context.Context,
	cloned gov.OwnerCloned,
	prop motionproto.Motion,
	decision motionproto.Decision,
	args ...any,

) (motionproto.Report, notice.Notices) {

	// ensure the set of eligible concerns is valid
	_, uNotices := x.Update(ctx, cloned, prop)

	// was the PR merged or not
	isMerged := decision.IsAccept()

	approvalPollName := pmp_0.ProposalApprovalPollName(prop.ID)
	adt := loadPropApprovalPollTally(ctx, cloned.PublicClone(), prop)

	if isMerged {

		// accepting a proposal against the popular vote?
		againstPopular := adt.Tally.Scores[pmp_0.ProposalBallotChoice] < 0

		// close the referendum for the motion
		closeApprovalPoll := ballotapi.Close_StageOnly(
			ctx,
			cloned,
			approvalPollName,
			pmp_0.ProposalRewardAccountID(prop.ID),
		)

		// close all concerns resolved by the motion, and
		// transfer their escrows into the bounty account
		resolved := loadResolvedConcerns(ctx, cloned, prop)
		bounty := closeResolvedConcerns(ctx, cloned, prop, resolved)

		// transfer bounty to author
		var bountyDonated bool
		bountyReceipt := metric.Receipt{
			Type:   metric.ReceiptTypeBounty,
			Amount: bounty.MetricHolding(),
		}
		if prop.Author.IsNone() {
			account.Transfer_StageOnly(
				ctx,
				cloned.PublicClone(),
				pmp_0.ProposalBountyAccountID(prop.ID),
				pmp_0.MatchingPoolAccountID,
				bounty,
				fmt.Sprintf("bounty for proposal %v", prop.ID),
			)
			bountyDonated = true
			bountyReceipt.To = pmp_0.MatchingPoolAccountID.MetricAccountID()
		} else {
			account.Transfer_StageOnly(
				ctx,
				cloned.PublicClone(),
				pmp_0.ProposalBountyAccountID(prop.ID),
				member.UserAccountID(prop.Author),
				bounty,
				fmt.Sprintf("bounty for proposal %v", prop.ID),
			)
			bountyReceipt.To = member.UserAccountID(prop.Author).MetricAccountID()
		}

		// distribute rewards
		rewards := disberseRewards(ctx, cloned, prop)

		// metrics
		metric.Log_StageOnly(ctx, cloned.PublicClone(), &metric.Event{
			Motion: &metric.MotionEvent{
				Close: &metric.MotionClose{
					ID:       metric.MotionID(prop.ID),
					Type:     "proposal",
					Policy:   metric.MotionPolicy(prop.Policy),
					Decision: decision.MetricDecision(),
					Receipts: append(rewards.MetricReceipts(), bountyReceipt),
				},
			},
		})

		return &CloseReport{
				Accepted:            true,
				ApprovalPollOutcome: closeApprovalPoll.Result,
				Resolved:            resolved,
				Bounty:              bounty,
				BountyDonated:       bountyDonated,
				Rewarded:            rewards,
			}, append(
				uNotices,
				closeNotice(ctx, prop, againstPopular, closeApprovalPoll.Result, resolved, bounty, bountyDonated, rewards)...,
			)

	} else {

		// rejecting a proposal against the popular vote?
		againstPopular := adt.Tally.Scores[pmp_0.ProposalBallotChoice] > 0

		// cancel the referendum for the motion (refunds voters)
		cancelApprovalPoll := ballotapi.Cancel_StageOnly(
			ctx,
			cloned,
			approvalPollName,
		)

		// metrics
		metric.Log_StageOnly(ctx, cloned.PublicClone(), &metric.Event{
			Motion: &metric.MotionEvent{
				Close: &metric.MotionClose{
					ID:       metric.MotionID(prop.ID),
					Type:     "proposal",
					Policy:   metric.MotionPolicy(prop.Policy),
					Decision: decision.MetricDecision(),
					Receipts: cancelApprovalPoll.Result.RefundedHistoryReceipts(),
				},
			},
		})

		return &CloseReport{
				Accepted:            false,
				ApprovalPollOutcome: cancelApprovalPoll.Result,
				Resolved:            nil,
				Bounty:              account.H(account.PluralAsset, 0.0),
				BountyDonated:       false,
				Rewarded:            nil,
			}, append(
				uNotices,
				cancelNotice(ctx, prop, againstPopular, cancelApprovalPoll.Result)...,
			)
	}
}

func (x proposalPolicy) Cancel(
	ctx context.Context,
	cloned gov.OwnerCloned,
	prop motionproto.Motion,
	args ...any,

) (motionproto.Report, notice.Notices) {

	// cancel the referendum for the motion (and return credits to users)
	referendumName := pmp_0.ProposalApprovalPollName(prop.ID)
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
	}, notice.Noticef(ctx, "Cancelling management of this PR, managed as Gov4Git concern `%v`.", prop.ID)
}

type PolicyView struct {
	State          *pmp_0.ProposalState      `json:"state"`
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
	policyState := motionapi.LoadPolicyState_Local[*pmp_0.ProposalState](ctx, cloned, prop.ID)

	// retrieve approval poll
	approvalPoll := loadPropApprovalPollTally(ctx, cloned, prop)

	return PolicyView{
			State:          policyState,
			ApprovalPoll:   approvalPoll,
			ApprovalMargin: *ballotapi.GetMargin_Local(ctx, cloned, approvalPoll.Ad.ID),
			BountyAccount:  pmp_0.ProposalBountyAccountID(prop.ID),
			RewardAccount:  pmp_0.ProposalRewardAccountID(prop.ID),
		}, motionproto.MotionBallots{
			motionproto.MotionBallot{
				Label:         "approval_poll",
				BallotID:      policyState.ApprovalPoll,
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

	if refType != pmp_0.ClaimsRefType {
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

	if refType != pmp_0.ClaimsRefType {
		return nil, nil
	}

	return nil, notice.Noticef(ctx, "This PR no longer references %v, managed as Gov4Git concern `%v`.", to.TrackerURL, to.ID)
}

func (x proposalPolicy) Freeze(
	ctx context.Context,
	cloned gov.OwnerCloned,
	prop motionproto.Motion,
	args ...any,

) (motionproto.Report, notice.Notices) {

	return nil, notice.Noticef(ctx, "This PR, managed by Gov4Git proposal `%v`, has been frozen ❄️", prop.ID)
}

func (x proposalPolicy) Unfreeze(
	ctx context.Context,
	cloned gov.OwnerCloned,
	prop motionproto.Motion,
	args ...any,

) (motionproto.Report, notice.Notices) {

	return nil, notice.Noticef(ctx, "This PR, managed by Gov4Git proposal `%v`, has been unfrozen 🌤️", prop.ID)
}
