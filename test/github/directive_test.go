package github

import (
	"fmt"
	"testing"

	"github.com/google/go-github/v55/github"
	govgh "github.com/gov4git/gov4git/github"
	"github.com/gov4git/gov4git/proto/balance"
	"github.com/gov4git/gov4git/proto/ballot/qv"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/runtime"
	"github.com/gov4git/gov4git/test"
	"github.com/gov4git/lib4git/base"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/testutil"
	"github.com/migueleliasweb/go-github-mock/src/mock"
)

var (
	testDirectiveOrganizerGithubUser = "organizer"
	testDirectiveEditIssue           = []any{
		&github.Issue{},
		&github.Issue{},
	}
	testDirectivePostComments = []any{
		&github.IssueComment{},
		&github.IssueComment{},
	}
)

func TestDirective(t *testing.T) {
	base.LogVerbosely()

	// init governance
	ctx := testutil.NewCtx(t, runtime.TestWithCache)
	cty := test.NewTestCommunity(t, ctx, 2)

	testIssueAmount := 20.0
	testTransferAmount := 10.0
	testDirectiveGetIssues := []any{
		[]*github.Issue{
			{
				ID:     github.Int64(111),
				Number: github.Int(1),
				Title:  github.String("Issue directive"),
				URL:    github.String("https://test/issue/1"),
				Labels: []*github.Label{{Name: github.String(govgh.DirectiveLabel)}},
				Locked: github.Bool(false),
				State:  github.String("open"),
				Body: github.String(
					fmt.Sprintf("issue %v voting credits to @%v", testIssueAmount, cty.MemberUser(0)),
				),
				User: &github.User{Login: github.String(testDirectiveOrganizerGithubUser)},
			},
			{
				ID:     github.Int64(222),
				Number: github.Int(2),
				Title:  github.String("Transfer directive"),
				URL:    github.String("https://test/issue/2"),
				Labels: []*github.Label{{Name: github.String(govgh.DirectiveLabel)}},
				Locked: github.Bool(false),
				State:  github.String("open"),
				Body: github.String(
					fmt.Sprintf("transfer %v voting credits from @%v to @%v", testTransferAmount, cty.MemberUser(0), cty.MemberUser(1)),
				),
				User: &github.User{Login: github.String(testDirectiveOrganizerGithubUser)},
			},
		},
	}

	// init mock github
	mockedHTTPClient := mock.NewMockedHTTPClient(
		// fetch issues
		mock.WithRequestMatch(mock.GetReposIssuesByOwnerByRepo,
			testDirectiveGetIssues...),
		// issue + transfer directives execution
		mock.WithRequestMatch(mock.PostReposIssuesCommentsByOwnerByRepoByIssueNumber,
			testDirectivePostComments...),
		mock.WithRequestMatch(mock.PatchReposIssuesByOwnerByRepoByIssueNumber,
			testDirectiveEditIssue...),
	)
	ghRepo := govgh.GithubRepo{Owner: "owner1", Name: "repo1"}
	ghClient := github.NewClient(mockedHTTPClient)

	// process directives
	chg := govgh.ProcessDirectiveIssues(ctx, ghRepo, ghClient, cty.Organizer(), []string{testDirectiveOrganizerGithubUser})
	if len(chg.Result) != 2 {
		t.Fatalf("expecting 2 directives")
	}
	fmt.Println(form.SprintJSON(chg.Result))

	b1 := balance.Get(ctx, gov.GovAddress(cty.Organizer().Public), cty.MemberUser(0), qv.VotingCredits)
	if b1 != 10.0 {
		t.Errorf("expecting %v, got %v", 10.0, b1)
	}
	b2 := balance.Get(ctx, gov.GovAddress(cty.Organizer().Public), cty.MemberUser(1), qv.VotingCredits)
	if b2 != 10.0 {
		t.Errorf("expecting %v, got %v", 10.0, b2)
	}

	// <-(chan int(nil))
}
