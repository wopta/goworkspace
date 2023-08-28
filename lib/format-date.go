package lib

import (
	"fmt"
	"strconv"
	"time"
)

func FormatDate01032008(d time.Time) string {
	var res string
	res = fmt.Sprintf("%02d", d.Day()) + fmt.Sprintf("%02d", int(d.Month())) + strconv.Itoa(d.Year())
	return res

}
