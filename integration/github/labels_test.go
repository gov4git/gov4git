//go:build integration
// +build integration

package github_test

import (
	"context"
	"testing"

	"github.com/google/go-github/v55/github"
	govgh "github.com/gov4git/gov4git/v2/github"
)

func TestCreateLabel(t *testing.T) {
	ctx := context.Background()
	testLabel := "xyz:test-label"

	client.Issues.DeleteLabel(ctx, TestRepo.Owner, TestRepo.Name, testLabel)

	label := &github.Label{Name: github.String(testLabel)}

	_, _, err := client.Issues.CreateLabel(ctx, TestRepo.Owner, TestRepo.Name, label)
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = client.Issues.CreateLabel(ctx, TestRepo.Owner, TestRepo.Name, label)
	if err == nil {
		t.Fatalf("error is expected")
	}

	if !govgh.IsLabelAlreadyExists(err) {
		t.Errorf("not expecting %v", err)
	}
}
