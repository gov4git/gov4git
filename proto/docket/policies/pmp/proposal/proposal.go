package proposal

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/proto/ballot/ballot"
	"github.com/gov4git/gov4git/proto/ballot/qv"
	"github.com/gov4git/gov4git/proto/docket/policies/pmp"
	"github.com/gov4git/gov4git/proto/docket/policy"
	"github.com/gov4git/gov4git/proto/docket/schema"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/member"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/ns"
)

func init() {
	policy.Install(context.Background(), proposalPolicy{})
}

const ProposalPolicyName = schema.PolicyName("pmp-proposal-policy")

type proposalPolicy struct{}

func (x proposalPolicy) Name() string {
	return ProposalPolicyName.String()
}

func (x proposalPolicy) Open(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion schema.Motion,
	policyNS ns.NS,

) {

	// initialize state
	state := NewProposalState(motion.ID)
	SaveState_StageOnly(ctx, cloned.Public.Tree(), policyNS, state)

	// open a poll for the motion
	ballot.Open_StageOnly(
		ctx,
		qv.QV{},
		cloned,
		state.ApprovalReferendum,
		fmt.Sprintf("Approval referendum for motion %v", motion.ID),
		fmt.Sprintf("Up/down vote the approval vote for proposal (pull request) %v", motion.ID),
		[]string{schema.MotionPollBallotChoice},
		member.Everybody,
	)
}

func (x proposalPolicy) Score(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion schema.Motion,
	policyNS ns.NS,

) schema.Score {

	state := LoadState_Local(ctx, cloned.Public.Tree(), policyNS)

	// compute score
	ads := ballot.Show_Local(ctx, cloned.Public.Tree(), state.ApprovalReferendum)
	attention := ads.Tally.Attention()

	return schema.Score{
		Attention: attention,
	}
}

func (x proposalPolicy) Close(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion schema.Motion,
	policyNS ns.NS,

) {

	// close the referendum for the motion
	referendumName := pmp.ProposalReferendumBallotName(motion.ID)
	ballot.Close_StageOnly(
		ctx,
		cloned,
		referendumName,
		false,
	)

	panic("XXX") // XXX: apply reward mechanism
}

func (x proposalPolicy) Cancel(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion schema.Motion,
	policyNS ns.NS,

) {

	// cancel the referendum for the motion (and return credits to users)
	referendumName := pmp.ProposalReferendumBallotName(motion.ID)
	ballot.Close_StageOnly(
		ctx,
		cloned,
		referendumName,
		true,
	)
}

func (x proposalPolicy) Show(
	ctx context.Context,
	cloned gov.Cloned,
	motion schema.Motion,
	policyNS ns.NS,

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
) {
}

func (x proposalPolicy) AddRefFrom(
	ctx context.Context,
	cloned gov.OwnerCloned,
	refType schema.RefType,
	from schema.Motion,
	to schema.Motion,
	fromPolicyNS ns.NS,
	toPolicyNS ns.NS,
) {
}

func (x proposalPolicy) RemoveRefTo(
	ctx context.Context,
	cloned gov.OwnerCloned,
	refType schema.RefType,
	from schema.Motion,
	to schema.Motion,
	fromPolicyNS ns.NS,
	toPolicyNS ns.NS,
) {
}

func (x proposalPolicy) RemoveRefFrom(
	ctx context.Context,
	cloned gov.OwnerCloned,
	refType schema.RefType,
	from schema.Motion,
	to schema.Motion,
	fromPolicyNS ns.NS,
	toPolicyNS ns.NS,
) {
}
