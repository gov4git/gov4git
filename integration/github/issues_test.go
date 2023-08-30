//go:build integration
// +build integration

package github_test

import (
	"context"
	"sort"
	"testing"

	"github.com/google/go-github/v54/github"
	govgh "github.com/gov4git/gov4git/github"
	"github.com/gov4git/lib4git/util"
)

func TestGithubIssueStructure(t *testing.T) {
	issues, _, err := client.Issues.ListByRepo(context.Background(), TestRepo.Owner, TestRepo.Name, &github.IssueListByRepoOptions{State: "all"})
	if err != nil {
		t.Fatalf("Issues.ListByRepo returned error: %v", err)
	}
	issuesByNumber(issues).Sort()
	// fmt.Println("ISSUES", form.SprintJSON(issues))

	if len(issues) < 4 {
		t.Fatalf("Expected at least 4 issues, got %v", len(issues))
	}

	// test issue 1
	if issues[0].GetNumber() != 1 {
		t.Fatalf("Expected issue number 1, got %v", issues[0].GetNumber())
	}
	if issues[0].GetState() != "open" {
		t.Fatalf("Expected issue state 'open', got %v", issues[0].GetState())
	}
	if !util.IsIn(govgh.PrioritizeIssueByGovernanceLabel, govgh.LabelsToStrings(issues[0].Labels)...) {
		t.Fatalf("Expected issue to be prioritized")
	}

	// test issue 2
	if issues[1].GetNumber() != 2 {
		t.Fatalf("Expected issue number 2, got %v", issues[1].GetNumber())
	}
	if !issues[1].GetLocked() {
		t.Fatalf("Expected issue to be locked")
	}

	// test issue 3
	if issues[2].GetNumber() != 3 {
		t.Fatalf("Expected issue number 3, got %v", issues[2].GetNumber())
	}
	if issues[2].GetState() != "closed" {
		t.Fatalf("Expected issue state 'open', got %v", issues[2].GetState())
	}

}

type issuesByNumber []*github.Issue

func (x issuesByNumber) Sort() {
	sort.Sort(x)
}

func (x issuesByNumber) Len() int {
	return len(x)
}

func (x issuesByNumber) Less(i, j int) bool {
	return x[i].GetNumber() < x[j].GetNumber()
}

func (x issuesByNumber) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}
