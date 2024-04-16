package broker

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/wopta/goworkspace/lib"
	plc "github.com/wopta/goworkspace/policy"
)

func GetPolicyFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.SetPrefix("[GetPolicyFx] ")
	defer log.SetPrefix("")

	log.Println("Handler start -----------------------------------------------")

	policyUid := chi.URLParam(r, "uid")
	log.Printf("policyUid: %s", policyUid)

	firePolicy := lib.GetDatasetByEnv(r.Header.Get("Origin"), lib.PolicyCollection)

	policy, _ := plc.GetPolicy(policyUid, firePolicy)
	res, err := json.Marshal(policy)

	log.Println("Handler end -------------------------------------------------")

	return string(res), policy, err
}
