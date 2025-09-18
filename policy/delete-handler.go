package policy

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models"
)

type PolicyDeleteRequest struct {
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Code        string    `json:"code"`
}

func DeletePolicyFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err     error
		request PolicyDeleteRequest
	)

	log.AddPrefix("DeletePolicyFx")
	defer log.PopPrefix()

	log.Println("Handler start -----------------------------------------------")

	policyUid := chi.URLParam(r, "uid")
	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	err = json.Unmarshal(body, &request)
	if err != nil {
		log.ErrorF("error unmarshaling request: %s", err.Error())
		return "", nil, err
	}

	policy, err := GetPolicy(policyUid)
	if err != nil {
		return "", nil, err
	}

	deletePolicy(&policy, request)

	firePolicy := lib.PolicyCollection

	log.Println("setting policy to delete in firestore...")
	err = lib.SetFirestoreErr(firePolicy, policyUid, policy)
	if err != nil {
		log.ErrorF("error saving policy to firestore: %s", err.Error())
		return "", nil, err
	}
	log.Println("policy set to deleted in firestore")

	log.Println("setting policy to delete in bigquery...")
	policy.BigquerySave()

	guaranteFire := lib.GuaranteeCollection
	log.Println("updating policy's guarantees to delete in bigquery...")
	models.SetGuaranteBigquery(policy, "delete", guaranteFire)

	log.Println("Handler end -------------------------------------------------")
	policy.AddSystemNote(models.GetDeletePolicyNote)
	return "{}", nil, nil
}

func deletePolicy(p *models.Policy, request PolicyDeleteRequest) {
	p.IsDeleted = true
	p.DeleteDesc = request.Description
	p.DeleteDate = request.Date
	p.DeleteCode = request.Code
	p.Status = models.PolicyStatusDeleted
	p.StatusHistory = append(p.StatusHistory, p.Status)
	p.Updated = time.Now().UTC()
}
