package git

import (
	"context"
	"testing"
)

func TestCloneMissingBranch(t *testing.T) {
	local := Local{Path: "tmp"}
	if err := local.Clone(context.Background(), "git@github.com:petar/soul.pub.git", "main"); err != nil {
		t.Fatal(err)
	}
}

func TestCloneLocal(t *testing.T) {
	XXX
}
