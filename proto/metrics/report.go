package metrics

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/gov4git/gov4git/v2/proto/account"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/history/metric"
	"github.com/gov4git/gov4git/v2/proto/journal"
	"github.com/gov4git/gov4git/v2/proto/motion/motionpolicies/pmp_0"
)

func AssembleReport(
	ctx context.Context,
	addr gov.Address,
	urlCalc AssetURLCalculator,
	earliest time.Time,
	latest time.Time,

) *ReportAssets {

	cloned := gov.Clone(ctx, addr)
	return AssembleReport_Local(ctx, cloned, urlCalc, earliest, latest)
}

type ReportSeries struct {
	Last30Days *Series
	AllTime    *Series
}

func AssembleReport_Local(
	ctx context.Context,
	cloned gov.Cloned,
	urlCalc AssetURLCalculator,
	earliest time.Time,
	latest time.Time,

) *ReportAssets {

	entries := loadHistory_Local(ctx, cloned, earliest, latest)
	last30DaysSeries := ComputeSeries(entries, latest.AddDate(0, -1, 0), latest)
	allTimeSeries := ComputeSeries(entries, earliest, latest)

	matchAccount := account.Get_Local(ctx, cloned, pmp_0.MatchingPoolAccountID)
	matchFunds := matchAccount.Balance(account.PluralAsset).Quantity

	capTable := GetCapTable_Local(ctx, cloned)

	var w bytes.Buffer

	fmt.Fprintf(&w, "## Community economics\n\n")

	fmt.Fprintf(&w, "Currently, there are `%0.6f` credits in the __matching fund__.\n\n", matchFunds)

	fmt.Fprintf(&w, "## Last 30-days\n\n")

	fmt.Fprintf(&w, "### Aggregates\n\n")

	fmt.Fprintf(&w, "| Indicator|  30-day aggregate |\n")
	fmt.Fprintf(&w, "|  ---:|  :--- |\n")
	fmt.Fprintf(&w, "| Number of opened issues/PRs | %d |\n", int(last30DaysSeries.DailyNumMotionOpen.Total()))
	fmt.Fprintf(&w, "| Number of closed issues/PRs | %d |\n", int(last30DaysSeries.DailyNumMotionClose.Total()))
	fmt.Fprintf(&w, "| Number of cancelled issues/PRs | %d |\n\n", int(last30DaysSeries.DailyNumMotionCancel.Total()))

	fmt.Fprintf(&w, "| Indicator|  30-day aggregate |\n")
	fmt.Fprintf(&w, "|  ---:|  :--- |\n")
	fmt.Fprintf(&w, "| Number of votes on issues | %d |\n", int(last30DaysSeries.DailyNumConcernVotes.Total()))
	fmt.Fprintf(&w, "| Number of votes on PRs | %d |\n", int(last30DaysSeries.DailyNumProposalVotes.Total()))
	fmt.Fprintf(&w, "| Number of votes on other | %d |\n", int(last30DaysSeries.DailyNumOtherVotes.Total()))
	fmt.Fprintf(&w, "| Credits spent on issue votes | %0.6f |\n", last30DaysSeries.DailyConcernVoteCharges.Total())
	fmt.Fprintf(&w, "| Credits spent on PR votes | %0.6f |\n", last30DaysSeries.DailyProposalVoteCharges.Total())
	fmt.Fprintf(&w, "| Credits spent on other votes | %0.6f |\n\n", last30DaysSeries.DailyOtherVoteCharges.Total())

	fmt.Fprintf(&w, "| Indicator|  30-day aggregate |\n")
	fmt.Fprintf(&w, "|  ---:|  :--- |\n")
	fmt.Fprintf(&w, "| Credits cleared in bounties | %0.6f |\n", last30DaysSeries.DailyClearedBounties.Total())
	fmt.Fprintf(&w, "| Credits cleared in rewards | %0.6f |\n", last30DaysSeries.DailyClearedRewards.Total())
	fmt.Fprintf(&w, "| Credits cleared in refunds | %0.6f |\n\n", last30DaysSeries.DailyClearedRefunds.Total())

	fmt.Fprintf(&w, "| Indicator|  30-day aggregate |\n")
	fmt.Fprintf(&w, "|  ---:|  :--- |\n")
	fmt.Fprintf(&w, "| Number of new members | %d |\n\n", int(last30DaysSeries.DailyNumJoins.Total()))

	fmt.Fprintf(&w, "| Indicator|  30-day aggregate |\n")
	fmt.Fprintf(&w, "|  ---:|  :--- |\n")
	fmt.Fprintf(&w, "| Credits issued | %0.6f |\n", last30DaysSeries.DailyCreditsIssued.Total())
	fmt.Fprintf(&w, "| Credits burned | %0.6f |\n", last30DaysSeries.DailyCreditsBurned.Total())
	fmt.Fprintf(&w, "| Credits transferred | %0.6f |\n\n", last30DaysSeries.DailyCreditsTransferred.Total())

	fmt.Printf("### Daily breakdown\n\n")

	fmt.Fprintf(&w, "<img alt=%q src=%s width=650 />\n\n", "Daily issues/PRs opened/closed/cancelled", urlCalc("daily_motions.png"))
	fmt.Fprintf(&w, "<img alt=%q src=%s width=650 />\n\n", "Daily vote counts", urlCalc("daily_votes.png"))
	fmt.Fprintf(&w, "<img alt=%q src=%s width=650 />\n\n", "Daily vote charges", urlCalc("daily_charges.png"))
	fmt.Fprintf(&w, "<img alt=%q src=%s width=650 />\n\n", "Daily credits cleared in bounties/rewards/refunds", urlCalc("daily_cleared.png"))
	fmt.Fprintf(&w, "<img alt=%q src=%s width=650 />\n\n", "Daily new community members", urlCalc("daily_joins.png"))
	fmt.Fprintf(&w, "<img alt=%q src=%s width=650 />\n\n", "Daily credits issued/burned/transferred", urlCalc("daily_credits.png"))

	fmt.Fprintf(&w, "## All time\n\n")

	fmt.Fprintf(&w, "| Indicator|  All time aggregate |\n")
	fmt.Fprintf(&w, "|  ---:|  :--- |\n")
	fmt.Fprintf(&w, "| Number of opened issues/PRs | %d |\n", int(allTimeSeries.DailyNumMotionOpen.Total()))
	fmt.Fprintf(&w, "| Number of closed issues/PRs | %d |\n", int(allTimeSeries.DailyNumMotionClose.Total()))
	fmt.Fprintf(&w, "| Number of cancelled issues/PRs | %d |\n\n", int(allTimeSeries.DailyNumMotionCancel.Total()))

	fmt.Fprintf(&w, "| Indicator|  All time aggregate |\n")
	fmt.Fprintf(&w, "|  ---:|  :--- |\n")
	fmt.Fprintf(&w, "| Number of votes on issues | %d |\n", int(allTimeSeries.DailyNumConcernVotes.Total()))
	fmt.Fprintf(&w, "| Number of votes on PRs | %d |\n", int(allTimeSeries.DailyNumProposalVotes.Total()))
	fmt.Fprintf(&w, "| Number of votes on other | %d |\n", int(allTimeSeries.DailyNumOtherVotes.Total()))
	fmt.Fprintf(&w, "| Credits spent on issue votes | %0.6f |\n", allTimeSeries.DailyConcernVoteCharges.Total())
	fmt.Fprintf(&w, "| Credits spent on PR votes | %0.6f |\n", allTimeSeries.DailyProposalVoteCharges.Total())
	fmt.Fprintf(&w, "| Credits spent on other votes | %0.6f |\n\n", allTimeSeries.DailyOtherVoteCharges.Total())

	fmt.Fprintf(&w, "| Indicator|  All time aggregate |\n")
	fmt.Fprintf(&w, "|  ---:|  :--- |\n")
	fmt.Fprintf(&w, "| Credits cleared in bounties | %0.6f |\n", allTimeSeries.DailyClearedBounties.Total())
	fmt.Fprintf(&w, "| Credits cleared in rewards | %0.6f |\n", allTimeSeries.DailyClearedRewards.Total())
	fmt.Fprintf(&w, "| Credits cleared in refunds | %0.6f |\n\n", allTimeSeries.DailyClearedRefunds.Total())

	fmt.Fprintf(&w, "| Indicator|  All time aggregate |\n")
	fmt.Fprintf(&w, "|  ---:|  :--- |\n")
	fmt.Fprintf(&w, "| Number of new members | %d |\n\n", int(allTimeSeries.DailyNumJoins.Total()))

	fmt.Fprintf(&w, "| Indicator|  All time aggregate |\n")
	fmt.Fprintf(&w, "|  ---:|  :--- |\n")
	fmt.Fprintf(&w, "| Credits issued | %0.6f |\n", allTimeSeries.DailyCreditsIssued.Total())
	fmt.Fprintf(&w, "| Credits burned | %0.6f |\n", allTimeSeries.DailyCreditsBurned.Total())
	fmt.Fprintf(&w, "| Credits transferred | %0.6f |\n\n", allTimeSeries.DailyCreditsTransferred.Total())

	if len(capTable) > 0 {
		fmt.Printf("### Capitalization table\n\n")
		fmt.Fprintf(&w, "| User|  Capitalization |\n")
		fmt.Fprintf(&w, "|  :--- |  :--- |\n")
		for _, uc := range capTable {
			fmt.Fprintf(&w, "|  @%s |  %0.6f |\n", uc.Name, uc.Cap)
		}
		fmt.Fprintln(&w)
	}

	return &ReportAssets{
		Series: &ReportSeries{
			Last30Days: last30DaysSeries,
			AllTime:    allTimeSeries,
		},
		ReportMD: w.String(),
		Assets: map[string][]byte{
			"daily_motions.png": plotDailyMotionsPNG(ctx, last30DaysSeries),
			"daily_votes.png":   plotDailyVotesPNG(ctx, last30DaysSeries),
			"daily_charges.png": plotDailyChargesPNG(ctx, last30DaysSeries),
			"daily_cleared.png": plotDailyClearedPNG(ctx, last30DaysSeries),
			"daily_joins.png":   plotDailyJoinsPNG(ctx, last30DaysSeries),
			"daily_credits.png": plotDailyCreditsPNG(ctx, last30DaysSeries),
		},
	}
}

func loadHistory_Local(
	ctx context.Context,
	cloned gov.Cloned,
	earliest time.Time,
	latest time.Time,

) journal.Entries[*metric.Event] {

	all := metric.List_Local(ctx, cloned)
	s := journal.Entries[*metric.Event]{}
	for _, entry := range all {
		if isNotBefore(entry.Stamp, earliest) && isNotAfter(entry.Stamp, latest) {
			s = append(s, entry)
		}
	}
	return s
}
