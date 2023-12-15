package boot

import (
	"testing"

	"github.com/gov4git/gov4git/v2/proto/member"
	"github.com/gov4git/gov4git/v2/runtime"
	"github.com/gov4git/gov4git/v2/test"
	"github.com/gov4git/lib4git/base"
	"github.com/gov4git/lib4git/testutil"
)

func TestBoot(t *testing.T) {
	base.LogVerbosely()
	ctx := testutil.NewCtx(t, runtime.TestWithCache)
	cty := test.NewTestCommunity(t, ctx, 2)

	if !member.IsGroup(ctx, cty.Gov(), member.Everybody) {
		t.Errorf("expecting group %v", member.Everybody)
	}
}
