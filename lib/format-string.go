package lib

import "github.com/dustin/go-humanize"

func HumanaizePriceEuro(price float64) string {
	return "€ " + humanize.FormatFloat("#.###,##", price)
}
