package broker

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"cloud.google.com/go/firestore"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	lib "github.com/wopta/goworkspace/lib"
	models "github.com/wopta/goworkspace/models"
	"google.golang.org/api/iterator"
)

func init() {
	log.Println("INIT Broker")
	functions.HTTP("Broker", Broker)
}

func Broker(w http.ResponseWriter, r *http.Request) {

	log.Println("Broker")
	lib.EnableCors(&w, r)
	w.Header().Set("Access-Control-Allow-Methods", "POST")
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

func GetNumberCompany(w http.ResponseWriter, r *http.Request) (string, interface{}) {

	return "", nil
}

type BrokerResponse struct {
	EnvelopSignId string `json:"envelopSignId"`
	LinkGcs       string `json:"linkGcs"`
	Bytes         string `json:"bytes"`
}

func ToListData(query *firestore.DocumentIterator) []models.Policy {
	var result []models.Policy
	for {
		d, err := query.Next()
		log.Println("for")
		if err != nil {
			log.Println("error")
		}
		if err != nil {
			if err == iterator.Done {
				log.Println("iterator.Done")
				break
			}

		}

		var value models.Policy

		e := d.DataTo(&value)

		log.Println("todata")
		lib.CheckError(e)
		result = append(result, value)

		log.Println(len(result))
	}
	return result
}
func GetSequenceByProduct(name string) (string, int) {
	var companyDefault string
	switch name {
	case "global":
		companyDefault = "49999999"
	}
	var numberCompany string
	var number int
	log.Println("GetSequenceByProduct")
	rn, e := lib.OrderWhereLimitFirestoreErr("policy", "company", "numberCompany", "==", name, firestore.Desc, 1)
	lib.CheckError(e)
	log.Println("RN")
	policy := ToListData(rn)
	if len(policy) == 0 {
		log.Println("len(policy) == 0")
		numberCompany = companyDefault
	} else {
		log.Println("else")
		log.Println(rn)
		log.Println("policy use company")
		log.Println(len(policy))
		intNumberCompany, e := strconv.Atoi(policy[0].NumberCompany)
		lib.CheckError(e)
		numberCompany = fmt.Sprint(intNumberCompany + 1)
		number = policy[0].Number + 1
	}
	r, e := lib.OrderLimitFirestoreErr("policy", "number", firestore.Desc, 1)
	lib.CheckError(e)
	policyCompany := ToListData(r)
	if len(policyCompany) == 0 {
		log.Println("len(policy) == 0")
		number = 1
	} else {
		log.Println("policy use number")

		number = policyCompany[0].Number + 1
	}
	return numberCompany, number
}
func GetSequenceProposal(name string) int {
	var number int
	log.Println("GetSequenceProposal")
	r, e := lib.OrderLimitFirestoreErr("policy", "proposalNumber", firestore.Desc, 1)
	lib.CheckError(e)
	policyCompany := ToListData(r)
	if len(policyCompany) == 0 {
		log.Println("len(policy) == 0")
		number = 1
	} else {
		log.Println("policy use number")

		number = policyCompany[0].Number + 1
	}
	return number
}
