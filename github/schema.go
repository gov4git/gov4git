package github

import (
	"context"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/docket/schema"
	"github.com/gov4git/lib4git/must"
)

type Repo struct {
	Owner string `json:"github_repo_owner"`
	Name  string `json:"github_repo_name"`
}

func (x Repo) HTTPS() string {
	return `https://github.com/` + x.Owner + `/` + x.Name
}

func ParseRepo(ctx context.Context, s string) Repo {
	first, second, ok := strings.Cut(s, "/")
	must.Assertf(ctx, ok, "not a github repo: %v", s)
	return Repo{Owner: first, Name: second}
}

const (
	// prioritizing issues by ballot
	PrioritizeIssueByGovernanceLabel = "gov4git:prioritize"
	PrioritizeBallotChoice           = "prioritize"

	// member join
	JoinRequestLabel        = "gov4git:join"
	JoinRequestApprovalWord = "approve"

	// organizer directives
	DirectiveLabel = "gov4git:directive"

	// labels for issues that are managed
	IssueIsManagedLabel = "gov4git:managed"

	// Github deploy environment
	DeployEnvName = "gov4git:governance"
)

var GovernanceLabels = []string{
	PrioritizeIssueByGovernanceLabel,
	JoinRequestLabel,
	DirectiveLabel,
	IssueIsManagedLabel,
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
	IsPullRequest     bool `json:"is_pull_request"`
	IsManaged         bool `json:"is_governed"`
	ForPrioritization bool `json:"for_prioritization"`
}

func (x ImportedIssue) Key() string {
	return strconv.Itoa(int(x.Number))
}

func IssueNumberToMotionID(no int64) schema.MotionID {
	return schema.MotionID(strconv.Itoa(int(no)))
}

const (
	ImportedGithubPrefix = "github"
	ImportedIssuePrefix  = "issues"
	ImportedPullPrefix   = "pull"
)

func (x ImportedIssue) BallotName() common.BallotName {
	if x.IsPullRequest {
		return common.BallotName{ImportedGithubPrefix, ImportedPullPrefix, x.Key()}
	} else {
		return common.BallotName{ImportedGithubPrefix, ImportedIssuePrefix, x.Key()}
	}
}

func (x ImportedIssue) MotionType() schema.MotionType {
	if x.IsPullRequest {
		return schema.MotionProposalType
	} else {
		return schema.MotionConcernType
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
