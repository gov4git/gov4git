package proposal

import (
	"bytes"
	"context"
	"fmt"
	"slices"

	"github.com/gov4git/gov4git/v2/proto/account"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballot"
	"github.com/gov4git/gov4git/v2/proto/ballot/common"
	"github.com/gov4git/gov4git/v2/proto/ballot/load"
	"github.com/gov4git/gov4git/v2/proto/docket/ops"
	"github.com/gov4git/gov4git/v2/proto/docket/policies/pmp"
	"github.com/gov4git/gov4git/v2/proto/docket/policy"
	"github.com/gov4git/gov4git/v2/proto/docket/schema"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/history"
	"github.com/gov4git/gov4git/v2/proto/member"
	"github.com/gov4git/gov4git/v2/proto/notice"
	"github.com/gov4git/gov4git/v2/proto/purpose"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/ns"
)

func init() {
	policy.Install(context.Background(), ProposalPolicyName, proposalPolicy{})
}

const ProposalPolicyName = schema.PolicyName("pmp-proposal-policy")

type proposalPolicy struct{}

func (x proposalPolicy) Open(
	ctx context.Context,
	cloned gov.OwnerCloned,
	prop schema.Motion,
	policyNS ns.NS,
	args ...any,

) (policy.Report, notice.Notices) {

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
	ballot.Open_StageOnly(
		ctx,
		load.QVStrategyName,
		cloned,
		state.ApprovalPoll,
		pmp.ProposalAccountID(prop.ID),
		purpose.Proposal,
		fmt.Sprintf("Approval referendum for motion %v", prop.ID),
		fmt.Sprintf("Up/down vote the approval vote for proposal (pull request) %v", prop.ID),
		[]string{pmp.ProposalBallotChoice},
		member.Everybody,
	)

	// metrics
	history.Log_StageOnly(ctx, cloned.PublicClone(), &history.Event{
		Motion: &history.MotionEvent{
			Open: &history.MotionOpen{
				ID:   history.MotionID(prop.ID),
				Type: "proposal",
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
	motion schema.Motion,
	policyNS ns.NS,
	args ...any,

) (schema.Score, notice.Notices) {

	state := LoadState_Local(ctx, cloned.Public.Tree(), policyNS)

	// compute score
	ads := ballot.Show_Local(ctx, cloned.Public.Tree(), state.ApprovalPoll)
	attention := ads.Tally.Attention()

	return schema.Score{
		Attention: attention,
	}, nil
}

func (x proposalPolicy) Update(
	ctx context.Context,
	cloned gov.OwnerCloned,
	prop schema.Motion,
	policyNS ns.NS,
	args ...any,

) (policy.Report, notice.Notices) {

	notices := notice.Notices{}
	state := LoadState_Local(ctx, cloned.Public.Tree(), policyNS)

	// update approval score

	ads := ballot.Show_Local(ctx, cloned.Public.Tree(), state.ApprovalPoll)
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
	if !slices.Equal[schema.Refs](eligible, state.EligibleConcerns) {
		// display list of eligible concerns
		if len(eligible) == 0 {
			notices = append(
				notices,
				notice.Noticef(ctx, "The set of eligible issues addressed by this PR is now empty.\n")...,
			)
		} else {
			var w bytes.Buffer
			for _, ref := range eligible {
				con := ops.LookupMotion_Local(ctx, cloned.PublicClone(), ref.To)
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

	return nil, notices
}

func (x proposalPolicy) Close(
	ctx context.Context,
	cloned gov.OwnerCloned,
	prop schema.Motion,
	policyNS ns.NS,
	decision schema.Decision,
	args ...any,

) (policy.Report, notice.Notices) {

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
		closeApprovalPoll := ballot.Close_StageOnly(
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
		bountyReceipt := history.Receipt{
			Type:   history.ReceiptTypeBounty,
			Amount: bounty.HistoryHolding(),
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
		history.Log_StageOnly(ctx, cloned.PublicClone(), &history.Event{
			Motion: &history.MotionEvent{
				Close: &history.MotionClose{
					ID:       history.MotionID(prop.ID),
					Type:     "proposal",
					Receipts: append(rewards.HistoryReceipts(), bountyReceipt),
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
		cancelApprovalPoll := ballot.Cancel_StageOnly(
			ctx,
			cloned,
			approvalPollName,
		)

		// metrics
		history.Log_StageOnly(ctx, cloned.PublicClone(), &history.Event{
			Motion: &history.MotionEvent{
				Close: &history.MotionClose{
					ID:       history.MotionID(prop.ID),
					Type:     "proposal",
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
	prop schema.Motion,
	policyNS ns.NS,
	args ...any,

) (policy.Report, notice.Notices) {

	// cancel the referendum for the motion (and return credits to users)
	referendumName := pmp.ProposalApprovalPollName(prop.ID)
	chg := ballot.Cancel_StageOnly(
		ctx,
		cloned,
		referendumName,
	)

	// metrics
	history.Log_StageOnly(ctx, cloned.PublicClone(), &history.Event{
		Motion: &history.MotionEvent{
			Cancel: &history.MotionCancel{
				ID:       history.MotionID(prop.ID),
				Type:     "proposals",
				Receipts: chg.Result.RefundedHistoryReceipts(),
			},
		},
	})

	return &CancelReport{
		ApprovalPollOutcome: chg.Result,
	}, notice.Noticef(ctx, "Cancelling management of this PR, managed as Gov4Git concern `%v`.", prop.ID)
}

type PolicyView struct {
	State         *ProposalState    `json:"state"`
	ApprovalPoll  common.AdTally    `json:"approval_poll"`
	BountyAccount account.AccountID `json:"bounty_account"`
	RewardAccount account.AccountID `json:"reward_account"`
}

func (x proposalPolicy) Show(
	ctx context.Context,
	cloned gov.Cloned,
	motion schema.Motion,
	policyNS ns.NS,
	args ...any,

) form.Form {

	// retrieve policy state
	policyState := LoadState_Local(ctx, cloned.Tree(), policyNS)

	// retrieve approval poll
	approvalPoll := loadPropApprovalPollTally(ctx, cloned, motion)

	return PolicyView{
		State:         policyState,
		ApprovalPoll:  approvalPoll,
		BountyAccount: pmp.ProposalBountyAccountID(motion.ID),
		RewardAccount: pmp.ProposalRewardAccountID(motion.ID),
	}
}

func (x proposalPolicy) AddRefTo(
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

func (x proposalPolicy) AddRefFrom(
	ctx context.Context,
	cloned gov.OwnerCloned,
	refType schema.RefType,
	from schema.Motion,
	to schema.Motion,
	fromPolicyNS ns.NS,
	toPolicyNS ns.NS,
	args ...any,

) (policy.Report, notice.Notices) {

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
	refType schema.RefType,
	from schema.Motion,
	to schema.Motion,
	fromPolicyNS ns.NS,
	toPolicyNS ns.NS,
	args ...any,

) (policy.Report, notice.Notices) {

	return nil, nil
}

func (x proposalPolicy) RemoveRefFrom(
	ctx context.Context,
	cloned gov.OwnerCloned,
	refType schema.RefType,
	from schema.Motion,
	to schema.Motion,
	fromPolicyNS ns.NS,
	toPolicyNS ns.NS,
	args ...any,

) (policy.Report, notice.Notices) {

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
	motion schema.Motion,
	policyNS ns.NS,
	args ...any,

) (policy.Report, notice.Notices) {

	return nil, notice.Noticef(ctx, "This PR, managed by Gov4Git proposal `%v`, has been frozen ❄️", motion.ID)
}

func (x proposalPolicy) Unfreeze(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion schema.Motion,
	policyNS ns.NS,
	args ...any,

) (policy.Report, notice.Notices) {

	return nil, notice.Noticef(ctx, "This PR, managed by Gov4Git proposal `%v`, has been unfrozen 🌤️", motion.ID)
}
