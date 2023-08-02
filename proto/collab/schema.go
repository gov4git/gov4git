package collab

import (
	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/kv"
)

var (
	collabNS = proto.RootNS.Sub("collab")

	issuesNS = collabNS.Sub("issue")
	issuesKV = kv.KV[IssueName, IssueState]{}

	prNS = collabNS.Sub("pr")
	prKV = kv.KV[IssueName, PRState]{}
)

// QUESTIONS:
//	What if PRs only address an issue partially?

// IssueName is the name of an issue within the gov4git system.
// IssueNames must be unique within a community.
type IssueName string

// IssueState is the current state of an issue.
type IssueState struct {
	Name        IssueName `json:"name"`        // name of issue
	TrackerURL  string    `json:"tracker_url"` // link to issue on an external issue tracker, such as GitHub
	Closed      bool      `json:"closed"`
	AddressedBy []string  `json:"addressed_by"` // prs addressing this issue
}

// PRName is the name of a pull request within the gov4git system.
// PRNames must be unique within a community.
type PRName string

// PRState is the current state of a pull request.
type PRState struct {
	Name       PRName      `json:"name"`        // name of pr
	TrackerURL string      `json:"tracker_url"` // link to pr on an external pr tracker, such as GitHub
	Closed     bool        `json:"closed"`
	Addresses  []IssueName `json:"addresses"` // issues addressed by this pr
}
