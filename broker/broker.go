package broker

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	lib "github.com/wopta/goworkspace/lib"
	models "github.com/wopta/goworkspace/models"
)

func init() {
	log.Println("INIT Broker")
	functions.HTTP("Broker", Broker)
}

func Broker(w http.ResponseWriter, r *http.Request) {

	log.Println("Callback")
	lib.EnableCors(&w, r)
	//w.Header().Set("Access-Control-Allow-Methods", "POST")
	route := lib.RouteData{
		Routes: []lib.Route{

			{
				Route:   "/v1/policy/proposal",
				Hendler: Proposal,
			},
			{
				Route:   "/v1/policy/emit",
				Hendler: Emit,
			},
		},
	}
	route.Router(w, r)

}
func Proposal(w http.ResponseWriter, r *http.Request) (string, interface{}) {
	policy := models.Policy{
		ID:            "",
		IdSign:        "",
		IdPay:         "",
		Uid:           "",
		Number:        "",
		NumberCompany: "",
		Status:        "",
		StatusHistory: []string{""},
		Transactions:  []string{""},
		Company:       "",
		Name:          "",
		StartDate:     "",
		EndDate:       "",
		CreationDate:  "",
		Updated:       "",
		Payment:       "",
		PaymentType:   "",
		PaymentSplit:  "",
		IsPay:         false,
		IsSign:        false,
		CoverageType:  "",
		Voucher:       "",
		Channel:       "",
		Covenant:      "",
		TaxAmount:     0,
		PriceNett:     0,
		PriceGross:    0,
		Contractor:    &models.User{},
		DocumentName:  "",
		Statements:    []models.Statement{{}},
		Attachments:   []models.Attachment{{}},
		Assets:        []models.Asset{{}},
		Claim:         []models.Claim{{}},
	}
	log.Println(policy)
	b, e := json.Marshal(policy)
	lib.CheckError(e)
	ref, _ := lib.PutFirestore("policy", "", policy)
	log.Println(ref)
	return string(b), policy
}
func Emit(w http.ResponseWriter, r *http.Request) (string, interface{}) {

	return "", nil
}
func GetNumberCompany(w http.ResponseWriter, r *http.Request) (string, interface{}) {

	return "", nil
}
