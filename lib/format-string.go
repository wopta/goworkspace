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

func TrimSpace(s string) string {
	return strings.TrimSpace(s)
}

func ToUpper(s string) string {
	return strings.ToUpper(TrimSpace(s))
}

func ToLower(s string) string {
	return strings.ToLower(TrimSpace(s))
}

func ReplaceAll(s, old, new string) string {
	return strings.ReplaceAll(s, old, new)
}
