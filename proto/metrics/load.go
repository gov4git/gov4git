package metrics

import (
	"context"
	"time"

	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/history"
	"github.com/gov4git/gov4git/v2/proto/journal"
)

func LoadHistorySince_Local(
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

func ComputeMetrics_Local(
	ctx context.Context,
	cloned gov.Cloned,
	earliest time.Time,
	latest time.Time,

) *Series {

	entries := LoadHistorySince_Local(ctx, cloned, earliest, latest)
	return ComputeSeries(entries, earliest, latest)
}
