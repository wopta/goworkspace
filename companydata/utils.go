package companydata

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"time"

	lib "github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
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
func ExtractUserDataFromFiscalCode(fiscalCode string) (string, models.User, error) {
	var (
		codes map[string]map[string]string
		user  = models.User{}
	)

	log.Println("Decode")

	if len(fiscalCode) < 15 {
		return "", models.User{}, fmt.Errorf("invalid fiscal code")
	}

	b, err := os.ReadFile(lib.GetAssetPathByEnv("companyData") + "/reverse-codes.json")
	lib.CheckError(err)

	err = json.Unmarshal(b, &codes)
	lib.CheckError(err)

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
