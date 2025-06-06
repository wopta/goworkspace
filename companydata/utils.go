package companydata

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"github.com/xuri/excelize/v2"
	lib "gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
)

func CheckStructNil[T interface{}](s interface{}) T {
	var result T
	result1 := new(T)
	result = *result1
	//log.Println(reflect.TypeOf(s))
	if reflect.TypeOf(s) != nil {
		log.Println("is not nill")
		result = s.(T)
	}
	log.Println(s)
	log.Println(result)
	return result
}
func StringMapping(v string, m map[string]string) string {
	return m[v]
}
func ExtractUserDataFromFiscalCode(fiscalCode string, codes map[string]map[string]string) (string, models.User, error) {
	var (
		user = models.User{}
	)
	user.FiscalCode = fiscalCode
	if len(fiscalCode) < 15 {
		return "", models.User{}, fmt.Errorf("invalid fiscal code")
	}

	day, _ := strconv.Atoi(fiscalCode[9:11])

	if day > 40 {
		user.Gender = "F"
	} else {
		user.Gender = "M"
	}

	birthPlaceCode := fiscalCode[11:15]
	if birthPlaceCode == "" {
		return "", models.User{}, fmt.Errorf("invalid birth place code")
	}
	user.BirthCity = codes[birthPlaceCode]["city"]
	user.BirthProvince = codes[birthPlaceCode]["province"]

	user.BirthDate = lib.ExtractBirthdateFromItalianFiscalCode(user.FiscalCode).Format(time.RFC3339)

	outJson, err := json.Marshal(&user)
	lib.CheckError(err)

	return string(outJson), user, nil
}
func CreateExcel(sheet [][]string, filePath string, sheetname string) ([]byte, error) {
	log.Println("CreateExcel")
	f := excelize.NewFile()
	alfabet := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
		"AA", "AB", "AC", "AD", "AE", "AF", "AG", "AH", "AI", "AJ", "AK", "AL", "AM", "AN", "AO", "AP", "AQ", "AR", "AS", "AT", "AU", "AV", "AW", "AX", "AY", "AZ",
		"BA", "BB", "BC", "BD", "BE", "BF", "BG", "BH", "BI", "BJ", "BK", "BL", "BM", "BN", "BO", "BP", "BQ", "BR", "BS", "BT", "BU", "BV", "BW", "BX", "BY", "BZ"}
	// Create a new sheet.

	index, err := f.NewSheet(sheetname)
	lib.CheckError(err)
	for x, row := range sheet {
		for i, cel := range row {
			f.SetCellValue(sheetname, alfabet[i]+""+strconv.Itoa(x+1), cel)
		}
	}
	//Set active sheet of the workbook.
	f.SetActiveSheet(index)
	err = f.DeleteSheet("Sheet1")
	lib.CheckError(err)
	//Save spreadsheet by the given path.
	err = f.SaveAs(filePath)
	lib.CheckError(err)
	resByte, err := f.WriteToBuffer()

	return resByte.Bytes(), err
}
func getRequestData(req []byte) (time.Time, bool) {
	var (
		obj    DataReq
		upload bool
	)

	json.Unmarshal([]byte(req), &obj)

	now := time.Now()

	if obj.Day == "" {
		now = time.Now()
		upload = true
	} else {
		date, _ := time.Parse("2006-01-02", obj.Day)
		now = date
		upload = obj.Upload
	}
	log.Println("end CreateExcel")
	return now, upload
}
func getCompanyDataReq(req []byte) (time.Time, bool, DataReq) {
	var (
		obj    DataReq
		upload bool
	)

	json.Unmarshal([]byte(req), &obj)

	now := time.Now()

	if obj.Day == "" {
		now = time.Now()
		upload = true
	} else {
		date, _ := time.Parse("2006-01-02", obj.Day)
		now = date
		upload = obj.Upload
	}
	return now, upload, obj
}
