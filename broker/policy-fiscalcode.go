package broker

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"

	lib "github.com/wopta/goworkspace/lib"
	models "github.com/wopta/goworkspace/models"
)

func PolicyFiscalcode(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	w.Header().Set("Access-Control-Allow-Methods", "GET")

	log.Println("GetPolicyByFiscalCode")
	log.Println(r.RequestURI)

	var policies []models.Policy
	fiscalCode := r.Header.Get("fiscalCode")
	fiscalCodeRegex, _ := regexp.Compile("^(?:[A-Z][AEIOU][AEIOUX]|[AEIOU]X{2}|[B-DF-HJ-NP-TV-Z]{2}[A-Z]){2}(?:[\\dLMNP-V]{2}(?:[A-EHLMPR-T](?:[04LQ][1-9MNP-V]|[15MR][\\dLMNP-V]|[26NS][0-8LMNP-U])|[DHPS][37PT][0L]|[ACELMRT][37PT][01LM]|[AC-EHLMPR-T][26NS][9V])|(?:[02468LNQSU][048LQU]|[13579MPRTV][26NS])B[26NS][9V])(?:[A-MZ][1-9MNP-V][\\dLMNP-V]{2}|[A-M][0L](?:[1-9MNP-V][\\dLMNP-V]|[0L][1-9MNP-V]))[A-Z]$")
	
	if !fiscalCodeRegex.Match([]byte(fiscalCode)) {
		return `{}`, nil, nil
	}

	// get all policies from firestore
	docsnap := lib.WhereFirestore("policy", "contractor.fiscalCode", "==", fiscalCode)
	policies = models.PolicyToListData(docsnap)
	

	// get all policies from wise
	// wiseDoc := wiseProxy.WiseProxyObj()

	res, _ := json.Marshal(policies)

	return string(res), policies, nil
}
