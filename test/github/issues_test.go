package github

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-github/v54/github"
	govgh "github.com/gov4git/gov4git/github"
	"github.com/gov4git/gov4git/proto/ballot/ballot"
	"github.com/gov4git/gov4git/test"
	"github.com/gov4git/lib4git/base"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/testutil"
	"github.com/migueleliasweb/go-github-mock/src/mock"
)

func TestGithubMock(t *testing.T) {

	mockedHTTPClient := mock.NewMockedHTTPClient(
		mock.WithRequestMatch(
			mock.GetReposIssuesByOwnerByRepo,
			[]github.Issue{
				{ // issue without governance
					ID:     github.Int64(123),
					Number: github.Int(1),
					Title:  github.String("Issue 1"),
					URL:    github.String("https://test/issue/1"),
					Labels: []*github.Label{},
					Locked: github.Bool(false),
					//
					State: github.String("open"),
				},
				{ // open, non-frozen issue with governance
					ID:     github.Int64(456),
					Number: github.Int(2),
					Title:  github.String("Issue 2"),
					URL:    github.String("https://test/issue/2"),
					Labels: []*github.Label{{Name: github.String(govgh.PrioritizeIssueByGovernanceLabel)}},
					Locked: github.Bool(false),
					//
					State: github.String("open"),
				},
			},
			[]github.Issue{},
		),
	)

	c := github.NewClient(mockedHTTPClient)

	ctx := context.Background()

	issues1, _, repo1Err := c.Issues.ListByRepo(ctx, "owner1", "repo1", &github.IssueListByRepoOptions{})
	if len(issues1) != 2 || repo1Err != nil {
		t.Errorf("unexpected")
	}

	issues2, _, repo2Err := c.Issues.ListByRepo(ctx, "owner1", "repo2", &github.IssueListByRepoOptions{})
	if len(issues2) != 0 || repo2Err != nil {
		t.Errorf("unexpected")
	}
}

var (
	testImportIssuesForPrioritization = []interface{}{
		// request #1
		[]github.Issue{
			{ // issue without governance
				ID:     github.Int64(111),
				Number: github.Int(1),
				Title:  github.String("Issue 1"),
				URL:    github.String("https://test/issue/1"),
				Labels: []*github.Label{},
				Locked: github.Bool(false),
				State:  github.String("open"),
			},
			{ // issue with governance, open, not-frozen
				ID:     github.Int64(222),
				Number: github.Int(2),
				Title:  github.String("Issue 2"),
				URL:    github.String("https://test/issue/2"),
				Labels: []*github.Label{{Name: github.String(govgh.PrioritizeIssueByGovernanceLabel)}},
				Locked: github.Bool(false),
				State:  github.String("open"),
			},
			{ // issue with governance, open, frozen
				ID:     github.Int64(333),
				Number: github.Int(3),
				Title:  github.String("Issue 3"),
				URL:    github.String("https://test/issue/3"),
				Labels: []*github.Label{{Name: github.String(govgh.PrioritizeIssueByGovernanceLabel)}},
				Locked: github.Bool(true),
				State:  github.String("open"),
			},
			{ // issue with governance, closed
				ID:     github.Int64(444),
				Number: github.Int(4),
				Title:  github.String("Issue 4"),
				URL:    github.String("https://test/issue/4"),
				Labels: []*github.Label{{Name: github.String(govgh.PrioritizeIssueByGovernanceLabel)}},
				Locked: github.Bool(true),
				State:  github.String("closed"),
			},
		},
		// request #2
		[]github.Issue{
			{ // issue without governance -> with governance, open, not-frozen
				ID:     github.Int64(111),
				Number: github.Int(1),
				Title:  github.String("Issue 1"),
				URL:    github.String("https://test/issue/1"),
				Labels: []*github.Label{{Name: github.String(govgh.PrioritizeIssueByGovernanceLabel)}},
				Locked: github.Bool(false),
				State:  github.String("open"),
			},
			{ // issue with governance, open, not-frozen -> without governance, open, frozen (XXX: this hits a bug: during filtering it is removed and not considered)
				ID:     github.Int64(222),
				Number: github.Int(2),
				Title:  github.String("Issue 2"),
				URL:    github.String("https://test/issue/2"),
				Labels: []*github.Label{},
				Locked: github.Bool(false),
				State:  github.String("open"),
			},
			{ // issue with governance, open, frozen -> with governance, closed
				ID:     github.Int64(333),
				Number: github.Int(3),
				Title:  github.String("Issue 3"),
				URL:    github.String("https://test/issue/3"),
				Labels: []*github.Label{{Name: github.String(govgh.PrioritizeIssueByGovernanceLabel)}},
				Locked: github.Bool(true),
				State:  github.String("open"),
			},
			{ // issue with governance, closed -> with govrnance reopen
				ID:     github.Int64(444),
				Number: github.Int(4),
				Title:  github.String("Issue 4"),
				URL:    github.String("https://test/issue/4"),
				Labels: []*github.Label{{Name: github.String(govgh.PrioritizeIssueByGovernanceLabel)}},
				Locked: github.Bool(true),
				State:  github.String("closed"),
			},
		},
	}
)

func TestImportIssuesForPrioritization(t *testing.T) {
	base.LogVerbosely()

	// init mock github
	mockedHTTPClient := mock.NewMockedHTTPClient(
		mock.WithRequestMatch(mock.GetReposIssuesByOwnerByRepo, testImportIssuesForPrioritization...),
	)
	ghRepo := govgh.GithubRepo{Owner: "owner1", Name: "repo1"}
	ghClient := github.NewClient(mockedHTTPClient)

	// init governance
	ctx := testutil.NewCtx()
	cty := test.NewTestCommunity(t, ctx, 2)

	// import issues #1
	chg1 := govgh.ImportIssuesForPrioritization(ctx, ghRepo, ghClient, cty.Organizer())
	fmt.Println("IMPORT#1", form.SprintJSON(chg1.Result))

	// list #1
	ads1 := ballot.List(ctx, cty.Gov())
	fmt.Println("ADS#1", form.SprintJSON(ads1))
	if len(ads1) != 3 {
		t.Errorf("expecting 3, got %v", len(ads1))
	}

	// import issues #2
	chg2 := govgh.ImportIssuesForPrioritization(ctx, ghRepo, ghClient, cty.Organizer())
	fmt.Println("IMPORT#2", form.SprintJSON(chg2.Result))

	// list #2
	ads2 := ballot.List(ctx, cty.Gov())
	fmt.Println("ADS#2", form.SprintJSON(ads2))
	if len(ads2) != 4 {
		t.Errorf("expecting 4, got %v", len(ads2))
	}

	// testutil.Hang()
}
