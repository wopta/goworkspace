package policy

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/wiseproxy"
)

func GetPolicyByFiscalCodeFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		policies           []models.Policy
		wiseToken          *string = nil
		e                  error
		wiseSimplePolicies *[]WiseSimplePolicy
		response           GetPolicesByFiscalCodeResponse
	)

	log.SetPrefix("[GetPolicyByFiscalCodeFx] ")
	defer log.SetPrefix("")

	log.Println("Handler start -----------------------------------------------")

	policyFire := lib.GetDatasetByEnv(r.Header.Get("origin"), lib.PolicyCollection)
	fiscalCode := chi.URLParam(r, "fiscalcode")
	fiscalCodeRegex, _ := regexp.Compile("^(?:[A-Z][AEIOU][AEIOUX]|[AEIOU]X{2}|[B-DF-HJ-NP-TV-Z]{2}[A-Z]){2}(?:[\\dLMNP-V]{2}(?:[A-EHLMPR-T](?:[04LQ][1-9MNP-V]|[15MR][\\dLMNP-V]|[26NS][0-8LMNP-U])|[DHPS][37PT][0L]|[ACELMRT][37PT][01LM]|[AC-EHLMPR-T][26NS][9V])|(?:[02468LNQSU][048LQU]|[13579MPRTV][26NS])B[26NS][9V])(?:[A-MZ][1-9MNP-V][\\dLMNP-V]{2}|[A-M][0L](?:[1-9MNP-V][\\dLMNP-V]|[0L][1-9MNP-V]))[A-Z]$")

	if !fiscalCodeRegex.Match([]byte(fiscalCode)) {
		return `{}`, nil, nil
	}

	policies = GetPoliciesFromFirebase(fiscalCode, policyFire)

	wiseToken, wiseSimplePolicies, e = getAllSimplePoliciesForUserFromWise(fiscalCode)

	if e != nil {
		return "{}", nil, nil
	}

	wisePolicies := getCompletePoliciesFromWise(*wiseSimplePolicies, wiseToken)
	policies = append(policies, wisePolicies...)

	response.Policies = policies
	res, _ := json.Marshal(response)

	fmt.Printf("Found %d policies for this fiscal code: %s", len(policies), fiscalCode)

	log.Println("Handler end -------------------------------------------------")

	return string(res), response, nil
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
		2,
		func(input WiseProxyInput) models.Policy {
			responseReader, wiseToken = wiseproxy.WiseBatch(input.Endpoint, input.Request, input.Method, wiseToken)
			jsonData, _ := io.ReadAll(responseReader)

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

	responseReader, wiseToken = wiseproxy.WiseBatch("WebApiProduct/Api/RicercaPolizzaCliente", request, "POST", nil)
	jsonData, e := io.ReadAll(responseReader)

	if e != nil {
		return nil, nil, e
	}

	e = json.Unmarshal(jsonData, &wiseSimplePolicyResponse)

	livePolicies := lib.SliceFilter(wiseSimplePolicyResponse.Policies, func(pol WiseSimplePolicy) bool {
		return strings.ToUpper(pol.State) == "POLIZZA IN VITA"
	})

	return wiseToken, &livePolicies, e
}

func GetPoliciesFromFirebase(fiscalCode string, policyFire string) []models.Policy {
	q := lib.Firequeries{
		Queries: []lib.Firequery{
			{
				Field:      "contractor.fiscalCode",
				Operator:   "==",
				QueryValue: fiscalCode,
			},
			{
				Field:      "companyEmit",
				Operator:   "==",
				QueryValue: true,
			},
			{
				Field:      "isPay",
				Operator:   "==",
				QueryValue: true,
			},
			{
				Field:      "isSign",
				Operator:   "==",
				QueryValue: true,
			},
		},
	}
	docsnap, _ := q.FirestoreWherefields(policyFire)
	return models.PolicyToListData(docsnap)
}

type GetPolicesByFiscalCodeResponse struct {
	Policies []models.Policy `json:"policies"`
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
