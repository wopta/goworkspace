package broker

import (
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	lib "github.com/wopta/goworkspace/lib"
	models "github.com/wopta/goworkspace/models"
)

func init() {
	log.Println("INIT Broker")
	functions.HTTP("Broker", Broker)
}

func Broker(w http.ResponseWriter, r *http.Request) {

	log.Println("Broker")
	lib.EnableCors(&w, r)
	route := lib.RouteData{

		Routes: []lib.Route{
			{
				Route:   "/v1/policies/fiscalcode/:fiscalcode",
				Handler: PolicyFiscalcode,
				Method:  "GET",
			},
			{
				Route:   "/v1/policy/:uid",
				Handler: GetPolicy,
				Method:  "GET",
			},

			{
				Route:   "/v1/policy/proposal",
				Handler: Proposal,
				Method:  "POST",
			},

			{
				Route:   "/v1/policy/emit",
				Handler: Emit,
				Method:  "POST",
			},
		},
	}
	route.Router(w, r)

}

func GetNumberCompany(w http.ResponseWriter, r *http.Request) (string, interface{}) {

	return "", nil
}

type BrokerResponse struct {
	EnvelopSignId string `json:"envelopSignId"`
	LinkGcs       string `json:"linkGcs"`
	Bytes         string `json:"bytes"`
}

func GetSequenceByCompany(name string) (string, int, int) {
	var (
		codeCompany         string
		companyDefault      int
		companyPrefix       string
		companyPrefixLenght string
		numberCompany       int
		number              int
	)
	switch name {
	case "global":
		companyDefault = 1
		companyPrefix = "WB"
		companyPrefixLenght = "%07d"
	case "axa":
		companyDefault = 1
		companyPrefix = "WB"
		companyPrefixLenght = "%07d"
	}

	rn, e := lib.OrderWhereLimitFirestoreErr("policy", "company", "numberCompany", "==", name, firestore.Desc, 1)
	lib.CheckError(e)
	policy := models.PolicyToListData(rn)

	if len(policy) == 0 {
		//WE0000001
		numberCompany = companyDefault
		codeCompany = companyPrefix + fmt.Sprintf(companyPrefixLenght, numberCompany)
		number = 1
	} else {
		numberCompany = policy[0].NumberCompany + 1
		codeCompany = companyPrefix + fmt.Sprintf(companyPrefixLenght, numberCompany)
		number = policy[0].Number + 1
	}
	r, e := lib.OrderLimitFirestoreErr("policy", "number", firestore.Desc, 1)
	lib.CheckError(e)
	policyCompany := models.PolicyToListData(r)
	if len(policyCompany) == 0 {
		number = 1
	} else {

		number = policyCompany[0].Number + 1
	}
	log.Println("GetSequenceByCompany: ", codeCompany)

	return codeCompany, numberCompany, number
}
func GetSequenceProposal(name string) int {
	var number int
	r, e := lib.OrderLimitFirestoreErr("policy", "proposalNumber", firestore.Desc, 1)
	lib.CheckError(e)
	policyCompany := models.PolicyToListData(r)
	if len(policyCompany) == 0 {
		number = 1
	} else {

		number = policyCompany[0].ProposalNumber + 1
	}
	log.Println("GetSequenceProposal: ", number)
	return number
}
