// Package pmp implements the Plural Management Protocol.
package pmp

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/proto/account"
	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/docket/schema"
	"github.com/gov4git/gov4git/proto/gov"
)

var (
	ConcernBallotChoice  = "rank"
	ProposalBallotChoice = "rank"

	// We chose not to use "resolves", because it triggers GitHub to close the resolved issue immediately.
	ResolvesRefType = schema.RefType("addresses")
)

func ConcernPollBallotName(id schema.MotionID) common.BallotName {
	return common.BallotName{"pmp", "motion", "priority_poll", id.String()}
}

func ProposalApprovalPollName(id schema.MotionID) common.BallotName {
	return common.BallotName{"pmp", "motion", "approval_poll", id.String()}
}

func ProposalBountyAccountID(motionID schema.MotionID) account.AccountID {
	return account.AccountIDFromLine(
		account.Cat(
			account.Pair("motion", motionID.String()),
			account.Term("pmp-proposal-bounty"),
		),
	)
}

func ProposalRewardAccountID(motionID schema.MotionID) account.AccountID {
	return account.AccountIDFromLine(
		account.Cat(
			account.Pair("motion", motionID.String()),
			account.Term("pmp-proposal-reward"),
		),
	)
}

var (
	OwnerID = account.AccountIDFromLine(
		account.Term("pmp"),
	)
	BurnPoolAccountID = account.AccountIDFromLine(
		account.Cat(
			account.Term("pmp"),
			account.Term("burn"),
		),
	)
	TaxPoolAccountID = account.AccountIDFromLine(
		account.Cat(
			account.Term("pmp"),
			account.Term("tax"),
		),
	)
	MatchingPoolAccountID = account.AccountIDFromLine(
		account.Cat(
			account.Term("pmp"),
			account.Term("matching"),
		),
	)
)

func Boot_StageOnly(ctx context.Context, cloned gov.Cloned) {

	// create burn pool account
	account.Create_StageOnly(
		ctx,
		cloned,
		BurnPoolAccountID,
		account.OwnerID(OwnerID),
		fmt.Sprintf("burn account for PMP"),
	)

	// create tax pool account
	account.Create_StageOnly(
		ctx,
		cloned,
		TaxPoolAccountID,
		account.OwnerID(OwnerID),
		fmt.Sprintf("tax account for PMP"),
	)

	// create matching pool account
	account.Create_StageOnly(
		ctx,
		cloned,
		MatchingPoolAccountID,
		account.OwnerID(OwnerID),
		fmt.Sprintf("matching pool account for PMP"),
	)

}
