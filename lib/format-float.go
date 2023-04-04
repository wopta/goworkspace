package lib

import (
	"fmt"
	"strconv"
)

func RoundFloatTwoDecimals(in float64) float64 {
	res, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", in), 64)
	return res
}
