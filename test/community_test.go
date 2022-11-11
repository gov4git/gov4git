package test

import (
	"testing"

	"github.com/gov4git/gov4git/lib/testutil"
)

func TestTestCommunity(t *testing.T) {
	ctx := testutil.NewCtx()
	NewTestCommunity(t, ctx, 3)
	testutil.Hang()
}
