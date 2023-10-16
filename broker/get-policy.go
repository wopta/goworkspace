package broker

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	plc "github.com/wopta/goworkspace/policy"
)

func GetPolicyFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	firePolicy := lib.GetDatasetByEnv(r.Header.Get("origin"), models.PolicyCollection)
	log.Println("GetPolicy")
	log.Println(r.RequestURI)
	policyUid := r.Header.Get("uid")
	log.Println(policyUid)
	policy, _ := plc.GetPolicy(policyUid, firePolicy)
	res, err := json.Marshal(policy)
	return string(res), policy, err
}
