package user

import (
	"encoding/json"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"os"
	"time"
)

func extractUserDataFromFiscalCode(user models.User) (string, models.User) {
	var (
		codes map[string]map[string]string
	)

	b, err := os.ReadFile(lib.GetAssetPathByEnv("user") + "/reverse.json")
	lib.CheckError(err)

	err = json.Unmarshal(b, &codes)
	lib.CheckError(err)

	birthPlaceCode := user.FiscalCode[11:15]
	user.BirthCity = codes[birthPlaceCode]["city"]
	user.BirthProvince = codes[birthPlaceCode]["province"]

	user.BirthDate = lib.ExtractBirthdateFromItalianFiscalCode(user.FiscalCode).Format(time.RFC3339)

	outJson, err := json.Marshal(&user)
	lib.CheckError(err)

	return string(outJson), user
}
