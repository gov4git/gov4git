package metrics

import (
	"sort"
	"time"
)

type DailyBuckets map[time.Time]float64

func (db DailyBuckets) Add(t time.Time, v float64) {
	t = truncateDay(t)
	u, _ := db[t]
	db[t] = u + v
}

func (db DailyBuckets) XY() (ds DailySeries) {

	if len(db) == 0 {
		return DailySeries{}
	}

	// compute first and last day
	var earliest, latest time.Time
	var earliestHave, latestHave bool
	for t := range db {
		if earliestHave {
			earliest = minTime(t, earliest)
		} else {
			earliest = t
			earliestHave = true
		}
		if latestHave {
			latest = maxTime(t, latest)
		} else {
			latest = t
			latestHave = true
		}
	}

	// backfill missing days
	for t := earliest; !t.After(latest); t.AddDate(0, 0, 1) {
		if _, ok := db[t]; !ok {
			db[t] = 0.0
		}
	}

	// order
	sv := make(stampedValues, 0, len(db))
	for t, v := range db {
		sv = append(sv, stampedValue{Stamp: t, Value: v})
	}
	sv.Sort()

	// produce plot data
	ds.X, ds.Y = make([]time.Time, len(sv)), make([]float64, len(sv))
	ds.Total = 0.0
	for i := range sv {
		ds.X[i] = sv[i].Stamp
		ds.Y[i] = sv[i].Value
		ds.Total += sv[i].Value
	}

	return ds
}

func truncateDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func minTime(p, q time.Time) time.Time {
	if p.Before(q) {
		return p
	}
	return q
}

func maxTime(p, q time.Time) time.Time {
	if p.After(q) {
		return p
	}
	return q
}

type stampedValue struct {
	Stamp time.Time
	Value float64
}

type stampedValues []stampedValue

func (x stampedValues) Len() int {
	return len(x)
}

func (x stampedValues) Less(i, j int) bool {
	return x[i].Stamp.Before(x[j].Stamp)
}

func (x stampedValues) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

func (x stampedValues) Sort() {
	sort.Sort(x)
}
