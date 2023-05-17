package lib

import (
	"strconv"
	"strings"
	"time"
)

func ExtractBirthdateFromItalianFiscalCode(fiscalCode string) time.Time {
	year, _ := strconv.Atoi(fiscalCode[6:8])
	month := getMonth(fiscalCode[8:9])
	day, _ := strconv.Atoi(fiscalCode[9:11])

	if day > 40 {
		day -= 40
	}

	if year < time.Now().Year()-2000 {
		year += 2000
	} else {
		year += 1900
	}

	birthdate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	return birthdate
}

func getMonth(monthCode string) int {
	monthMap := map[string]int{
		"A": 1,
		"B": 2,
		"C": 3,
		"D": 4,
		"E": 5,
		"H": 6,
		"L": 7,
		"M": 8,
		"P": 9,
		"R": 10,
		"S": 11,
		"T": 12,
	}

	return monthMap[strings.ToUpper(monthCode)]
}
