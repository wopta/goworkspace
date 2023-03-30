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
		policies           []models.Policy
		wiseToken          *string = nil
		e                  error
		wiseSimplePolicies *[]WiseSimplePolicy
	)

	log.Println("GetPolicyByFiscalCode")
	log.Println(r.RequestURI)

	fiscalCode := r.Header.Get("fiscalCode")
	fiscalCodeRegex, _ := regexp.Compile("^(?:[A-Z][AEIOU][AEIOUX]|[AEIOU]X{2}|[B-DF-HJ-NP-TV-Z]{2}[A-Z]){2}(?:[\\dLMNP-V]{2}(?:[A-EHLMPR-T](?:[04LQ][1-9MNP-V]|[15MR][\\dLMNP-V]|[26NS][0-8LMNP-U])|[DHPS][37PT][0L]|[ACELMRT][37PT][01LM]|[AC-EHLMPR-T][26NS][9V])|(?:[02468LNQSU][048LQU]|[13579MPRTV][26NS])B[26NS][9V])(?:[A-MZ][1-9MNP-V][\\dLMNP-V]{2}|[A-M][0L](?:[1-9MNP-V][\\dLMNP-V]|[0L][1-9MNP-V]))[A-Z]$")

	if !fiscalCodeRegex.Match([]byte(fiscalCode)) {
		return `{}`, nil, nil
	}

	policies = GetPoliciesFromFirebase(fiscalCode)

	wiseToken, wiseSimplePolicies, e = getAllSimplePoliciesForUserFromWise(fiscalCode)

	if e != nil {
		return "{}", nil, nil
	}

	wisePolicies := getCompletePoliciesFromWise(*wiseSimplePolicies, wiseToken)
	policies = append(policies, wisePolicies...)

	res, _ := json.Marshal(policies)

	fmt.Printf("Found %d policies for this fiscal code: %s", len(policies), fiscalCode)

	return string(res), policies, nil
}

func getCompletePoliciesFromWise(simplePolicies []WiseSimplePolicy, wiseToken *string) []models.Policy {
	var (
		wiseCompletePolicyResponse WiseCompletePolicyResponse
		wiseProxyInputs            []WiseProxyInput
		request                    []byte
		responseReader             io.ReadCloser
	)

	for _, simplePolicy := range simplePolicies {
		request = []byte(`{"idPolizza": "` + fmt.Sprint(simplePolicy.Id) + `", "cdLingua": "it"}`)
		wiseProxyInputs = append(wiseProxyInputs, WiseProxyInput{"WebApiProduct/Api/GetPolizzaCompleta", request, "POST"})
	}

	return lib.ExecuteInBatches(
		wiseProxyInputs,
		1,
		func(input WiseProxyInput) models.Policy {
			responseReader, wiseToken = wiseProxy.WiseBatch(input.Endpoint, input.Request, input.Method, wiseToken)
			jsonData, _ := ioutil.ReadAll(responseReader)

			_ = json.Unmarshal(jsonData, &wiseCompletePolicyResponse)
			return wiseCompletePolicyResponse.Policy.ToDomain()
		},
	)
}

func getAllSimplePoliciesForUserFromWise(fiscalCode string) (*string, *[]WiseSimplePolicy, error) {
	var (
		wiseToken                *string
		responseReader           io.ReadCloser
		wiseSimplePolicyResponse WiseSimplePolicyResponse
	)

	request := []byte(`{
		"idNodo": "1",
		"codiceFiscalePIva": "` + fiscalCode + `",
		"cdLingua": "it"
	}`)

	responseReader, wiseToken = wiseProxy.WiseBatch("WebApiProduct/Api/RicercaPolizzaCliente", request, "POST", nil)
	jsonData, e := ioutil.ReadAll(responseReader)

	if e != nil {
		return nil, nil, e
	}

	e = json.Unmarshal(jsonData, &wiseSimplePolicyResponse)

	livePolicies := filter(wiseSimplePolicyResponse.Policies, func(pol WiseSimplePolicy) bool {
		return strings.ToUpper(pol.State) == "POLIZZA IN VITA"
	})

	return wiseToken, &livePolicies, e
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
	q := lib.Firequeries{
		Queries: []lib.Firequery{
			{
				Field:      "contractor.fiscalCode",
				Operator:   "==",
				QueryValue: fiscalCode,
			},
			{
				Field:      "status",
				Operator:   "==",
				QueryValue: models.PolicyStatusCompanyEmit,
			},
		},
	}
	docsnap, _ := q.FirestoreWherefields("policy")
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
