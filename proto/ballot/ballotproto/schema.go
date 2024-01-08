package ballotproto

import (
	"github.com/gov4git/gov4git/v2/proto/account"
	"github.com/gov4git/gov4git/v2/proto/gov"
)

func BallotEscrowAccountID(ballotName BallotID) account.AccountID {
	return account.AccountIDFromLine(account.Pair("ballot_escrow", ballotName.GitPath()))
}

func BallotTopic(ballotName BallotID) string {
	// BallotTopic must produce the same string on every OS.
	// It is essential to use ballotName.GitPath, instead of ballotName.Path which is OS-specific.
	return "ballot:" + ballotName.GitPath()
}

type BallotAddress struct {
	Gov  gov.Address
	Name BallotID
}

type AdTally struct {
	Ad    Ad    `json:"ballot_advertisement"`
	Tally Tally `json:"ballot_tally"`
}
