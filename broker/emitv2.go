package broker

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func Emitv2Fx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println("[EmitFx] Handler start ----------------------------------------")

	var (
		result     EmitRequest
		e          error
		firePolicy string
		policy     models.Policy
	)

	origin := r.Header.Get("origin")
	firePolicy = lib.GetDatasetByEnv(origin, "policy")
	request := lib.ErrorByte(io.ReadAll(r.Body))

	log.Printf("[EmitFx] Request: %s", string(request))
	json.Unmarshal([]byte(request), &result)

	uid := result.Uid
	log.Printf("[EmitFx] Uid: %s", uid)

	docsnap := lib.GetFirestore(firePolicy, string(uid))
	docsnap.DataTo(&policy)
	policyJsonLog, _ := policy.Marshal()
	log.Printf("[EmitFx] Policy %s JSON: %s", uid, string(policyJsonLog))

	responseEmit := EmitV2(&policy, result, origin)
	b, e := json.Marshal(responseEmit)
	log.Println("[EmitFx] Response: ", string(b))

	return string(b), responseEmit, e
}

func EmitV2(policy *models.Policy, request EmitRequest, origin string) EmitResponse {
	var responseEmit EmitResponse

	firePolicy := lib.GetDatasetByEnv(origin, "policy")
	guaranteFire := lib.GetDatasetByEnv(origin, "guarante")
	policy.Uid = request.Uid // we should enforce the setting of the ID on proposal

	if policy.IsReserved && policy.Status != models.PolicyStatusWaitForApproval {
		emitApproval(policy)
	} else {
		log.Printf("[Emit] Policy Uid %s", request.Uid)

		emitBase(policy, origin)

		emitSign(policy, request, origin)

		emitPay(policy, request, origin)

		responseEmit = EmitResponse{UrlPay: policy.PayUrl, UrlSign: policy.SignUrl}
		policyJson, _ := policy.Marshal()
		log.Printf("[Emit] Policy %s: %s", request.Uid, string(policyJson))
	}

	policy.Updated = time.Now().UTC()
	lib.SetFirestore(firePolicy, request.Uid, policy)
	policy.BigquerySave(origin)
	models.SetGuaranteBigquery(*policy, "emit", guaranteFire)

	return responseEmit
}
