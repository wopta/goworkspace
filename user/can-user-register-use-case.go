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

func CanUserRegisterUseCase(fiscalCode string) (bool, *string, *string) {
	var (
		wiseResponse models.WiseUserRegistryResponseDto
		e            error
		userId       string
		canRegister  bool
		hasWiseData  bool
		email        string
	)

	hasWiseData, e = userHasDataInWise(fiscalCode, &wiseResponse)
	if hasWiseData {
		log.Println("Found a register in wise for the user")
		canRegister = true
	}

	// Check for policies in Firestore
	woptaUser, e := GetUserByFiscalCode(fiscalCode)

	if e == nil {
		userId = woptaUser.Uid
		email = woptaUser.Mail
		log.Println(`Found user in Wopta DB for the user ` + userId)
		canRegister = true
	} else {
		log.Printf(`Error trying to find user in Firebase %v\n`, e)
	}

	return canRegister, &userId, &email
}

func userHasDataInWise(fiscalCode string, wiseResponse *models.WiseUserRegistryResponseDto) (bool, error) {
	request := []byte(`{
		"idNodo": "1",
		"cdFiscale": "` + fiscalCode + `",
		"cdLingua": "it"
	}`)

	responseReader := wiseProxy.WiseProxyObj("WebApiProduct/Api/RicercaAnagSemplice", request, "POST")
	jsonData, e := ioutil.ReadAll(responseReader)
	
	if e != nil {
		return false, e
	}
	
	log.Printf("%s", jsonData)
	e = json.Unmarshal(jsonData, &wiseResponse)
	
	if wiseResponse == nil || len(*wiseResponse.UserRegistries) == 0 {
		return false, e
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
	log.Printf("%s", jsonData)
	lib.CheckError(e)

	return len(*wiseResponse.UserRegistries) > 0, e
}
