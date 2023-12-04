package proposal

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/proto/account"
	"github.com/gov4git/gov4git/proto/ballot/ballot"
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

	// create an account for the proposal
	account.Create_StageOnly(
		ctx,
		cloned.PublicClone(),
		schema.PolicyAccountID(motion.ID),
		schema.MotionOwnerID(motion.ID),
	)

	// open a poll for the motion
	ballot.Open_StageOnly(
		ctx,
		load.QVStrategyName,
		cloned,
		state.ApprovalReferendum,
		fmt.Sprintf("Approval referendum for motion %v", motion.ID),
		fmt.Sprintf("Up/down vote the approval vote for proposal (pull request) %v", motion.ID),
		[]string{schema.MotionPollBallotChoice},
		member.Everybody,
	)

	return notice.Noticef("Started managing this PR as Gov4Git proposal #%v.", motion.ID)
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
	ads := ballot.Show_Local(ctx, cloned.Public.Tree(), state.ApprovalReferendum)
	attention := ads.Tally.Attention()

	return schema.Score{
		Attention: attention,
	}, notice.Noticef("Updated approval tally to %v.", ads.Tally.Scores[schema.MotionPollBallotChoice])
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
	motion schema.Motion,
	policyNS ns.NS,
	args ...any,

) notice.Notices {

	must.Assertf(ctx, len(args) == 1, "proposal closure missing argument")
	_, ok := args[0].(bool) // isMerged
	must.Assertf(ctx, ok, "proposal closure unrecognized argument")

	// close the referendum for the motion
	referendumName := pmp.ProposalReferendumBallotName(motion.ID)
	ballot.Close_StageOnly(
		ctx,
		cloned,
		referendumName,
		false,
	)

	// XXX: apply reward mechanism

	return notice.Noticef("Closing managment of this PR, managed as Gov4Git proposal #%v).", motion.ID)
}

func (x proposalPolicy) Cancel(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion schema.Motion,
	policyNS ns.NS,
	args ...any,

) notice.Notices {

	// cancel the referendum for the motion (and return credits to users)
	referendumName := pmp.ProposalReferendumBallotName(motion.ID)
	ballot.Close_StageOnly(
		ctx,
		cloned,
		referendumName,
		true,
	)

	return notice.Noticef("Cancelling management of this PR, managed as Gov4Git concern #%v.", motion.ID)
}

func (x proposalPolicy) Show(
	ctx context.Context,
	cloned gov.Cloned,
	motion schema.Motion,
	policyNS ns.NS,
	args ...any,

) form.Map {

	// retrieve policy state
	policyState := LoadState_Local(ctx, cloned.Tree(), policyNS)

	// retrieve referendum state
	referendumName := pmp.ProposalReferendumBallotName(motion.ID)
	referendumState := ballot.Show_Local(ctx, cloned.Tree(), referendumName)

	return form.Map{
		"pmp_proposal_policy_state":              policyState,
		"pmp_proposal_approval_referendum_state": referendumState,
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

	return notice.Noticef("This PR, managed by Gov4Git proposal #%v, has been frozen.", motion.ID)
}

func (x proposalPolicy) Unfreeze(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion schema.Motion,
	policyNS ns.NS,
	args ...any,

) notice.Notices {

	return notice.Noticef("This PR, managed by Gov4Git proposal #%v, has been unfrozen.", motion.ID)
}
