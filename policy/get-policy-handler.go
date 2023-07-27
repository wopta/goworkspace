package policy

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	lib "github.com/wopta/goworkspace/lib"
	models "github.com/wopta/goworkspace/models"
)

func GetPolicyFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		result map[string]string
	)
	firePolicy := lib.GetDatasetByEnv(r.Header.Get("origin"), "policy")
	log.Println("GetPolicy")
	log.Println(r.RequestURI)
	log.Println(r.Header.Get("uid"))
	requestPath := strings.Split(r.RequestURI, "/")
	request := lib.ErrorByte(ioutil.ReadAll(r.Body))
	json.Unmarshal([]byte(request), &result)
	log.Println(requestPath[2])
	policy, _ := GetPolicy(r.Header.Get("uid"), firePolicy)
	res, _ := json.Marshal(policy)
	return string(res), policy, nil
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
