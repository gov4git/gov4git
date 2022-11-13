package test

import (
	"testing"

	"github.com/gov4git/lib4git/testutil"
)

func TestTestCommunity(t *testing.T) {
	// base.LogVerbosely()
	ctx := testutil.NewCtx()
	NewTestCommunity(t, ctx, 3)
	// testutil.Hang()
}
