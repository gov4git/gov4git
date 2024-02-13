package metrics

import (
	"fmt"
	"os"
	"testing"

	"github.com/gov4git/gov4git/v2/proto/metrics"
	"github.com/gov4git/gov4git/v2/runtime"
	"github.com/gov4git/gov4git/v2/test"
	pmp "github.com/gov4git/gov4git/v2/test/motion/pmp_0"
	"github.com/gov4git/lib4git/base"
	"github.com/gov4git/lib4git/testutil"
)

func TestDashboardPMP(t *testing.T) {
	t.SkipNow()
	os.Setenv("PATH", "/opt/homebrew/bin")

	base.LogVerbosely()
	ctx := testutil.NewCtx(t, runtime.TestWithCache)
	cty := test.NewTestCommunity(t, ctx, 2)

	pmp.SetupTest(t, ctx, cty)

	urlCalc := func(assetRepoPath string) (url string) {
		return assetRepoPath
	}
	report := metrics.AssembleReport(ctx, cty.Gov(), urlCalc, metrics.TimeDailyLowerBound, metrics.Today().AddDate(0, 0, 1))
	fmt.Println(report.ReportMD)

	if report.Series.AllTime.DailyNumConcernVotes.Total() != 2 {
		t.Errorf("expecting 2, got %v", report.Series.AllTime.DailyNumConcernVotes.Total())
	}
}
