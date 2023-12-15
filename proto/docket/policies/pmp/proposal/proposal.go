package proposal

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/proto/account"
	"github.com/gov4git/gov4git/proto/ballot/ballot"
	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/ballot/load"
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
	policy.Install(context.Background(), ProposalPolicyName, proposalPolicy{})
}

const ProposalPolicyName = schema.PolicyName("pmp-proposal-policy")

type proposalPolicy struct{}

func (x proposalPolicy) Open(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion schema.Motion,
	policyNS ns.NS,
	args ...any,

) notice.Notices {

	// initialize state
	state := NewProposalState(motion.ID)
	SaveState_StageOnly(ctx, cloned.Public.Tree(), policyNS, state)

	// create a bounty account for the proposal
	account.Create_StageOnly(
		ctx,
		cloned.PublicClone(),
		pmp.ProposalBountyAccountID(motion.ID),
		schema.MotionOwnerID(motion.ID),
	)

	// create a reward account for the proposal
	account.Create_StageOnly(
		ctx,
		cloned.PublicClone(),
		pmp.ProposalRewardAccountID(motion.ID),
		schema.MotionOwnerID(motion.ID),
	)

	// open a poll for the motion
	ballot.Open_StageOnly(
		ctx,
		load.QVStrategyName,
		cloned,
		state.ApprovalPoll,
		fmt.Sprintf("Approval referendum for motion %v", motion.ID),
		fmt.Sprintf("Up/down vote the approval vote for proposal (pull request) %v", motion.ID),
		[]string{pmp.ProposalBallotChoice},
		member.Everybody,
	)

	return notice.Noticef("Started managing this PR as Gov4Git proposal `%v`.", motion.ID)
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
	}, notice.Noticef("Updated approval tally to %v.", ads.Tally.Scores[pmp.ProposalBallotChoice])
}

func (x proposalPolicy) Update(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion schema.Motion,
	policyNS ns.NS,
	args ...any,

) notice.Notices {

	return nil
}

func (x proposalPolicy) Close(
	ctx context.Context,
	cloned gov.OwnerCloned,
	prop schema.Motion,
	policyNS ns.NS,
	args ...any,

) notice.Notices {

	// was the PR merged or not
	must.Assertf(ctx, len(args) == 1, "proposal closure missing argument")
	isMerged, ok := args[0].(bool) // isMerged
	must.Assertf(ctx, ok, "proposal closure unrecognized argument")

	referendumName := pmp.ProposalApprovalPollName(prop.ID)

	if !isMerged {

		// cancel the referendum for the motion (refunds voters)
		closeChg := ballot.Cancel_StageOnly(
			ctx,
			cloned,
			referendumName,
		)

		return cancelNotice(ctx, prop, closeChg.Result)

	} else {

		// close the referendum for the motion
		referendumName := pmp.ProposalApprovalPollName(prop.ID)
		closeChg := ballot.Close_StageOnly(
			ctx,
			cloned,
			referendumName,
			pmp.ProposalRewardAccountID(prop.ID),
		)

		// close all concerns resolved by the motion, and
		// transfer their escrows into the bounty account
		resolved := loadResolvedConcerns(ctx, cloned, prop)
		bounty := closeResolvedConcerns(ctx, cloned, prop, resolved)

		// transfer bounty to author
		account.Transfer_StageOnly(
			ctx,
			cloned.PublicClone(),
			pmp.ProposalBountyAccountID(prop.ID),
			member.UserAccountID(prop.Author),
			bounty,
		)

		// distribute rewards
		rewards := disberseRewards(ctx, cloned, prop)

		return closeNotice(ctx, prop, closeChg.Result, resolved, bounty, rewards)
	}
}

func (x proposalPolicy) Cancel(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion schema.Motion,
	policyNS ns.NS,
	args ...any,

) notice.Notices {

	// cancel the referendum for the motion (and return credits to users)
	referendumName := pmp.ProposalApprovalPollName(motion.ID)
	ballot.Cancel_StageOnly(
		ctx,
		cloned,
		referendumName,
	)

	return notice.Noticef("Cancelling management of this PR, managed as Gov4Git concern `%v`.", motion.ID)
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

) notice.Notices {

	return nil
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

) notice.Notices {

	return nil
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

) notice.Notices {

	return nil
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

) notice.Notices {

	return nil
}

func (x proposalPolicy) Freeze(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion schema.Motion,
	policyNS ns.NS,
	args ...any,

) notice.Notices {

	return notice.Noticef("This PR, managed by Gov4Git proposal `%v`, has been frozen ❄️", motion.ID)
}

func (x proposalPolicy) Unfreeze(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion schema.Motion,
	policyNS ns.NS,
	args ...any,

) notice.Notices {

	return notice.Noticef("This PR, managed by Gov4Git proposal `%v`, has been unfrozen 🌤️", motion.ID)
}
