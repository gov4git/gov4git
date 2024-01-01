package proposal

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
	"github.com/gov4git/gov4git/v2/proto/motion/motionapi"
	"github.com/gov4git/gov4git/v2/proto/motion/motionpolicies/pmp"
	"github.com/gov4git/gov4git/v2/proto/motion/motionpolicy"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
	"github.com/gov4git/gov4git/v2/proto/notice"
	"github.com/gov4git/gov4git/v2/proto/purpose"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/ns"
)

func init() {
	motionpolicy.Install(context.Background(), ProposalPolicyName, proposalPolicy{})
}

const (
	ProposalPolicyName               motionproto.PolicyName   = "pmp-proposal"
	ProposalApprovalPollStrategyName ballotproto.StrategyName = "pmp-proposal-approval"
)

type proposalPolicy struct{}

func (x proposalPolicy) Open(
	ctx context.Context,
	cloned gov.OwnerCloned,
	prop motionproto.Motion,
	policyNS ns.NS,
	args ...any,

) (motionpolicy.Report, notice.Notices) {

	// initialize state
	state := NewProposalState(prop.ID)
	SaveState_StageOnly(ctx, cloned.Public.Tree(), policyNS, state)

	// create a bounty account for the proposal
	account.Create_StageOnly(
		ctx,
		cloned.PublicClone(),
		pmp.ProposalBountyAccountID(prop.ID),
		pmp.ProposalAccountID(prop.ID),
		fmt.Sprintf("bounty account for proposal %v", prop.ID),
	)

	// create a reward account for the proposal
	account.Create_StageOnly(
		ctx,
		cloned.PublicClone(),
		pmp.ProposalRewardAccountID(prop.ID),
		pmp.ProposalAccountID(prop.ID),
		fmt.Sprintf("reward account for proposal %v", prop.ID),
	)

	// open a poll for the motion
	ballotapi.Open_StageOnly(
		ctx,
		ballotio.QVStrategyName,
		cloned,
		state.ApprovalPoll,
		pmp.ProposalAccountID(prop.ID),
		purpose.Proposal,
		prop.Policy,
		fmt.Sprintf("Approval referendum for motion %v", prop.ID),
		fmt.Sprintf("Up/down vote the approval vote for proposal (pull request) %v", prop.ID),
		[]string{pmp.ProposalBallotChoice},
		member.Everybody,
	)
	zeroState := ScoreKernelState{
		Bounty: 0.0,
	}
	ballotapi.SaveStrategyState_StageOnly[ScoreKernelState](
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
			pmp.Welcome, prop.ID, state.LatestApprovalScore)
}

func (x proposalPolicy) Score(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion motionproto.Motion,
	policyNS ns.NS,
	args ...any,

) (motionproto.Score, notice.Notices) {

	state := LoadState_Local(ctx, cloned.Public.Tree(), policyNS)

	// compute score
	ads := ballotapi.Show_Local(ctx, cloned.Public.Tree(), state.ApprovalPoll)
	attention := ads.Tally.Attention()

	return motionproto.Score{
		Attention: attention,
	}, nil
}

