package alcap

import (
	"github.com/gov4git/gov4git/v2/proto/account"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/motion"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
)

const (
	ConcernBallotChoice  = "rank"
	ProposalBallotChoice = "rank"

	ConcernPolicyName  motion.PolicyName = "al-capitan-concern"
	ProposalPolicyName motion.PolicyName = "al-capitan-proposal"

	ConcernPolicyGithubLabel  = "gov4git:al-cap" //XXX: add to github driver
	ProposalPolicyGithubLabel = ConcernPolicyGithubLabel

	ClaimsRefType = "claims"
)

func ConcernAccountID(id motionproto.MotionID) account.AccountID {
	return account.AccountIDFromLine(
		account.Cat(
			account.Pair("motion", id.String()),
			account.Term("al-cap-concern"),
		),
	)
}

func ProposalAccountID(id motionproto.MotionID) account.AccountID {
	return account.AccountIDFromLine(
		account.Cat(
			account.Pair("motion", id.String()),
			account.Term("al-cap-proposal"),
		),
	)
}

func ConcernPollBallotName(id motionproto.MotionID) ballotproto.BallotID {
	return ballotproto.BallotID("al-cap/motion/priority_poll/" + id.String())
}

func ProposalApprovalPollName(id motionproto.MotionID) ballotproto.BallotID {
	return ballotproto.BallotID("al-cap/motion/approval_poll/" + id.String())
}

func ProposalBountyAccountID(motionID motionproto.MotionID) account.AccountID {
	return account.AccountIDFromLine(
		account.Cat(
			account.Pair("motion", motionID.String()),
			account.Term("al-cap-proposal-bounty"),
		),
	)
}

func ProposalRewardAccountID(motionID motionproto.MotionID) account.AccountID {
	return account.AccountIDFromLine(
		account.Cat(
			account.Pair("motion", motionID.String()),
			account.Term("al-cap-proposal-reward"),
		),
	)
}
