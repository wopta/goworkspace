package broker

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	doc "github.com/wopta/goworkspace/document"
	lib "github.com/wopta/goworkspace/lib"
	models "github.com/wopta/goworkspace/models"
)

func Emit(w http.ResponseWriter, r *http.Request) (string, interface{}) {
	var (
		result map[string]string
	)

	log.Println("Emit")
	request := lib.ErrorByte(ioutil.ReadAll(r.Body))
	json.Unmarshal([]byte(request), &result)
	log.Println(result["uid"])
	var policy models.Policy
	docsnap := lib.GetFirestore("policy", string(result["uid"]))
	docsnap.DataTo(&policy)
	company, numb := GetSequenceByProduct("global")
	policy.NumberCompany = company
	policy.Number = numb
	policy.Updated = time.Now()
	p := <-doc.ContractObj(policy)
	log.Println(p.LinkGcs)
	policy.DocumentName = p.LinkGcs
	_, res := doc.NamirialOtp(policy)
	policy.IdSign = res.EnvelopeId
	lib.SetFirestore("policy", result["uid"], policy)
	b, e := json.Marshal(res)
	lib.CheckError(e)
	return string(b), res
}
