package common

import (
	"time"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/account"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/gov4git/proto/member"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/ns"
	"github.com/gov4git/lib4git/util"
)

var (
	BallotNS         = proto.RootNS.Append("ballot")
	AdFilebase       = "ballot_ad.json"
	StrategyFilebase = "ballot_strategy.json"
	TallyFilebase    = "ballot_tally.json"
	OutcomeFilebase  = "ballot_outcome.json"

	VoteLogNS = proto.RootNS.Append("votes") // namespace in voter's repo for recording votes
)

// VoteLog records the votes of a user to a ballot within a given governance.
type VoteLog struct {
	GovID         id.ID         `json:"governance_id"`
	GovAddress    gov.Address   `json:"governance_address"`
	Ballot        BallotName    `json:"ballot_name"`
	VoteEnvelopes VoteEnvelopes `json:"vote_envelopes"` // in the order in which they were sent
}

// VoterStatus reflects the state of an individual user's votes within a ballot.
type VoterStatus struct {
	GovID         id.ID             `json:"governance_id"`
	GovAddress    gov.Address       `json:"governance_address"`
	BallotName    BallotName        `json:"ballot_name"`
	AcceptedVotes AcceptedElections `json:"accepted_votes"`
	RejectedVotes RejectedElections `json:"rejected_votes"`
	PendingVotes  Elections         `json:"pending_votes"`
}

func VoteLogPath(govID id.ID, ballotName BallotName) ns.NS {
	return VoteLogNS.Append(
		form.StringHashForFilename(string(govID)),
		form.StringHashForFilename(BallotTopic(ballotName)),
	)
}

func BallotOwnerID(ballotName BallotName) account.OwnerID {
	return account.OwnerIDFromLine(account.Pair("ballot", ballotName.GitPath()))
}

func BallotEscrowAccountID(ballotName BallotName) account.AccountID {
	return account.AccountIDFromLine(account.Pair("ballot_escrow", ballotName.GitPath()))
}

func BallotTopic(ballotName BallotName) string {
	// BallotTopic must produce the same string on every OS.
	// It is essential to use ballotName.GitPath, instead of ballotName.Path which is OS-specific.
	return "ballot:" + ballotName.GitPath()
}

func BallotPath(name BallotName) ns.NS {
	return BallotNS.Join(name.NS())
}

type BallotName ns.NS

func (x BallotName) OSPath() string {
	return ns.NS(x).OSPath()
}

func (x BallotName) GitPath() string {
	return ns.NS(x).GitPath()
}

func (x BallotName) NS() ns.NS {
	return ns.NS(x)
}

func ParseBallotNameFromPath(p string) BallotName {
	return BallotName(ns.ParseFromGitPath(p))
}

type Advertisement struct {
	Gov          gov.Address    `json:"community"`
	Name         BallotName     `json:"name"`
	Title        string         `json:"title"`
	Description  string         `json:"description"`
	Choices      []string       `json:"choices"`
	Strategy     string         `json:"strategy"`
	Participants member.Group   `json:"participants_group"`
	Frozen       bool           `json:"frozen"` // if frozen, the ballot is not accepting votes
	Closed       bool           `json:"closed"` // closed ballots cannot be re-opened
	Cancelled    bool           `json:"cancelled"`
	ParentCommit git.CommitHash `json:"parent_commit"`
}

type BallotAddress struct {
	Gov  gov.Address
	Name BallotName
}

type Election struct {
	VoteID             id.ID     `json:"vote_id"`
	VoteTime           time.Time `json:"vote_time"`
	VoteChoice         string    `json:"vote_choice"`
	VoteStrengthChange float64   `json:"vote_strength_change"`
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

type VoteEnvelope struct {
	AdCommit  git.CommitHash `json:"ballot_ad_commit"`
	Ad        Advertisement  `json:"ballot_ad"`
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

type AdStrategyTally struct {
	Ad       Advertisement `json:"ballot_advertisement"`
	Strategy Strategy      `json:"ballot_strategy"`
	Tally    Tally         `json:"ballot_tally"`
}

type Summary string

type Outcome struct {
	Summary string             `json:"ballot_summary"`
	Scores  map[string]float64 `json:"ballot_scores"`
}
