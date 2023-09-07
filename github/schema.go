package github

import (
	"context"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/lib4git/must"
)

type GithubRepo struct {
	Owner string `json:"github_repo_owner"`
	Name  string `json:"github_repo_name"`
}

func ParseGithubRepo(ctx context.Context, s string) GithubRepo {
	first, second, ok := strings.Cut(s, "/")
	must.Assertf(ctx, ok, "not a github repo: %v", s)
	return GithubRepo{Owner: first, Name: second}
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

	// Github deploy environment
	GithubDeployEnvName = "gov4git:governance"
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
	IsPullRequest     bool
}

func (x GithubIssueBallot) Key() string {
	return strconv.Itoa(int(x.Number))
}

const (
	ImportedGithubPrefix = "github"
	ImportedIssuePrefix  = "issues"
	ImportedPullPrefix   = "pull"
)

func (x GithubIssueBallot) BallotName() common.BallotName {
	if x.IsPullRequest {
		return common.BallotName{ImportedGithubPrefix, ImportedPullPrefix, x.Key()}
	} else {
		return common.BallotName{ImportedGithubPrefix, ImportedIssuePrefix, x.Key()}
	}
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
