package metrics

import "time"

var TimeDailyLowerBound = time.Date(2023, 10, 9, 0, 0, 0, 0, time.UTC)

func Today() time.Time {
	return Dailify(time.Now())
}

func Dailify(t time.Time) time.Time {
	t = t.UTC()
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}
