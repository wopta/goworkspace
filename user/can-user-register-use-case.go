package user

import (
	"encoding/json"
	"io/ioutil"
	"log"

	lib "github.com/wopta/goworkspace/lib"
	models "github.com/wopta/goworkspace/models"
	wiseProxy "github.com/wopta/goworkspace/wiseproxy"
)

func CanUserRegisterUseCase(fiscalCode string) (bool, *string) {
	var (
		wiseResponse models.WiseUserRegistryResponseDto
		e            error
		userId       string
		canRegister  bool
	)

	request := []byte(`{
		"idNodo": "1",
		"cdFiscale": "` + fiscalCode + `",
		"cdLingua": "it"
	}`)

	responseReader := wiseProxy.WiseProxyObj("WebApiProduct/Api/RicercaAnagSemplice", request, "POST")
	jsonData, e := ioutil.ReadAll(responseReader)
	lib.CheckError(e)

	e = json.Unmarshal(jsonData, &wiseResponse)
	lib.CheckError(e)
	if len(*wiseResponse.UserRegistries) > 0 {
		log.Println("Found a register in wise for the user")
		canRegister = true
	}

	// Check for policies in Firestore
	rn := lib.WhereLimitFirestore("policy", "contractor.fiscalCode", "==", fiscalCode, 1)
	lib.CheckError(e)
	policies := models.PolicyToListData(rn)
	if len(policies) > 0 {
		userId = policies[0].Contractor.Uid
		log.Println(`Found a policy in Wopta DB for the user ` + userId)
		canRegister = true
	}

	return canRegister, &userId
}
