package lib

import (
	"strconv"
	"strings"
	"time"

	"gitlab.dev.wopta.it/goworkspace/lib/log"
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

func MonthsDifference(t1, t2 time.Time) int {
	// Ensure t1 is before t2
	if t1.After(t2) {
		t1, t2 = t2, t1
	}

	// Calculate the difference in years and months
	yearDiff := t2.Year() - t1.Year()
	monthDiff := int(t2.Month()) - int(t1.Month())

	// Adjust the difference if the day of the month is smaller in t1
	if t2.Day() < t1.Day() {
		monthDiff--
	}

	// Calculate the total difference in months
	totalMonths := (yearDiff * 12) + monthDiff

	return totalMonths
}

func GetPreviousMonth(t time.Time) time.Time {
	return t.AddDate(0, -1, 0)
}

func GetFirstDay(t time.Time) time.Time {
	year, month, _ := t.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
}
func ParseDateDDMMYYYY(date string) time.Time {
	var (
		res time.Time
	)
	log.Println(" ParseDateDDMMYYYY date:", date)
	log.Println(" ParseDateDDMMYYYY len(date):", len(date))
	if len(date) < 8 {
		date = "0" + date
	}
	if len(date) >= 8 {
		d, e := strconv.Atoi(date[0:1])
		m, e := strconv.Atoi(date[2:3])
		y, e := strconv.Atoi(date[4:7])

		res = time.Date(y, time.Month(m),
			d, 0, 0, 0, 0, time.UTC)
		log.Println(e)
	}
	return res

}

func SetDateToStartOfDay(date time.Time) time.Time {
	return date.Truncate(24 * time.Hour)
}

func AddMonths(t time.Time, m int) time.Time {
	x := t.AddDate(0, m, 0)
	if d := x.Day(); d != t.Day() {
		return x.AddDate(0, 0, -d)
	}
	return x
}

func IsEqual(d1, d2 time.Time) bool {
	return d1.Truncate(time.Hour * 24).Equal(d2.Truncate(time.Hour * 24))
}

func ExtractLocalMonth(t time.Time) string {
	month := map[string]string{
		"january":   "Gennaio",
		"february":  "Febbraio",
		"march":     "Marzo",
		"april":     "Aprile",
		"may":       "Maggio",
		"june":      "Giugno",
		"july":      "Luglio",
		"august":    "Agosto",
		"september": "Settembre",
		"october":   "Ottobre",
		"november":  "Novembre",
		"december":  "Dicembre",
	}
	return month[strings.ToLower(t.Month().String())]
}
