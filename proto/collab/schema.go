package collab

import (
	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/kv"
)

var (
	collabNS = proto.RootNS.Sub("collab")

	concernNS = collabNS.Sub("concern")
	concernKV = kv.KV[ConcernName, ConcernState]{}

	prNS = collabNS.Sub("proposal")
	prKV = kv.KV[ConcernName, ProposalState]{}
)

// ConcernName is the name of an issue within the gov4git system.
// ConcernNames must be unique within a community.
type ConcernName string

// ConcernState is the current state of an issue.
type ConcernState struct {
	Name        ConcernName `json:"name"`        // name of issue
	TrackerURL  string      `json:"tracker_url"` // link to issue on an external issue tracker, such as GitHub
	Closed      bool        `json:"closed"`
	Cancelled   bool        `json:"cancelled"`
	AddressedBy []string    `json:"addressed_by"` // prs addressing this issue
}

// ProposalName is the name of a pull request within the gov4git system.
// ProposalNames must be unique within a community.
type ProposalName string

// ProposalState is the current state of a pull request.
type ProposalState struct {
	Name       ProposalName  `json:"name"`        // name of pr
	TrackerURL string        `json:"tracker_url"` // link to pr on an external pr tracker, such as GitHub
	Closed     bool          `json:"closed"`
	Cancelled  bool          `json:"cancelled"`
	Addresses  []ConcernName `json:"addresses"` // issues addressed by this pr
}
