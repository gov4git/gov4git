//go:build integration
// +build integration

package github_test

import (
	"context"
	"testing"
)

func TestGithubIssueByUserWithCamelCaseLogin(t *testing.T) {
	issue, _, err := client.Issues.Get(context.Background(), TestRepo.Owner, TestRepo.Name, 15)
	if err != nil {
		t.Fatal(err)
	}
	if issue.GetUser() == nil {
		t.Fatalf("no user")
	}
	login := issue.GetUser().GetLogin()
	exp := "Gov4GitTestUser"
	if login != exp {
		t.Errorf("expecting %v, got %v", exp, login)
	}
}
