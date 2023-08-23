package collab

import (
	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/kv"
)

var (
	collabNS = proto.RootNS.Sub("collab")

	concernNS = collabNS.Sub("concern")
	concernKV = kv.KV[ConcernName, Concern]{}

	proposalNS = collabNS.Sub("proposal")
	proposalKV = kv.KV[ProposalName, Proposal]{}
)

type Name string

func (x Name) String() string {
	return string(x)
}

// ConcernName is the name of a concern within the gov4git system.
// ConcernNames must be unique within a community.
type ConcernName = Name

// Concern is the current state of an concern.
type Concern struct {
	Name        ConcernName    `json:"name"` // name of concern
	Title       string         `json:"title"`
	TrackerURL  string         `json:"tracker_url"` // link to concern on an external concern tracker, such as a GitHub issue
	Closed      bool           `json:"closed"`
	Cancelled   bool           `json:"cancelled"`
	Priority    Priority       `json:"priority"`
	AddressedBy []ProposalName `json:"addressed_by"` // prs addressing this concern
}

func ConcernPriorityBallotName(concernName ConcernName) common.BallotName {
	return common.BallotName{"concern", "priority", concernName.String()}
}

// ProposalName is the name of a pull request within the gov4git system.
// ProposalNames must be unique within a community.
type ProposalName = Name

// Proposal is the current state of a pull request.
type Proposal struct {
	Name       ProposalName  `json:"name"` // name of proposal
	Title      string        `json:"title"`
	TrackerURL string        `json:"tracker_url"` // link to proposal on an external proposal tracker, such as GitHub
	Closed     bool          `json:"closed"`
	Cancelled  bool          `json:"cancelled"`
	Priority   Priority      `json:"priority"`
	Addresses  []ConcernName `json:"addresses"` // concerns addressed by this proposal
}

func ProposalPriorityBallotName(proposalName ProposalName) common.BallotName {
	return common.BallotName{"proposal", "priority", proposalName.String()}
}

// ProposalConcernPair is a pair of a proposal and a concern.
type ProposalConcernPair struct {
	Proposal Proposal `json:"proposal"`
	Concern  Concern  `json:"concern"`
}

// Priority describes how a concern or a proposal is prioritized.
type Priority struct {
	Fixed  *float64           `json:"fixed"`
	Ballot *common.BallotName `json:"ballot"`
}
