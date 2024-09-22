package companydata

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	lib "github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

type fn func([]interface{}) interface{}

func GetMapFx(name string, value []interface{}) interface{} {

	res := map[string]fn{
		"test":                         Test,
		"formatDateDDMMYYYYSlash":      formatDateDDMMYYYYSlash,
		"formatISO8601toDDMMYYYYSlash": formatISO8601toDDMMYYYYSlash,
		"mapWorkCodeGlobal":            mapWorkCodeGlobal,
		"combineValuesWithSpace":       combineValuesWithSpace,
	}
	return res[name](value)
}
func Test(s []interface{}) interface{} {
	return ""
}
func formatDateDDMMYYYYSlash(s []interface{}) interface{} {
	//2024-09-04T00:00:00Z
	return s[0].(time.Time).Format("02/01/2006")
}
func formatISO8601toDDMMYYYYSlash(d []interface{}) interface{} {
	var res string
	log.Println("value d: ", d)
	log.Println("value d: ", reflect.TypeOf(d))
	if d[0].(string) != "" {
		splitD := strings.Split(d[0].(string), "-")
		split2 := strings.Split(splitD[2], "T")
		day, _ := strconv.Atoi(split2[0])
		month, _ := strconv.Atoi(splitD[1])
		res = fmt.Sprintf("%02d", day) + "/" + fmt.Sprintf("%02d", month) + "/" + splitD[0]
	}

	return res

}
func formatSplitPaymentNumber(s []interface{}) interface{} {
	var res string
	if s[0] == "monthly" {
		res = "1"

	}
	if s[0] == "yearly" {
		res = "12"

	}
	return res
}
func mapWorkCodeGlobal(s []interface{}) interface{} {
	var res string

	works := lib.GetFilesByEnv("enrich/work-code-global.csv")

	df := lib.CsvToDataframe(works)
	fil := df.Filter(
		dataframe.F{Colidx: 1, Colname: "Settore", Comparator: series.Eq, Comparando: s[0].(models.User).WorkType},
		dataframe.F{Colidx: 2, Colname: "Tipo", Comparator: series.Eq, Comparando: "Lavoratore " + s[0].(models.User).Work},
	)
	if fil.Nrow() > 0 {
		res = fil.Elem(0, 0).String()
	}
	return res
}
func combineValuesWithSpace(s []interface{}) interface{} {
	var res string
	for _, value := range s {
		res = res + " " + value.(string)
	}
	return res
}
func getNextPayRate(s []interface{}) interface{} {
	var (
		res     string
		resTime time.Time
	)

	t := s[0].(string) + "00:00"
	//2024-09-13T00:00:00Z
	//RFC3339	“2006-01-02T15:04:05Z07:00”
	parseTime, err := time.Parse(time.RFC3339, t)
	lib.CheckError(err)
	if s[1].(string) == "monthly" {
		resTime = parseTime.AddDate(0, 1, 0)

	} else {
		resTime = parseTime
	}
	res = resTime.Format("02/01/2006")
	return res
}
