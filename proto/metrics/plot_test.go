package metrics

import (
	"context"
	"math/rand"
	"os"
	"testing"
	"time"
)

func TestSeriesPlot(t *testing.T) {
	t.SkipNow()

	ctx := context.Background()
	os.Setenv("PATH", "/opt/homebrew/bin")

	pngMotions := plotDailyMotionsPNG(ctx, testSeries)
	os.WriteFile("plotDailyMotionsPNG.png", pngMotions, 0644)

	pngCredits := plotDailyCreditsPNG(ctx, testSeries)
	os.WriteFile("plotDailyCreditsPNG.png", pngCredits, 0644)

	pngCleared := plotDailyClearedPNG(ctx, testSeries)
	os.WriteFile("plotDailyClearedPNG.png", pngCleared, 0644)

	pngVotes := plotDailyVotesPNG(ctx, testSeries)
	os.WriteFile("plotDailyVotesPNG.png", pngVotes, 0644)

	pngCharges := plotDailyChargesPNG(ctx, testSeries)
	os.WriteFile("plotDailyChargesPNG.png", pngCharges, 0644)

	pngJoins := plotDailyJoinsPNG(ctx, testSeries)
	os.WriteFile("plotDailyJoinsPNG.png", pngJoins, 0644)

}

func testRandFloatArray(n int) []float64 {
	r := make([]float64, n)
	for i := range r {
		r[i] = rand.NormFloat64() + 3
	}
	return r
}

var (
	testDates = []time.Time{
		time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 11, 2, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 11, 3, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 11, 4, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 11, 5, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 11, 6, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 11, 7, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 11, 8, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 11, 9, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 11, 10, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 11, 11, 0, 0, 0, 0, time.UTC),
	}
	testSeries = &Series{
		DailyNumJoins: DailySeries{
			X: testDates,
			Y: testRandFloatArray(len(testDates)),
		},
		DailyNumMotionOpen: DailySeries{
			X: testDates,
			Y: testRandFloatArray(len(testDates)),
		},
		DailyNumMotionClose: DailySeries{
			X: testDates,
			Y: testRandFloatArray(len(testDates)),
		},
		DailyNumMotionCancel: DailySeries{
			X: testDates,
			Y: testRandFloatArray(len(testDates)),
		},
		DailyCreditsIssued: DailySeries{
			X: testDates,
			Y: testRandFloatArray(len(testDates)),
		},
		DailyCreditsBurned: DailySeries{
			X: testDates,
			Y: testRandFloatArray(len(testDates)),
		},
		DailyCreditsTransferred: DailySeries{
			X: testDates,
			Y: testRandFloatArray(len(testDates)),
		},
		DailyClearedBounties: DailySeries{
			X: testDates,
			Y: testRandFloatArray(len(testDates)),
		},
		DailyClearedRewards: DailySeries{
			X: testDates,
			Y: testRandFloatArray(len(testDates)),
		},
		DailyClearedRefunds: DailySeries{
			X: testDates,
			Y: testRandFloatArray(len(testDates)),
		},
		DailyNumConcernVotes: DailySeries{
			X: testDates,
			Y: testRandFloatArray(len(testDates)),
		},
		DailyNumProposalVotes: DailySeries{
			X: testDates,
			Y: testRandFloatArray(len(testDates)),
		},
		DailyNumOtherVotes: DailySeries{
			X: testDates,
			Y: testRandFloatArray(len(testDates)),
		},
		DailyConcernVoteCharges: DailySeries{
			X: testDates,
			Y: testRandFloatArray(len(testDates)),
		},
		DailyProposalVoteCharges: DailySeries{
			X: testDates,
			Y: testRandFloatArray(len(testDates)),
		},
		DailyOtherVoteCharges: DailySeries{
			X: testDates,
			Y: testRandFloatArray(len(testDates)),
		},
	}
)
