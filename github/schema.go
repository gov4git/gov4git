package github

import (
	"strconv"
	"time"

	"github.com/gov4git/lib4git/ns"
)

type GithubRepo struct {
	Owner string `json:"github_repo_owner"`
	Name  string `json:"github_repo_name"`
}

const PrioritizeIssueByGovernanceLabel = "gov:prioritize"

type GithubBallotIssue struct {
	URL       string
	Number    int64
	Title     string
	Body      string
	Labels    []string
	ClosedAt  *time.Time
	CreatedAt *time.Time
	UpdatedAt *time.Time
	Locked    bool
	Closed    bool
}

func (x GithubBallotIssue) BallotName() ns.NS {
	return ns.NS{"issue", strconv.Itoa(int(x.Number))}
}

type GithubBallotIssues []GithubBallotIssue
