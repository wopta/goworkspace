package broker

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func UpdatePolicy(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err       error
		policy    models.Policy
		policyUID string
		input     map[string]interface{}
	)
	log.Println("UpdatePolicy")

	firePolicy := lib.GetDatasetByEnv(r.Header.Get("origin"), "policy")
	policyUID = r.Header.Get("uid")

	b := lib.ErrorByte(io.ReadAll(r.Body))
	err = json.Unmarshal(b, &policy)
	if err != nil {
		return `{"uid":"` + policyUID + `", "success":"false"}`, `{"uid":"` + policyUID + `", "success":"false"}`, err
	}

	input = make(map[string]interface{}, 0)
	input["assets"] = policy.Assets
	input["contractor"] = policy.Contractor
	if policy.Surveys != nil {
		input["surveys"] = policy.Surveys
	}
	if policy.Statements != nil {
		input["statements"] = policy.Statements
	}
	input["updated"] = time.Now().UTC()

	lib.FireUpdate(firePolicy, policyUID, input)

	return `{"uid":"` + policyUID + `", "success":"true"}`, `{"uid":"` + policyUID + `", "success":"true"}`, err
}

func PatchPolicy(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err          error
		policyUID    string
		updateValues map[string]interface{}
	)
	log.Println("UpdatePolicy")

	firePolicy := lib.GetDatasetByEnv(r.Header.Get("origin"), "policy")
	policyUID = r.Header.Get("uid")

	b := lib.ErrorByte(io.ReadAll(r.Body))
	err = json.Unmarshal(b, &updateValues)
	if err != nil {
		return `{"uid":"` + policyUID + `", "success":"false"}`, `{"uid":"` + policyUID + `", "success":"false"}`, err
	}

	updateValues["updated"] = time.Now().UTC()

	err = lib.UpdateFirestoreErr(firePolicy, policyUID, updateValues)
	if err != nil {
		return `{"uid":"` + policyUID + `", "success":"false"}`, `{"uid":"` + policyUID + `", "success":"false"}`, err
	}

	return `{"uid":"` + policyUID + `", "success":"true"}`, `{"uid":"` + policyUID + `", "success":"true"}`, err
}
func DeletePolicy(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err       error
		policy    models.Policy
		policyUID string
		request   PolicyDeleteReq
	)
	log.Println("DeletePolicy")
	guaranteFire := lib.GetDatasetByEnv(r.Header.Get("origin"), "guarante")
	req := lib.ErrorByte(ioutil.ReadAll(r.Body))
	json.Unmarshal(req, &request)
	firePolicy := lib.GetDatasetByEnv(r.Header.Get("origin"), "policy")
	policyUID = r.Header.Get("uid")
	docsnap := lib.GetFirestore(firePolicy, string(policyUID))
	docsnap.DataTo(&policy)
	policy.IsDeleted = true
	policy.DeleteDesc = request.DeleteDesc
	lib.SetFirestore(firePolicy, policyUID, policy)
	policy.BigquerySave(r.Header.Get("origin"))
	models.SetGuaranteBigquery(policy, "delete", guaranteFire)
	return `{"uid":"` + policyUID + `", "success":"true"}`, `{"uid":"` + policyUID + `", "success":"true"}`, err
}

type PolicyDeleteReq struct {
	DeleteDesc string `json:"deleteDesc,omitempty"`
}
