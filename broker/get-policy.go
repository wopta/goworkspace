package broker

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func GetPolicyFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	firePolicy := lib.GetDatasetByEnv(r.Header.Get("origin"), "policy")
	log.Println("GetPolicy")
	log.Println(r.RequestURI)
	policyUid := r.Header.Get("uid")
	log.Println(policyUid)
	policy, _ := GetPolicy(policyUid, firePolicy)
	res, err := json.Marshal(policy)
	return string(res), policy, err
}

func GetPolicy(uid string, origin string) (models.Policy, error) {
	var (
		policy models.Policy
		err    error
	)
	firePolicy := lib.GetDatasetByEnv(origin, "policy")
	docsnap, err := lib.GetFirestoreErr(firePolicy, uid)
	if err != nil {
		return policy, err
	}
	err = docsnap.DataTo(&policy)
	return policy, err
}
