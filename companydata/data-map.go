package companydata

import (
	"fmt"
	"github.com/wopta/goworkspace/lib/log"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	lib "github.com/wopta/goworkspace/lib"
)

type fn func([]interface{}) interface{}

func GetMapFx(name string, value []interface{}) interface{} {

	res := map[string]fn{
		"test":                         Test,
		"formatDateDDMMYYYYSlash":      formatDateDDMMYYYYSlash,
		"formatISO8601toDDMMYYYYSlash": formatISO8601toDDMMYYYYSlash,
		"mapWorkCodeGlobal":            mapWorkCodeGlobal,
		"combineValuesWithSpace":       combineValuesWithSpace,
		"getNextPayDate":               getNextPayDate,
		"getNextPayRate":               getNextPayRate,
		"ifZeroEmpty":                  ifZeroEmpty,
		"dotToComma":                   dotToComma,
		"personaDeductible":            personaDeductible,
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
	if d[0].(string) != "" {
		splitD := strings.Split(d[0].(string), "-")
		split2 := strings.Split(splitD[2], "T")
		day, _ := strconv.Atoi(split2[0])
		month, _ := strconv.Atoi(splitD[1])
		res = fmt.Sprintf("%02d", day) + "/" + fmt.Sprintf("%02d", month) + "/" + splitD[0]
	}

	return res

}
func mapWorkCodeGlobal(s []interface{}) interface{} {
	var res string

	works := lib.GetFilesByEnv("enrich/work-code-global.csv")

	df := lib.CsvToDataframe(works)
	fil := df.Filter(
		dataframe.F{Colidx: 1, Colname: "Settore", Comparator: series.Eq, Comparando: "Lavoratore " + s[0].(map[string]interface{})["workType"].(string)},
	)
	fil = fil.Filter(

		dataframe.F{Colidx: 2, Colname: "Tipo", Comparator: series.Eq, Comparando: s[0].(map[string]interface{})["work"].(string)},
	)
	log.Println("fil.Nrow(): ", fil.Nrow())

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
func getNextPayDate(s []interface{}) interface{} {
	var (
		res     string
		resTime time.Time
	)

	t := s[0].(string)
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
func getNextPayRate(s []interface{}) interface{} {
	var (
		res    interface{}
		resOut interface{}
	)
	log.Println("getNextPayRate s[1].(string): ", s[1].(string))
	log.Println("getNextPayRate premiumGrossMonthly: ", s[0].(map[string]interface{})["premiumGrossMonthly"])
	log.Println("getNextPayRate premiumGrossYearly: ", s[0].(map[string]interface{})["premiumGrossYearly"])
	if s[1].(string) == "monthly" {
		res = s[0].(map[string]interface{})["premiumGrossMonthly"]
	} else {
		res = s[0].(map[string]interface{})["premiumGrossYearly"]
	}

	if reflect.TypeOf(res).String() == "float64" {
		s := fmt.Sprintf("%v", res.(float64))
		resOut = strings.Replace(s, ".", ",", -1)
	}
	return resOut
}
func ifZeroEmpty(s []interface{}) interface{} {
	var (
		res interface{}
	)

	if s[0] == 0 {
		res = " "

	} else {

		res = s[0]
	}

	return res
}
func personaDeductible(s []interface{}) interface{} {
	var (
		res interface{}
	)
	deductible := s[0].(map[string]interface{})["deductible"]
	deductibleType := s[0].(map[string]interface{})["deductibleType"]
	if deductible == "5" && deductibleType == "absorbable" {
		res = "1"
	}
	if deductible == "10" && deductibleType == "absorbable" {
		res = "2"

	}
	if deductible == "3" && deductibleType == "absolute" {
		res = "3"

	}
	if deductible == "5" && deductibleType == "absolute" {
		res = "4"

	}
	if deductible == "10" && deductibleType == "absolute" {
		res = "5"

	}

	return res
}
func dotToComma(s []interface{}) interface{} {
	var (
		res interface{}
	)

	if reflect.TypeOf(s[0]).String() == "float64" {
		s := fmt.Sprintf("%v", s[0].(float64))
		res = strings.Replace(s, ".", ",", -1)
	}

	return res
}
