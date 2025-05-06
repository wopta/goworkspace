package policy

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/lib/log"
	"github.com/wopta/goworkspace/models"
)

func GetPolicyFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		result map[string]string
	)

	log.AddPrefix("GetPolicyFx")
	defer log.PopPrefix()

	log.Println("Handler start -----------------------------------------------")

	firePolicy := lib.GetDatasetByEnv(r.Header.Get("Origin"), lib.PolicyCollection)
	uid := chi.URLParam(r, "uid")
	log.Println(uid)

	request := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	json.Unmarshal(request, &result)

	policy, _ := GetPolicy(uid, firePolicy)
	res, _ := json.Marshal(policy)

	log.Println("Handler end -------------------------------------------------")

	return string(res), policy, nil
}

func GetPolicy(uid string, origin string) (models.Policy, error) {
	var (
		policy models.Policy
		err    error
	)
	firePolicy := lib.GetDatasetByEnv(origin, lib.PolicyCollection)
	docsnap, err := lib.GetFirestoreErr(firePolicy, uid)
	if err != nil {
		return policy, err
	}
	err = docsnap.DataTo(&policy)
	return policy, err
}

// TODO: keep only one: GetPolicy or GetPolicyByUid
func GetPolicyByUid(policyUid string, origin string) models.Policy {
	firePolicy := lib.GetDatasetByEnv(origin, "policy")
	policyF := lib.GetFirestore(firePolicy, policyUid)
	var policy models.Policy
	policyF.DataTo(&policy)
	policyM, _ := policy.Marshal()
	log.Println("GetPolicyByUid: Policy "+policyUid+" found: ", string(policyM))

	return policy
}