func (x proposalPolicy) Update(
	ctx context.Context,
	cloned gov.OwnerCloned,
	prop motionproto.Motion,
	policyNS ns.NS,
	args ...any,

) (motionpolicy.Report, notice.Notices) {

	notices := notice.Notices{}
	state := LoadState_Local(ctx, cloned.Public.Tree(), policyNS)

	// update approval score

	ads := ballotapi.Show_Local(ctx, cloned.Public.Tree(), state.ApprovalPoll)
	latestApprovalScore := ads.Tally.Scores[pmp.ProposalBallotChoice]
	if latestApprovalScore != state.LatestApprovalScore {
		notices = append(
			notices,
			notice.Noticef(ctx, "This PR's __approval score__ was updated to `%0.6f`.", latestApprovalScore)...,
		)
	}
	state.LatestApprovalScore = latestApprovalScore

	// update eligible concerns

	eligible := computeEligibleConcerns(ctx, cloned.PublicClone(), prop)
	if !slices.Equal[motionproto.Refs](eligible, state.EligibleConcerns) {
		// display list of eligible concerns
		if len(eligible) == 0 {
			notices = append(
				notices,
				notice.Noticef(ctx, "The set of eligible issues addressed by this PR is now empty.\n")...,
			)
		} else {
			var w bytes.Buffer
			for _, ref := range eligible {
				con := motionapi.LookupMotion_Local(ctx, cloned.PublicClone(), ref.To)
				fmt.Fprintf(&w, "- %s, managed as Gov4Git motion `%v` with community attention of `%0.6f`\n",
					con.TrackerURL, con.ID, con.Score.Attention)
			}
			notices = append(
				notices,
				notice.Noticef(ctx, "The set of eligible issues addressed by this PR changed to:\n"+w.String())...,
			)
		}
	}
	state.EligibleConcerns = eligible

	//

	SaveState_StageOnly(ctx, cloned.Public.Tree(), policyNS, state)

	// update ScoreKernelState
	currentState := ScoreKernelState{
		Bounty: calcBounty(
			ctx,
			cloned,
			prop,
			state,
		),
	}
	ballotapi.SaveStrategyState_StageOnly[ScoreKernelState](
		ctx,
		cloned.PublicClone(),
		state.ApprovalPoll,
		currentState,
	)

	return nil, notices
}

func calcBounty(
	ctx context.Context,
	cloned gov.OwnerCloned,
	prop motionproto.Motion,
	state *ProposalState,
) float64 {

	bounty := 0.0
	for _, ref := range state.EligibleConcerns {
		adt := ballotapi.Show_Local(ctx, cloned.PublicClone().Tree(), pmp.ConcernPollBallotName(ref.To))
		bounty += adt.Tally.Capitalization()
	}
	return bounty
}

func (x proposalPolicy) Close(
	ctx context.Context,
	cloned gov.OwnerCloned,
	prop motionproto.Motion,
	policyNS ns.NS,
	decision motionproto.Decision,
	args ...any,

) (motionpolicy.Report, notice.Notices) {

	// update the policy state before closing the motion
	x.Update(ctx, cloned, prop, policyNS)

	// was the PR merged or not
	isMerged := decision.IsAccept()

	approvalPollName := pmp.ProposalApprovalPollName(prop.ID)
	adt := loadPropApprovalPollTally(ctx, cloned.PublicClone(), prop)

	if isMerged {

		// accepting a proposal against the popular vote?
		againstPopular := adt.Tally.Scores[pmp.ProposalBallotChoice] < 0

		// close the referendum for the motion
		approvalPollName := pmp.ProposalApprovalPollName(prop.ID)
		closeApprovalPoll := ballotapi.Close_StageOnly(
			ctx,
			cloned,
			approvalPollName,
			pmp.ProposalRewardAccountID(prop.ID),
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
				pmp.ProposalBountyAccountID(prop.ID),
				pmp.MatchingPoolAccountID,
				bounty,
				fmt.Sprintf("bounty for proposal %v", prop.ID),
			)
			bountyDonated = true
			bountyReceipt.To = pmp.MatchingPoolAccountID.HistoryAccountID()
		} else {
			account.Transfer_StageOnly(
				ctx,
				cloned.PublicClone(),
				pmp.ProposalBountyAccountID(prop.ID),
				member.UserAccountID(prop.Author),
				bounty,
				fmt.Sprintf("bounty for proposal %v", prop.ID),
			)
			bountyReceipt.To = member.UserAccountID(prop.Author).HistoryAccountID()
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
		}, closeNotice(ctx, prop, againstPopular, closeApprovalPoll.Result, resolved, bounty, bountyDonated, rewards)

	} else {

		// rejecting a proposal against the popular vote?
		againstPopular := adt.Tally.Scores[pmp.ProposalBallotChoice] > 0

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
		}, cancelNotice(ctx, prop, againstPopular, cancelApprovalPoll.Result)

	}
}

