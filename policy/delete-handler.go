package policy

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

type PolicyDeleteReq struct {
	DeleteDesc string `json:"deleteDesc,omitempty"`
}

func DeletePolicyFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err       error
		policy    models.Policy
		policyUID string
		request   PolicyDeleteReq
	)
	log.Println("DeletePolicy")
	guaranteFire := lib.GetDatasetByEnv(r.Header.Get("origin"), "guarante")
	req := lib.ErrorByte(io.ReadAll(r.Body))
	json.Unmarshal(req, &request)
	firePolicy := lib.GetDatasetByEnv(r.Header.Get("origin"), "policy")
	policyUID = r.Header.Get("uid")
	docsnap := lib.GetFirestore(firePolicy, string(policyUID))
	docsnap.DataTo(&policy)
	policy.IsDeleted = true
	policy.DeleteDesc = request.DeleteDesc
	policy.Status = models.PolicyStatusDeleted
	policy.StatusHistory = append(policy.StatusHistory, models.PolicyStatusDeleted)
	lib.SetFirestore(firePolicy, policyUID, policy)
	policy.BigquerySave(r.Header.Get("origin"))
	models.SetGuaranteBigquery(policy, "delete", guaranteFire)
	return `{"uid":"` + policyUID + `", "success":true}`, `{"uid":"` + policyUID + `", "success":true}`, err
}
