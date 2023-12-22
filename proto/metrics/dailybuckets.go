package metrics

import (
	"sort"
	"time"
)

type DailyBuckets map[time.Time]float64

func (db DailyBuckets) Add(t time.Time, v float64) {
	t = Dailify(t)
	u, _ := db[t]
	db[t] = u + v
}

func (db DailyBuckets) Earliest(anchor time.Time) time.Time {
	earliest := anchor
	for t := range db {
		earliest = minTime(earliest, t)
	}
	return earliest
}

func (db DailyBuckets) Latest(anchor time.Time) time.Time {
	latest := anchor
	for t := range db {
		latest = maxTime(latest, t)
	}
	return latest
}

func isNotBefore(q, earliest time.Time) bool {
	return !q.Before(earliest)
}

func isNotAfter(q, latest time.Time) bool {
	return !q.After(latest)
}

func (db DailyBuckets) XY(earliest, latest time.Time) (ds DailySeries) {

	// backfill missing days
	for t := earliest; isNotAfter(t, latest); t = t.AddDate(0, 0, 1) {
		if _, ok := db[t]; !ok {
			db[t] = 0.0
		}
	}

	// order
	sv := make(stampedValues, 0, len(db))
	for t, v := range db {
		if isNotBefore(t, earliest) && isNotAfter(t, latest) {
			sv = append(sv, stampedValue{Stamp: t, Value: v})
		}
	}
	sv.Sort()

	// produce plot data
	ds.X, ds.Y = make([]time.Time, len(sv)), make([]float64, len(sv))
	for i := range sv {
		ds.X[i] = sv[i].Stamp
		ds.Y[i] = sv[i].Value
	}

	return ds
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
