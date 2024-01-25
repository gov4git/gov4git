// Package pmp implements the Plural Management Protocol.
package pmp_0

import (
	"github.com/gov4git/gov4git/v2/proto/account"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
)

var (
	ConcernBallotChoice  = "rank"
	ProposalBallotChoice = "rank"

	// We explicitly avoid using "resolves" as the keyword for referencing issues/PRs, as
	// "resolves" triggers Github to automatically close resolved issues when a PR is closed, thereby
	// not giving Gov4Git a chance to close them as part of the PR closure/clearance procedure.
	ClaimsRefType = motionproto.RefType("claims")
)

func ConcernAccountID(id motionproto.MotionID) account.AccountID {
	return account.AccountIDFromLine(
		account.Cat(
			account.Pair("motion", id.String()),
			account.Term("pmp-concern"),
		),
	)
}

func ProposalAccountID(id motionproto.MotionID) account.AccountID {
	return account.AccountIDFromLine(
		account.Cat(
			account.Pair("motion", id.String()),
			account.Term("pmp-proposal"),
		),
	)
}

func ConcernPollBallotName(id motionproto.MotionID) ballotproto.BallotID {
	return ballotproto.BallotID("pmp/motion/priority_poll/" + id.String())
}

func ProposalApprovalPollName(id motionproto.MotionID) ballotproto.BallotID {
	return ballotproto.BallotID("pmp/motion/approval_poll/" + id.String())
}

func ProposalBountyAccountID(motionID motionproto.MotionID) account.AccountID {
	return account.AccountIDFromLine(
		account.Cat(
			account.Pair("motion", motionID.String()),
			account.Term("pmp-proposal-bounty"),
		),
	)
}

func ProposalRewardAccountID(motionID motionproto.MotionID) account.AccountID {
	return account.AccountIDFromLine(
		account.Cat(
			account.Pair("motion", motionID.String()),
			account.Term("pmp-proposal-reward"),
		),
	)
}
