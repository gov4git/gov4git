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
	since time.Time,

) journal.Entries[*history.Event] {

	all := history.List_Local(ctx, cloned)
	s := journal.Entries[*history.Event]{}
	for _, entry := range all {
		if !since.After(entry.Stamp) {
			s = append(s, entry)
		}
	}
	return s
}

func ComputeMetrics_Local(
	ctx context.Context,
	cloned gov.Cloned,
	since time.Time,

) *Series {

	entries := LoadHistorySince_Local(ctx, cloned, since)
	return ComputeSeries(entries)
}
