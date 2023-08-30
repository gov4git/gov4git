package github

import (
	"context"
	"testing"

	"github.com/google/go-github/v54/github"
	govgh "github.com/gov4git/gov4git/github"
	"github.com/gov4git/gov4git/proto/ballot/ballot"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/test"
	"github.com/gov4git/lib4git/base"
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
		[]github.Issue{
			{ // issue without governance -> with governance, open, not-frozen
				ID:     github.Int64(111),
				Number: github.Int(1),
				Title:  github.String("Issue 1"),
				URL:    github.String("https://test/issue/1"),
				Labels: []*github.Label{},
				Locked: github.Bool(false),
				State:  github.String("open"),
			},
			{ // issue with governance, open, not-frozen -> without governance (XXX: this hits a bug: during filtering it is removed and not considered)
				ID:     github.Int64(222),
				Number: github.Int(2),
				Title:  github.String("Issue 2"),
				URL:    github.String("https://test/issue/2"),
				Labels: []*github.Label{{Name: github.String(govgh.PrioritizeIssueByGovernanceLabel)}},
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

	// XXX: import issues #1
	govgh.ImportIssuesForPrioritization(XXX)

	// list
	ads := ballot.List(ctx, cty.Gov())
	if len(ads) != 1 {
		t.Errorf("expecting 1 ad, got %v", len(ads))
	}

	// XXX: import issues #2

	// verify ballots
	ast := ballot.Show(ctx, gov.GovAddress(cty.Organizer().Public), ballotName)
	if !ast.Ad.Closed {
		t.Errorf("expecting closed flag")
	}

	// testutil.Hang()
}
