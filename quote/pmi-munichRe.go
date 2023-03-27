package quote

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/civil"
	lib "github.com/wopta/goworkspace/lib"
)

func PmiMunichFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	jsonData, err := ioutil.ReadAll(r.Body)
	res := <-PmiMunich(jsonData)
	return res, nil, err

}
func PmiMunich(r []byte) <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		var urlstring = os.Getenv("MUNICHREBASEURL") + "/api/quote/rate/"
		client := lib.ClientCredentials(os.Getenv("MUNICHRECLIENTID"),
			os.Getenv("MUNICHRECLIENTSECRET"), os.Getenv("MUNICHRESCOPE"), os.Getenv("MUNICHRETOKENENDPOINT"))
		req, _ := http.NewRequest(http.MethodPost, urlstring, bytes.NewBuffer(r))
		req.Header.Set("Ocp-Apim-Subscription-Key", os.Getenv("MUNICHRESUBSCRIPTIONKEY"))
		req.Header.Set("Content-Type", "application/json")
		res, err := client.Do(req)
		lib.CheckError(err)

		if res != nil {
			body, err := ioutil.ReadAll(res.Body)
			lib.CheckError(err)
			res.Body.Close()
			log.Println("quote res")
			ch <- string(body)
			if res.StatusCode == 500 {
				log.Println("StatusCode == 500")

			}

		}

	}()
	return ch
}

func (r *MunichReQuoteRequest) Unmarshal(data []byte) (MunichReQuoteRequest, error) {
	var res MunichReQuoteRequest
	err := json.Unmarshal(data, &r)
	return res, err
}

func (r *MunichReQuoteRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *MunichReQuoteResponse) Unmarshal(data []byte) (MunichReQuoteResponse, error) {
	var res MunichReQuoteResponse
	err := json.Unmarshal(data, &res)
	return res, err
}

func (r *MunichReQuoteResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type MunichReQuoteRequest struct {
	SME SME `json:"sme"`
}

type SME struct {
	Ateco        string  `json:"ateco"`
	SubproductID int64   `json:"subproductId"`
	UWRole       string  `json:"UW_role"`
	Company      Company `json:"company"`
	Answers      Answers `json:"answers"`
}

type Answers struct {
	Step1 []Step1 `json:"step1"`
	Step2 []Step2 `json:"step2"`
}

type Step1 struct {
	Slug  string `json:"slug"`
	Value Value  `json:"value"`
}

type Step1Value struct {
	TypeOfSumInsured           *TypeOfSumInsured `json:"typeOfSumInsured,omitempty"`
	Deductible                 *string           `json:"deductible,omitempty"`
	SumInsuredLimitOfIndemnity *float64          `json:"sumInsuredLimitOfIndemnity,omitempty"`
	SelfInsurance              *string           `json:"selfInsurance,omitempty"`
	LegalDefence               *string           `json:"legalDefence,omitempty"`
}

type Step2 struct {
	BuildingID string     `json:"buildingId"`
	Value      Step2Value `json:"value"`
}

type Step2Value struct {
	BuildingType     string   `json:"buildingType"`
	NumberOfFloors   string   `json:"numberOfFloors"`
	ConstructionYear string   `json:"constructionYear"`
	Alarm            string   `json:"alarm"`
	TypeOfInsurance  string   `json:"typeOfInsurance"`
	Ateco            string   `json:"ateco"`
	Postcode         string   `json:"postcode"`
	Province         string   `json:"province"`
	Answer           []Answer `json:"answer"`
}

type Answer struct {
	Slug  string `json:"slug"`
	Value Value  `json:"value"`
}

type AnswerValue struct {
	TypeOfSumInsured           *TypeOfSumInsured `json:"typeOfSumInsured,omitempty"`
	Deductible                 *string           `json:"deductible,omitempty"`
	SumInsuredLimitOfIndemnity *float64          `json:"sumInsuredLimitOfIndemnity,omitempty"`
	SelfInsurance              *string           `json:"selfInsurance,omitempty"`
	Assistance                 *string           `json:"assistance,omitempty"`
	DailyAllowance             *string           `json:"dailyAllowance,omitempty"`
	LegalDefence               *string           `json:"legalDefence,omitempty"`
}

type Company struct {
	Country   string    `json:"country"`
	Vatnumber Vatnumber `json:"vatnumber"`
	OpreEur   Employees `json:"opre_eur"`
	Employees Employees `json:"employees"`
}

type Employees struct {
	Value int64 `json:"value"`
}

type Vatnumber struct {
	Value string `json:"value"`
}

type TypeOfSumInsured string

const (
	FirstLoss TypeOfSumInsured = "firstLoss"
)

type Value struct {
	TypeOfSumInsured           *string  `json:"typeOfSumInsured,omitempty"`
	Deductible                 *string  `json:"deductible,omitempty"`
	SumInsuredLimitOfIndemnity *float64 `json:"sumInsuredLimitOfIndemnity,omitempty"`
	SelfInsurance              *string  `json:"selfInsurance,omitempty"`
	LegalDefence               *string  `json:"legalDefence,omitempty"`
	Assistance                 *string  `json:"assistance,omitempty"`
	DailyAllowance             *string  `json:"dailyAllowance,omitempty"`
}
type MunichReQuoteResponse struct {
	Result Result `json:"result"`
}

type Result struct {
	Answers AnswersResponse `json:"answers"`
}
type AnswersResponse struct {
	Step1 []Step1Response `json:"step1"`
	Step2 []Step2Response `json:"step2"`
}

type Step1Response struct {
	Slug  string        `json:"slug"`
	Value ValueResponse `json:"value"`
}

type ValueResponse struct {
	PremiumNet       float64 `json:"premiumNet"`
	PremiumTaxAmount float64 `json:"premiumTaxAmount"`
	PremiumGross     float64 `json:"premiumGross"`
}

type Step2Response struct {
	BuildingID string          `json:"buildingId"`
	Value      []Step1Response `json:"value"`
}
type MunichReQuotePmiDWHCall struct {
	CreationDate      civil.DateTime ` bigquery:"creationDate"`
	Status            int64          ` bigquery:"status"`
	RequestRules      string         ` bigquery:"requestRules"`
	ResponseRules     string         ` bigquery:"responseRules"`
	RequestQuote      string         ` bigquery:"requestQuote"`
	ResponseQuote     string         ` bigquery:"responseQuote"`
	RequestRulesJson  string         ` bigquery:"requestRulesJson"`
	ResponseRulesJson string         ` bigquery:"responseRulesJson"`
	RequestQuoteJson  string         ` bigquery:"requestQuoteJson"`
	ResponseQuoteJson string         ` bigquery:"responseQuoteJson"`
	Your              int64          ` bigquery:"your"`
	Base              int64          ` bigquery:"base"`
	Premium           int64          ` bigquery:"premium"`
}
