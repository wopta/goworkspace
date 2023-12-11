package lib

import (
	"github.com/dustin/go-humanize"
	"strings"
)

func HumanaizePriceEuro(price float64) string {
	return "€ " + humanize.FormatFloat("#.###,##", price)
}

func Capitalize(s string) string {
	if len(s) == 0 {
		return ""
	}
	return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
}
