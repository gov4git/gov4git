package boot

import (
	"testing"

	"github.com/gov4git/gov4git/proto/member"
	"github.com/gov4git/gov4git/test"
	"github.com/gov4git/lib4git/base"
	"github.com/gov4git/lib4git/testutil"
)

func TestBoot(t *testing.T) {
	base.LogVerbosely()
	ctx := testutil.NewCtx(t, false)
	cty := test.NewTestCommunity(t, ctx, 2)

	if !member.IsGroup(ctx, cty.Gov(), member.Everybody) {
		t.Errorf("expecting group %v", member.Everybody)
	}
}
