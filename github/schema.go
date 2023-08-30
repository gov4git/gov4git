package github

import (
	"sort"
	"strconv"
	"time"

	"github.com/gov4git/lib4git/ns"
)

type GithubRepo struct {
	Owner string `json:"github_repo_owner"`
	Name  string `json:"github_repo_name"`
}

const (
	PrioritizeIssueByGovernanceLabel = "gov:prioritize"
	PrioritizeBallotChoice           = "prioritize"
)

type GithubIssueBallot struct {
	ForPrioritization bool
	URL               string
	Number            int64
	Title             string
	Body              string
	Labels            []string
	ClosedAt          *time.Time
	CreatedAt         *time.Time
	UpdatedAt         *time.Time
	Locked            bool
	Closed            bool
}

func (x GithubIssueBallot) Key() string {
	return strconv.Itoa(int(x.Number))
}

func (x GithubIssueBallot) BallotName() ns.NS {
	return ns.NS{"issue", x.Key()}
}

type GithubIssueBallots []GithubIssueBallot

func (x GithubIssueBallots) Sort() {
	sort.Sort(x)
}

func (x GithubIssueBallots) Len() int {
	return len(x)
}

func (x GithubIssueBallots) Less(i, j int) bool {
	return x[i].Number < x[j].Number
}

func (x GithubIssueBallots) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}
