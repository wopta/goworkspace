package user

import (
	"encoding/json"
	"fmt"
	"io"
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
		hasWiseData  bool
		woptaUser    models.User
		wiseUser     *models.User = nil
	)

	if woptaUser, err = GetUserByFiscalCode(fiscalCode); err == nil {
		userId = &woptaUser.Uid
		email = &woptaUser.Mail
		log.Println(`Found user in Wopta DB for the user ` + *userId)
		return true, &woptaUser, userId, email
	}

	wiseUser, hasWiseData, err = userHasDataInWise(fiscalCode, &wiseResponse)
	if hasWiseData && err == nil {
		log.Println("Found a register in wise for the user")
		return true, wiseUser, userId, &wiseUser.Mail
	}

	if err != nil {
		log.Printf(`Cannot register: %s`, err.Error())
		return false, nil, nil, nil
	}

	log.Println(`User cannot register: no data found`)
	return false, nil, nil, nil
}

func userHasDataInWise(fiscalCode string, wiseResponse *models.WiseUserRegistryResponseDto) (*models.User, bool, error) {
	request := []byte(`{
		"idNodo": "1",
		"cdFiscale": "` + fiscalCode + `",
		"cdLingua": "it"
	}`)

	responseReader := wiseProxy.WiseProxyObj("WebApiProduct/Api/RicercaAnagSemplice", request, "POST")
	jsonData, e := io.ReadAll(responseReader)

	if e != nil {
		return nil, false, e
	}

	log.Printf("%s", jsonData)
	e = json.Unmarshal(jsonData, &wiseResponse)

	if wiseResponse == nil || len(*wiseResponse.UserRegistries) == 0 {
		return nil, false, e
	}

	idx := slices.IndexFunc((*wiseResponse.UserRegistries), func(registry models.WiseUserRegistryDto) bool { return registry.RegistryType == "PERSONA FISICA" })
	if idx == -1 {
		return nil, false, e
	}

	subjectId := (*wiseResponse.UserRegistries)[idx].Id
	request = []byte(`{
		"idNodo": "1",
		"cdFiscale": "` + fiscalCode + `",
		"idSoggetto": "` + fmt.Sprint(subjectId) + `",
		"cdLingua": "it"
	}`)
	responseReader = wiseProxy.WiseProxyObj("WebApiProduct/Api/RicercaAnagCompleta", request, "POST")
	jsonData, e = io.ReadAll(responseReader)
	e = json.Unmarshal(jsonData, &wiseResponse)
	lib.CheckError(e)

	if len(*wiseResponse.UserRegistries) > 0 {
		idx := slices.IndexFunc((*wiseResponse.UserRegistries), func(registry models.WiseUserRegistryDto) bool { return registry.RegistryType == "PERSONA FISICA" })
		return (*wiseResponse.UserRegistries)[idx].ToDomain(), true, e
	}

	return nil, false, e
}
