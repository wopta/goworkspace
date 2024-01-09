package policy

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

type PolicyDeleteRequest struct {
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Code        string    `json:"code"`
}

func DeletePolicyFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.SetPrefix("[DeletePolicyFx] ")

	var (
		err     error
		request PolicyDeleteRequest
	)

	log.Println("Handler start -----------------------------------------------")

	origin := r.Header.Get("Origin")
	policyUid := r.Header.Get("uid")
	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()
	err = json.Unmarshal(body, &request)
	if err != nil {
		log.Printf("error unmarshaling request: %s", err.Error())
		return "", nil, err
	}

	policy := GetPolicyByUid(policyUid, origin)

	deletePolicy(&policy, request)

	firePolicy := lib.GetDatasetByEnv(origin, models.PolicyCollection)

	log.Println("setting policy to delete in firestore...")
	err = lib.SetFirestoreErr(firePolicy, policyUid, policy)
	if err != nil {
		log.Printf("error saving policy to firestore: %s", err.Error())
		return "", nil, err
	}
	log.Println("policy set to deleted in firestore")

	log.Println("setting policy to delete in bigquery...")
	policy.BigquerySave(origin)

	guaranteFire := lib.GetDatasetByEnv(origin, models.GuaranteeCollection)
	log.Println("updating policy's guarantees to delete in bigquery...")
	models.SetGuaranteBigquery(policy, "delete", guaranteFire)

	log.Println("Handler end -------------------------------------------------")
	return "{}", nil, nil
}

func deletePolicy(p *models.Policy, request PolicyDeleteRequest) {
	p.IsDeleted = true
	p.DeleteDesc = request.Description
	p.DeleteDate = request.Date
	p.DeleteCode = request.Code
	p.Status = models.PolicyStatusDeleted
	p.StatusHistory = append(p.StatusHistory, p.Status)
}
