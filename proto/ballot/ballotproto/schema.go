package ballotproto

import (
	"sort"
	"time"

	"github.com/gov4git/gov4git/v2/proto/account"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/history/metric"
	"github.com/gov4git/gov4git/v2/proto/id"
	"github.com/gov4git/gov4git/v2/proto/member"
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

// VoterStatus reflects the state of an individual user's votes within a ballot.
type VoterStatus struct {
	GovID         id.ID             `json:"governance_id"`
	GovAddress    gov.Address       `json:"governance_address"`
	BallotID      BallotID          `json:"ballot_id"`
	AcceptedVotes AcceptedElections `json:"accepted_votes"`
	RejectedVotes RejectedElections `json:"rejected_votes"`
	PendingVotes  Elections         `json:"pending_votes"`
}

func VoteLogPath(govID id.ID, ballotName BallotID) ns.NS {
	return VoteLogNS.Append(
		form.StringHashForFilename(string(govID)),
		form.StringHashForFilename(BallotTopic(ballotName)),
	)
}

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

func OneElection(choice string, strength float64) Elections {
	return Elections{NewElection(choice, strength)}
}

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

type AdTally struct {
	Ad    Advertisement `json:"ballot_advertisement"`
	Tally Tally         `json:"ballot_tally"`
}

type Summary string

type Outcome struct {
	Summary      string                                      `json:"summary"`
	Scores       map[string]float64                          `json:"scores"`
	ScoresByUser map[member.User]map[string]StrengthAndScore `json:"scores_by_user"`
	Refunded     map[member.User]account.Holding             `json:"refunded"`
}

func (o Outcome) RefundedHistoryReceipts() metric.Receipts {
	r := metric.Receipts{}
	for user, h := range o.Refunded {
		r = append(r,
			metric.Receipt{
				To:     user.MetricAccountID(),
				Type:   metric.ReceiptTypeRefund,
				Amount: h.MetricHolding(),
			},
		)
	}
	return r
}

func FlattenRefunds(m map[member.User]account.Holding) Refunds {
	r := Refunds{}
	for k, v := range m {
		r = append(r, Refund{User: k, Amount: v})
	}
	r.Sort()
	return r
}

type Refund struct {
	User   member.User     `json:"user"`
	Amount account.Holding `json:"amount"`
}

type Refunds []Refund

func (x Refunds) Len() int {
	return len(x)
}

func (x Refunds) Less(i, j int) bool {
	return x[i].User < x[j].User
}

func (x Refunds) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

func (x Refunds) Sort() {
	sort.Sort(x)
}