func (x proposalPolicy) Cancel(
	ctx context.Context,
	cloned gov.OwnerCloned,
	prop motionproto.Motion,
	policyNS ns.NS,
	args ...any,

) (motionpolicy.Report, notice.Notices) {

	// cancel the referendum for the motion (and return credits to users)
	referendumName := pmp.ProposalApprovalPollName(prop.ID)
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
	State          *ProposalState      `json:"state"`
	ApprovalPoll   ballotproto.AdTally `json:"approval_poll"`
	ApprovalMargin ballotproto.Margin  `json:"priority_margin"`
	BountyAccount  account.AccountID   `json:"bounty_account"`
	RewardAccount  account.AccountID   `json:"reward_account"`
}

func (x proposalPolicy) Show(
	ctx context.Context,
	cloned gov.Cloned,
	motion motionproto.Motion,
	policyNS ns.NS,
	args ...any,

) form.Form {

	// retrieve policy state
	policyState := LoadState_Local(ctx, cloned.Tree(), policyNS)

	// retrieve approval poll
	approvalPoll := loadPropApprovalPollTally(ctx, cloned, motion)

	return PolicyView{
		State:          policyState,
		ApprovalPoll:   approvalPoll,
		ApprovalMargin: *ballotapi.GetMargin_Local(ctx, cloned, approvalPoll.Ad.ID),
		BountyAccount:  pmp.ProposalBountyAccountID(motion.ID),
		RewardAccount:  pmp.ProposalRewardAccountID(motion.ID),
	}
}

func (x proposalPolicy) AddRefTo(
	ctx context.Context,
	cloned gov.OwnerCloned,
	refType motionproto.RefType,
	from motionproto.Motion,
	to motionproto.Motion,
	fromPolicyNS ns.NS,
	toPolicyNS ns.NS,
	args ...any,

) (motionpolicy.Report, notice.Notices) {

	return nil, nil
}

func (x proposalPolicy) AddRefFrom(
	ctx context.Context,
	cloned gov.OwnerCloned,
	refType motionproto.RefType,
	from motionproto.Motion,
	to motionproto.Motion,
	fromPolicyNS ns.NS,
	toPolicyNS ns.NS,
	args ...any,

) (motionpolicy.Report, notice.Notices) {

	if !to.IsConcern() {
		return nil, nil
	}

	if refType != pmp.ClaimsRefType {
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
	fromPolicyNS ns.NS,
	toPolicyNS ns.NS,
	args ...any,

) (motionpolicy.Report, notice.Notices) {

	return nil, nil
}

func (x proposalPolicy) RemoveRefFrom(
	ctx context.Context,
	cloned gov.OwnerCloned,
	refType motionproto.RefType,
	from motionproto.Motion,
	to motionproto.Motion,
	fromPolicyNS ns.NS,
	toPolicyNS ns.NS,
	args ...any,

) (motionpolicy.Report, notice.Notices) {

	if !to.IsConcern() {
		return nil, nil
	}

	if refType != pmp.ClaimsRefType {
		return nil, nil
	}

	return nil, notice.Noticef(ctx, "This PR no longer references %v, managed as Gov4Git concern `%v`.", to.TrackerURL, to.ID)
}

func (x proposalPolicy) Freeze(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion motionproto.Motion,
	policyNS ns.NS,
	args ...any,

) (motionpolicy.Report, notice.Notices) {

	return nil, notice.Noticef(ctx, "This PR, managed by Gov4Git proposal `%v`, has been frozen ❄️", motion.ID)
}

func (x proposalPolicy) Unfreeze(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion motionproto.Motion,
	policyNS ns.NS,
	args ...any,

) (motionpolicy.Report, notice.Notices) {

	return nil, notice.Noticef(ctx, "This PR, managed by Gov4Git proposal `%v`, has been unfrozen 🌤️", motion.ID)
}
