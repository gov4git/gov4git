// Package pmp implements the Plural Management Protocol.
package pmp

import (
	"context"

	"github.com/gov4git/gov4git/proto/account"
	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/docket/schema"
	"github.com/gov4git/gov4git/proto/gov"
)

func ConcernPollBallotName(id schema.MotionID) common.BallotName {
	return common.BallotName{"pmp", "motion", "poll", id.String()}
}

func ProposalReferendumBallotName(id schema.MotionID) common.BallotName {
	return common.BallotName{"pmp", "motion", "referendum", id.String()}
}

var (
	OwnerID = account.AccountIDFromLine(
		account.Term("pmp"),
	)
	TaxPoolAccountID = account.AccountIDFromLine(
		account.Cat(
			account.Term("pmp"),
			account.Term("tax_pool"),
		),
	)
	MatchingPoolAccountID = account.AccountIDFromLine(
		account.Cat(
			account.Term("pmp"),
			account.Term("matching_pool"),
		),
	)
)

func Boot_StageOnly(ctx context.Context, cloned gov.Cloned) {

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
