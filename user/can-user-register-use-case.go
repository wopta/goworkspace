package user

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"golang.org/x/exp/slices"

	lib "github.com/wopta/goworkspace/lib"
	models "github.com/wopta/goworkspace/models"
	wiseProxy "github.com/wopta/goworkspace/wiseproxy"
)

func CanUserRegisterUseCase(fiscalCode string) (bool, *models.User, *string, *string) {
	var (
		wiseResponse models.WiseUserRegistryResponseDto
		err          error
		userId       *string = nil
		email        *string = nil
		canRegister  bool
		hasWiseData  bool
		woptaUser    models.User
		wiseUser     *models.User = nil
	)

	wiseUser, hasWiseData, err = userHasDataInWise(fiscalCode, &wiseResponse)
	if hasWiseData {
		log.Println("Found a register in wise for the user")
		canRegister = true
	}

	// Check for policies in Firestore

	if woptaUser, err = GetUserByFiscalCode(fiscalCode); err == nil {
		userId = &woptaUser.Uid
		email = &woptaUser.Mail
		log.Println(`Found user in Wopta DB for the user ` + *userId)
		canRegister = true
		return canRegister, &woptaUser, userId, email
	} else {
		log.Printf(`Error trying to find user in Firebase %v\n`, err)
		return canRegister, wiseUser, userId, &wiseUser.Mail
	}
}

func userHasDataInWise(fiscalCode string, wiseResponse *models.WiseUserRegistryResponseDto) (*models.User, bool, error) {
	request := []byte(`{
		"idNodo": "1",
		"cdFiscale": "` + fiscalCode + `",
		"cdLingua": "it"
	}`)

	responseReader := wiseProxy.WiseProxyObj("WebApiProduct/Api/RicercaAnagSemplice", request, "POST")
	jsonData, e := ioutil.ReadAll(responseReader)

	if e != nil {
		return nil, false, e
	}

	log.Printf("%s", jsonData)
	e = json.Unmarshal(jsonData, &wiseResponse)

	if wiseResponse == nil || len(*wiseResponse.UserRegistries) == 0 {
		return nil, false, e
	}

	idx := slices.IndexFunc((*wiseResponse.UserRegistries), func(registry models.WiseUserRegistryDto) bool { return registry.RegistryType == "PERSONA FISICA" })
	subjectId := (*wiseResponse.UserRegistries)[idx].Id
	request = []byte(`{
		"idNodo": "1",
		"cdFiscale": "` + fiscalCode + `",
		"idSoggetto": "` + fmt.Sprint(subjectId) + `",
		"cdLingua": "it"
	}`)
	responseReader = wiseProxy.WiseProxyObj("WebApiProduct/Api/RicercaAnagCompleta", request, "POST")
	jsonData, e = ioutil.ReadAll(responseReader)
	e = json.Unmarshal(jsonData, &wiseResponse)
	lib.CheckError(e)

	if len(*wiseResponse.UserRegistries) > 0 {
		idx := slices.IndexFunc((*wiseResponse.UserRegistries), func(registry models.WiseUserRegistryDto) bool { return registry.RegistryType == "PERSONA FISICA" })
		return (*wiseResponse.UserRegistries)[idx].ToDomain(), true, e
	}

	return nil, false, e
}
