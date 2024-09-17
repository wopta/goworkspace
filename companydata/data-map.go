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

type fn func(interface{}) interface{}

func GetMapFx(name string, value interface{}) interface{} {

	res := map[string]fn{
		"test":                         Test,
		"formatDateDDMMYYYYSlash":      formatDateDDMMYYYYSlash,
		"formatISO8601toDDMMYYYYSlash": formatISO8601toDDMMYYYYSlash,
		"mapWorkCodeGlobal":            mapWorkCodeGlobal,
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
func formatISO8601toDDMMYYYYSlash(d interface{}) interface{} {
	var res string
	log.Println("value d: ", d)
	log.Println("value d: ", reflect.TypeOf(d))
	if d.(string) != "" {
		splitD := strings.Split(d.(string), "-")
		split2 := strings.Split(splitD[2], "T")
		day, _ := strconv.Atoi(split2[0])
		month, _ := strconv.Atoi(splitD[1])
		res = fmt.Sprintf("%02d", day) + "/" + fmt.Sprintf("%02d", month) + "/" + splitD[0]
	}

	return res

}
func formatSplitPaymentNumber(s interface{}) interface{} {
	var res string
	if s == "monthly" {
		res = "1"

	}
	if s == "yearly" {
		res = "12"

	}
	return res
}
func mapWorkCodeGlobal(s interface{}) interface{} {
	var res string

	works := lib.GetFilesByEnv("enrich/works-code-global.csv")

	df := lib.CsvToDataframe(works)
	fil := df.Filter(
		dataframe.F{Colidx: 1, Colname: "Settore", Comparator: series.Eq, Comparando: s.(models.User).WorkType},
		dataframe.F{Colidx: 2, Colname: "Tipo", Comparator: series.Eq, Comparando: s.(models.User).Work},
	)
	if fil.Nrow() > 0 {
		res = fil.Elem(0, 0).String()
	}
	return res
}
