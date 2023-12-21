package metrics

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/history"
	"github.com/gov4git/gov4git/v2/proto/journal"
)

func AssembleReport_Local(
	ctx context.Context,
	cloned gov.Cloned,
	urlCalc AssetURLCalculator,
	earliest time.Time,
	latest time.Time,

) *ReportAssets {

	entries := loadHistory_Local(ctx, cloned, earliest, latest)
	series := ComputeSeries(entries, earliest, latest)

	var w bytes.Buffer

	fmt.Fprintf(&w, "## Past 30-days metrics\n\n")

	fmt.Fprintf(&w, "### Aggregates\n\n")

	fmt.Fprintf(&w, "| Indicator|  30-day aggregate |\n")
	fmt.Fprintf(&w, "|  ---:|  :--- |\n")
	fmt.Fprintf(&w, "| Number of opened motions | %d |\n", int(series.DailyNumMotionOpen.Total()))
	fmt.Fprintf(&w, "| Number of closed motions | %d |\n", int(series.DailyNumMotionClose.Total()))
	fmt.Fprintf(&w, "| Number of cancelled motions | %d |\n", int(series.DailyNumMotionCancel.Total()))
	fmt.Fprintf(&w, "|  ---:|  :--- |\n")
	fmt.Fprintf(&w, "| Number of votes made | %d |\n", int(series.DailyNumVotes.Total()))
	fmt.Fprintf(&w, "| Credits spent on votes | %0.6f |\n", series.DailyVoteCharges.Total())
	fmt.Fprintf(&w, "|  ---:|  :--- |\n")
	fmt.Fprintf(&w, "| Credits cleared in bounties | %0.6f |\n", series.DailyClearedBounties.Total())
	fmt.Fprintf(&w, "| Credits cleared in rewards | %0.6f |\n", series.DailyClearedRewards.Total())
	fmt.Fprintf(&w, "| Credits cleared in refunds | %0.6f |\n", series.DailyClearedRefunds.Total())
	fmt.Fprintf(&w, "|  ---:|  :--- |\n")
	fmt.Fprintf(&w, "| Number of new members | %d |\n", int(series.DailyNumJoins.Total()))
	fmt.Fprintf(&w, "|  ---:|  :--- |\n")
	fmt.Fprintf(&w, "| Credits issued | %0.6f |\n", series.DailyCreditsIssued.Total())
	fmt.Fprintf(&w, "| Credits burned | %0.6f |\n", series.DailyCreditsBurned.Total())
	fmt.Fprintf(&w, "| Credits transferred | %0.6f |\n", series.DailyCreditsTransferred.Total())

	fmt.Println("### Daily breakdown")

	fmt.Fprintf(&w, "![%s](%s)\n<hr>\n", "Daily motions opened/closed/cancelled", urlCalc("daily_motions.png"))
	fmt.Fprintf(&w, "![%s](%s)\n<hr>\n", "Daily vote counts", urlCalc("daily_votes.png"))
	fmt.Fprintf(&w, "![%s](%s)\n<hr>\n", "Daily vote charges", urlCalc("daily_charges.png"))
	fmt.Fprintf(&w, "![%s](%s)\n<hr>\n", "Daily credits cleared in bounties/rewards/refunds", urlCalc("daily_cleared.png"))
	fmt.Fprintf(&w, "![%s](%s)\n<hr>\n", "Daily new community members", urlCalc("daily_joins.png"))
	fmt.Fprintf(&w, "![%s](%s)\n<hr>\n", "Daily credits issued/burned/transferred", urlCalc("daily_credits.png"))

	return &ReportAssets{
		ReportMD: w.String(),
		Assets: map[string][]byte{
			"daily_motions.png": plotDailyMotionsPNG(ctx, series),
			"daily_votes.png":   plotDailyVotesPNG(ctx, series),
			"daily_charges.png": plotDailyChargesPNG(ctx, series),
			"daily_cleared.png": plotDailyClearedPNG(ctx, series),
			"daily_joins.png":   plotDailyJoinsPNG(ctx, series),
			"daily_credits.png": plotDailyCreditsPNG(ctx, series),
		},
	}
}

func loadHistory_Local(
	ctx context.Context,
	cloned gov.Cloned,
	earliest time.Time,
	latest time.Time,

) journal.Entries[*history.Event] {

	all := history.List_Local(ctx, cloned)
	s := journal.Entries[*history.Event]{}
	for _, entry := range all {
		if isNoEarlierThan(entry.Stamp, earliest) && isNoLaterThan(entry.Stamp, latest) {
			s = append(s, entry)
		}
	}
	return s
}
