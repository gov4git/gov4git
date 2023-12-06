// Package pmp implements the Plural Management Protocol.
package pmp

import (
	"context"

	"github.com/gov4git/gov4git/proto/account"
	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/docket/schema"
	"github.com/gov4git/gov4git/proto/gov"
)

var (
	ConcernBallotChoice  = "priority"
	ProposalBallotChoice = "approval"
)

func ConcernPollBallotName(id schema.MotionID) common.BallotName {
	return common.BallotName{"pmp", "motion", "poll", id.String()}
}

func ProposalApprovalPollName(id schema.MotionID) common.BallotName {
	return common.BallotName{"pmp", "motion", "referendum", id.String()}
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
	)

	// create tax pool account
	account.Create_StageOnly(
		ctx,
		cloned,
		TaxPoolAccountID,
		account.OwnerID(OwnerID),
	)

	// create matching pool account
	account.Create_StageOnly(
		ctx,
		cloned,
		MatchingPoolAccountID,
		account.OwnerID(OwnerID),
	)

}
