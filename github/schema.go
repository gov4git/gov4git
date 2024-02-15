package github

import (
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/motion/motionpolicies/pmp_0"
	"github.com/gov4git/gov4git/v2/proto/motion/motionpolicies/pmp_1"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
	"github.com/gov4git/gov4git/v2/proto/purpose"
)

const (
	// prioritizing issues by ballot
	PrioritizeIssueByGovernanceLabel = "gov4git:prioritize"
	PrioritizeBallotChoice           = "prioritize"

	// member join
	JoinRequestApprovalWord = "approve"

	// organizer directives
	DirectiveLabel = "gov4git:directive"

	// labels for issues that are managed
	// TODO: generate labels from the descriptors of installed motion policies
	IssueIsManagedLabel        = "gov4git:managed"
	IssueIsManagedByPMPv0Label = pmp_0.ConcernPolicyGithubLabel
	IssueIsManagedByPMPv1Label = pmp_1.ConcernPolicyGithubLabel

	// the issue with this label will be used as a dashboard display
	DashboardIssueLabel = "gov4git:dashboard"

	// Github deploy environment
	DeployEnvName = "gov4git:governance"
)

var GovernanceLabels = []string{
	PrioritizeIssueByGovernanceLabel,
	DirectiveLabel,
	IssueIsManagedLabel,
	IssueIsManagedByPMPv0Label,
	IssueIsManagedByPMPv1Label,
	DashboardIssueLabel,
}

type ImportedIssue struct {
	Author string `json:"author"`
	Number int64  `json:"number"`
	// meta
	URL    string   `json:"url"`
	Title  string   `json:"title"`
	Body   string   `json:"body"`
	Labels []string `json:"labels"`
	//
	ClosedAt  *time.Time `json:"closed_at"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	//
	Refs []ImportedRef `json:"refs"`
	//
	Locked            bool `json:"locked"`
	Closed            bool `json:"closed"`
	PullRequest       bool `json:"pull_request"`
	Merged            bool `json:"merged"`
	ManagedByPMPv0    bool `json:"managed_by_pmp_v0,omitempty"`
	ManagedByPMPv1    bool `json:"managed_by_pmp_v1,omitempty"`
	ForPrioritization bool `json:"for_prioritization"`
}

func (x ImportedIssue) IsManaged() bool {
	return x.ManagedByPMPv0 || x.ManagedByPMPv1
}

func (x ImportedIssue) Key() string {
	return strconv.Itoa(int(x.Number))
}

func (x ImportedIssue) MotionID() motionproto.MotionID {
	return IssueNumberToMotionID(x.Number)
}

func (x ImportedIssue) Purpose() purpose.Purpose {
	if x.PullRequest {
		return purpose.Proposal
	}
	return purpose.Concern
}

func MotionIDToIssueNumber(id motionproto.MotionID) (int, error) {
	return strconv.Atoi(id.String())
}

func IssueNumberToMotionID(no int64) motionproto.MotionID {
	return motionproto.MotionID(strconv.Itoa(int(no)))
}

const (
	ImportedGithubPrefix = "github"
	ImportedIssuePrefix  = "issues"
	ImportedPullPrefix   = "pull"
)

func (x ImportedIssue) BallotName() ballotproto.BallotID {
	if x.PullRequest {
		return ballotproto.ParseBallotID(strings.Join([]string{
			ImportedGithubPrefix,
			ImportedPullPrefix,
			x.Key(),
		}, "/"))
	} else {
		return ballotproto.ParseBallotID(strings.Join([]string{
			ImportedGithubPrefix,
			ImportedIssuePrefix,
			x.Key(),
		}, "/"))
	}
}

func (x ImportedIssue) MotionType() motionproto.MotionType {
	if x.PullRequest {
		return motionproto.MotionProposalType
	} else {
		return motionproto.MotionConcernType
	}
}

type ImportedRef struct {
	To   int64  `json:"to"`
	Type string `json:"type"`
}

type ImportedIssues []ImportedIssue

func (x ImportedIssues) Sort() {
	sort.Sort(x)
}

func (x ImportedIssues) Len() int {
	return len(x)
}

func (x ImportedIssues) Less(i, j int) bool {
	return x[i].Number < x[j].Number
}

func (x ImportedIssues) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}
