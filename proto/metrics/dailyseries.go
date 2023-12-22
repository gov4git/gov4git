package metrics

import "time"

type DailySeries struct {
	X []time.Time
	Y []float64
}

func (ds DailySeries) Len() int {
	return len(ds.X)
}

func (ds DailySeries) Total() float64 {
	t := 0.0
	for _, v := range ds.Y {
		t += v
	}
	return t
}
