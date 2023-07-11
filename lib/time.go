package lib

import (
	"time"
)

func Dateformat(t time.Time) string {
	layout := "02/01/2006"
	return t.Format(layout)
}

// computes the age/elapsed years between t1, and t2.
func ElapsedYears(t1 time.Time, t2 time.Time) int {
	if t1.After(t2) {
		t1, t2 = t2, t1
	}

	t1y, t1m, t1d := t1.Date()
	date1 := time.Date(t1y, t1m, t1d, 0, 0, 0, 0, time.UTC)

	t2y, t2m, t2d := t2.Date()
	date2 := time.Date(t2y, t2m, t2d, 0, 0, 0, 0, time.UTC)

	years := t2y - t1y
	anniversary := date1.AddDate(years, 0, 0)
	if anniversary.After(date2) {
		years--
	}

	return years
}
