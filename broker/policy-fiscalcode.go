package broker

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	lib "github.com/wopta/goworkspace/lib"
	models "github.com/wopta/goworkspace/models"
	wiseProxy "github.com/wopta/goworkspace/wiseproxy"
)

func PolicyFiscalcode(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	w.Header().Set("Access-Control-Allow-Methods", "GET")

	var (
		policies                   []models.Policy
		wiseSimplePolicyResponse   WiseSimplePolicyResponse
		wiseCompletePolicyResponse WiseCompletePolicyResponse
		wiseProxyInputs            []WiseProxyInput
		wiseToken                  *string
		responseReader             io.ReadCloser
		jsonData                   []byte
		e                          error
	)

	log.Println("GetPolicyByFiscalCode")
	log.Println(r.RequestURI)

	fiscalCode := r.Header.Get("fiscalCode")
	fiscalCodeRegex, _ := regexp.Compile("^(?:[A-Z][AEIOU][AEIOUX]|[AEIOU]X{2}|[B-DF-HJ-NP-TV-Z]{2}[A-Z]){2}(?:[\\dLMNP-V]{2}(?:[A-EHLMPR-T](?:[04LQ][1-9MNP-V]|[15MR][\\dLMNP-V]|[26NS][0-8LMNP-U])|[DHPS][37PT][0L]|[ACELMRT][37PT][01LM]|[AC-EHLMPR-T][26NS][9V])|(?:[02468LNQSU][048LQU]|[13579MPRTV][26NS])B[26NS][9V])(?:[A-MZ][1-9MNP-V][\\dLMNP-V]{2}|[A-M][0L](?:[1-9MNP-V][\\dLMNP-V]|[0L][1-9MNP-V]))[A-Z]$")

	if !fiscalCodeRegex.Match([]byte(fiscalCode)) {
		return `{}`, nil, nil
	}

	// get all policies from firestore
	policies = GetPoliciesFromFirebase(fiscalCode)

	request := []byte(`{
		"idNodo": "1",
		"codiceFiscalePIva": "` + fiscalCode + `",
		"cdLingua": "it"
	}`)

	responseReader, wiseToken = wiseProxy.WiseBatch("WebApiProduct/Api/RicercaPolizzaCliente", request, "POST", wiseToken)
	jsonData, e = ioutil.ReadAll(responseReader)

	if e != nil {
		return "", false, e
	}

	log.Printf("%s", jsonData)
	e = json.Unmarshal(jsonData, &wiseSimplePolicyResponse)

	livePolicies := filter(wiseSimplePolicyResponse.Policies, func(pol WiseSimplePolicy) bool {
		return strings.ToUpper(pol.State) != "POLIZZA IN VITA"
	})

	for _, simplePolicy := range livePolicies {
		request = []byte(`{"idPolizza": "` + fmt.Sprint(simplePolicy.Id) + `", "cdLingua": "it"}`)
		wiseProxyInputs = append(wiseProxyInputs, WiseProxyInput{"WebApiProduct/Api/GetPolizzaCompleta", request, "POST"})
	}

	policies = lib.ExecuteInBatches(
		wiseProxyInputs,
		1,
		func(input WiseProxyInput) models.Policy {
			responseReader, wiseToken = wiseProxy.WiseBatch(input.Endpoint, input.Request, input.Method, wiseToken)
			jsonData, _ := ioutil.ReadAll(responseReader)
			fmt.Println(string(jsonData))
			fmt.Println("==============================================")

			e = json.Unmarshal(jsonData, &wiseCompletePolicyResponse)
			return wiseCompletePolicyResponse.Policy.ToDomain()
		},
	)

	res, _ := json.Marshal(policies)

	fmt.Printf("Found %d policies for this fiscal code: %s", len(policies), fiscalCode)

	return string(res), policies, nil
}

func filter[T any](ss []T, test func(T) bool) (ret []T) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

func GetPoliciesFromFirebase(fiscalCode string) []models.Policy {
	docsnap := lib.WhereFirestore("policy", "contractor.fiscalCode", "==", fiscalCode)
	return models.PolicyToListData(docsnap)
}

type WiseSimplePolicyResponse struct {
	Policies []WiseSimplePolicy `json:"listRisultatoRicerca"`
}

type WiseSimplePolicy struct {
	Id    int    `json:"idPolizza"`
	State string `json:"statoPolizza"`
}

type WiseCompletePolicyResponse struct {
	Policy models.WiseCompletePolicy `json:"polizza"`
}

type WiseProxyInput struct {
	Endpoint string
	Request  []byte
	Method   string
}
