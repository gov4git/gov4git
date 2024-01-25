package ballotproto

import (
	"time"

	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/id"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/ns"
	"github.com/gov4git/lib4git/util"
)

// VoteLog records the votes of a user to a ballot within a given governance.
type VoteLog struct {
	GovID         id.ID         `json:"governance_id"`
	GovAddress    gov.Address   `json:"governance_address"`
	BallotID      BallotID      `json:"ballot_id"`
	VoteEnvelopes VoteEnvelopes `json:"vote_envelopes"` // in the order in which they were sent
}

func VoteLogPath(govID id.ID, ballotName BallotID) ns.NS {
	return VoteLogNS.Append(
		form.StringHashForFilename(string(govID)),
		form.StringHashForFilename(BallotTopic(ballotName)),
	)
}

// VoterStatus reflects the state of an individual user's votes within a ballot.
type VoterStatus struct {
	GovID         id.ID             `json:"governance_id"`
	GovAddress    gov.Address       `json:"governance_address"`
	BallotID      BallotID          `json:"ballot_id"`
	AcceptedVotes AcceptedElections `json:"accepted_votes"`
	RejectedVotes RejectedElections `json:"rejected_votes"`
	PendingVotes  Elections         `json:"pending_votes"`
}

type Election struct {
	VoteID             id.ID     `json:"vote_id"`
	VoteTime           time.Time `json:"vote_time"`
	VoteChoice         string    `json:"vote_choice"`
	VoteStrengthChange float64   `json:"vote_strength_change"` // this is the voter's payment with a sign to indicate direction of vote
}

func NewElection(choice string, strength float64) Election {
	return Election{
		VoteID:             id.GenerateRandomID(),
		VoteTime:           time.Now(),
		VoteChoice:         choice,
		VoteStrengthChange: strength,
	}
}

type Elections []Election

func OneElection(choice string, strength float64) Elections {
	return Elections{NewElection(choice, strength)}
}

type VoteEnvelope struct {
	AdCommit  git.CommitHash `json:"ballot_ad_commit"`
	Ad        Ad             `json:"ballot_ad"`
	Elections Elections      `json:"ballot_elections"`
}

type VoteEnvelopes []VoteEnvelope

// Verify verifies that elections are consistent with the ballot ad.
func (x VoteEnvelope) VerifyConsistency() bool {
	for _, v := range x.Elections {
		if !util.IsIn(v.VoteChoice, x.Ad.Choices...) {
			return false
		}
	}
	return true
}
