package github

import (
	"fmt"
	"testing"

	"github.com/google/go-github/v55/github"
	govgh "github.com/gov4git/gov4git/github"
	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/gov4git/runtime"
	"github.com/gov4git/gov4git/test"
	"github.com/gov4git/lib4git/base"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/testutil"
	"github.com/migueleliasweb/go-github-mock/src/mock"
)

var (
	testProcessJoinRequestsOrganizerGithubUser = "organizer"
	testProcessJoinRequestsApplicantGithubUser = "applicant"
	testProcessJoinRequestsGetComments         = []any{
		[]*github.IssueComment{
			{
				User: &github.User{Login: github.String(testProcessJoinRequestsOrganizerGithubUser)},
				Body: github.String("Approve."),
			},
		},
	}
	testProcessJoinRequestsCreateComments = []any{
		&github.IssueComment{
			User: &github.User{Login: github.String(testProcessJoinRequestsOrganizerGithubUser)},
			Body: github.String("Approve."),
		},
	}
	testProcessJoinRequestsEditIssue = []any{
		&github.Issue{},
	}
)

func TestProcessJoinRequests(t *testing.T) {
	base.LogVerbosely()

	// init governance
	ctx := testutil.NewCtx(t, runtime.TestWithCache)
	cty := test.NewTestCommunity(t, ctx, 2)

	// init join applicant's identity
	applicantID := id.NewTestID(ctx, t, git.MainBranch, true)
	id.Init(ctx, applicantID.OwnerAddress())

	testProcessJoinRequestsGetIssues := []any{
		[]*github.Issue{
			{ // issue without governance
				ID:     github.Int64(111),
				Number: github.Int(1),
				Title:  github.String("Issue 1"),
				URL:    github.String("https://test/issue/1"),
				Labels: []*github.Label{{Name: github.String(govgh.JoinRequestLabel)}},
				Locked: github.Bool(false),
				State:  github.String("open"),
				Body: github.String(
					fmt.Sprintf("### Your public repo\n\n%v\n\n### Your public branch\n\n%v\n\n### Your email (optional)\n\n%v",
						applicantID.Public.Dir(), git.MainBranch, "test@test"),
				),
				User:     &github.User{Login: github.String(testProcessJoinRequestsApplicantGithubUser)},
				Comments: github.Int(1),
			},
		},
	}

	// init mock github
	mockedHTTPClient := mock.NewMockedHTTPClient(
		mock.WithRequestMatch(mock.GetReposIssuesByOwnerByRepo, testProcessJoinRequestsGetIssues...),
		mock.WithRequestMatch(mock.GetReposIssuesCommentsByOwnerByRepoByIssueNumber, testProcessJoinRequestsGetComments...),
		mock.WithRequestMatch(mock.PostReposIssuesCommentsByOwnerByRepoByIssueNumber, testProcessJoinRequestsCreateComments...),
		mock.WithRequestMatch(mock.PatchReposIssuesByOwnerByRepoByIssueNumber, testProcessJoinRequestsEditIssue...),
	)
	ghRepo := govgh.GithubRepo{Owner: "owner1", Name: "repo1"}
	ghClient := github.NewClient(mockedHTTPClient)

	// process join requests
	chg := govgh.ProcessJoinRequestIssues(ctx, ghRepo, ghClient, cty.Organizer(), []string{testProcessJoinRequestsOrganizerGithubUser})
	if len(chg.Result) != 1 {
		t.Fatalf("expecting 1 join")
	}
	if chg.Result[0] != testProcessJoinRequestsApplicantGithubUser {
		t.Errorf("expecting %v, got %v", testProcessJoinRequestsApplicantGithubUser, chg.Result[0])
	}

	// <-(chan int(nil))
}
