package companydata

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type fn func(interface{}) interface{}

func GetMapFx(name string, value interface{}) interface{} {

	res := map[string]fn{
		"test": Test,
		"formatDateDDMMYYYYSlash": formatDateDDMMYYYYSlash,
		"formatBithdateDDMMYYYYSlash": formatISO8601toDDMMYYYYSlash,
	}
	return res[name](value)
}
func Test(s interface{}) interface{} {
	return ""
}
func formatDateDDMMYYYYSlash(s interface{}) interface{} {
//2024-09-04T00:00:00Z
	return s.(time.Time).Format("02/01/2006")
}
func formatISO8601toDDMMYYYYSlash(d interface{}) interface{}{
	var res string
	if d.(string) != "" {
		splitD := strings.Split(d.(string), "-")
		split2 := strings.Split(splitD[2], "T")
		day, _ := strconv.Atoi(split2[0])
		month, _ := strconv.Atoi(splitD[1])
		res = fmt.Sprintf("%02d", day) + "/" + fmt.Sprintf("%02d", month)+ "/" + splitD[0]
	}

	return res

}
func formatSplitPaymentNumber(s interface{}) interface{} {
	var res string
	if s=="monthly"{
		res="1"

	}
	if s=="yearly"{
		res="12"

	}
	return res
}