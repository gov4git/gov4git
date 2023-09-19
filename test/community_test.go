package test

import (
	"testing"

	"github.com/gov4git/gov4git/runtime"
	"github.com/gov4git/lib4git/testutil"
)

func TestTestCommunity(t *testing.T) {
	// base.LogVerbosely()
	ctx := testutil.NewCtx(t, runtime.TestWithCache)
	NewTestCommunity(t, ctx, 3)
	// testutil.Hang()
}
