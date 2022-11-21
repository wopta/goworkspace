package broker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	doc "github.com/wopta/goworkspace/document"
	lib "github.com/wopta/goworkspace/lib"
	mail "github.com/wopta/goworkspace/mail"
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
func Proposal(w http.ResponseWriter, r *http.Request) (string, interface{}) {
	log.Println("Proposal")
	var policy models.Policy
	req := lib.ErrorByte(ioutil.ReadAll(r.Body))

	e := json.Unmarshal([]byte(req), &policy)
	lib.CheckError(e)
	defer r.Body.Close()
	//policy, e := models.UnmarshalPolicy(req)
	policy.CreationDate = time.Now()
	policy.Updated = time.Now()
	policy.CreationDate = time.Now()
	policy.Status = models.Proposal
	log.Println("GetSequenceByProduct")
	company, numb := GetSequenceByProduct("global")
	log.Println(string(company))
	policy.NumberCompany = company
	policy.Number = numb
	log.Println("save")
	ref, _ := lib.PutFirestore("policy", policy)
	log.Println("saved")
	var obj mail.MailRequest
	obj.From = "noreply@wopta.it"
	obj.To = []string{policy.Contractor.Mail}
	obj.Message = `<p>ciao </p> `
	obj.Subject = "Wopta Proposta e set informantivo"
	obj.IsHtml = true
	mail.SendMail(obj)
	log.Println(ref.ID)

	return `{"uid":"` + ref.ID + `"}`, policy
}
func Emit(w http.ResponseWriter, r *http.Request) (string, interface{}) {
	var (
		result map[string]string
	)

	log.Println("PmiAllrisk")
	request := lib.ErrorByte(ioutil.ReadAll(r.Body))
	// Unmarshal or Decode the JSON to the interface.
	json.Unmarshal([]byte(request), &result)
	log.Println(result["uid"])
	var policy models.Policy
	docsnap := lib.GetFirestore("policy", string(result["uid"]))
	docsnap.DataTo(&policy)
	_, p := doc.ContractObj(policy)
	policy.DocumentName = p.(doc.DodumentResponse).LinkGcs
	doc.NamirialOtp(policy)
	return "", nil
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
			if err == iterator.Done {
				break
			}
			var value models.Policy
			e := d.DataTo(&value)
			log.Println("todata")
			lib.CheckError(e)
			result = append(result, value)

		}

	}
	return result
}
func GetSequenceByProduct(name string) (string, int) {
	var numberCompany string
	var number int
	log.Println("GetSequenceByProduct")
	rn, e := lib.OrderWhereLimitFirestoreErr("policy", "", "company", "==", name, firestore.Desc, 1)
	log.Println("RN")

	if e != nil {
		log.Println("e nil")
		numberCompany = "49999999"
	} else {
		log.Println("else")
		log.Println(rn)
		policy := ToListData(rn)
		intNumberCompany, e := strconv.Atoi(policy[1].NumberCompany)
		lib.CheckError(e)
		numberCompany = fmt.Sprint(intNumberCompany + 1)
		number = policy[1].Number + 1
	}
	r, e := lib.OrderLimitFirestoreErr("policy", "number", firestore.Desc, 1)
	if e != nil {
		number = 1
	} else {
		policy := ToListData(r)
		number = policy[1].Number + 1
	}
	return numberCompany, number
}
