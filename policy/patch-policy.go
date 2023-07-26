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

func PatchPolicyFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err       error
		policy    models.Policy
		policyUID string
		input     map[string]interface{}
	)
	log.Println("PatchPolicyFx")

	firePolicy := lib.GetDatasetByEnv(r.Header.Get("origin"), "policy")
	policyUID = r.Header.Get("uid")

	b := lib.ErrorByte(io.ReadAll(r.Body))
	err = json.Unmarshal(b, &policy)
	if err != nil {
		log.Println("PatchPolicyFx: unable to unmarshal request body")
		return `{"uid":"` + policyUID + `", "success":false}`, `{"uid":"` + policyUID + `", "success":false}`, err
	}

	input = make(map[string]interface{}, 0)
	input["assets"] = policy.Assets
	input["contractor"] = policy.Contractor
	input["fundsOrigin"] = policy.FundsOrigin
	if policy.Surveys != nil {
		input["surveys"] = policy.Surveys
	}
	if policy.Statements != nil {
		input["statements"] = policy.Statements
	}
	input["updated"] = time.Now().UTC()

	lib.FireUpdate(firePolicy, policyUID, input)

	return `{"uid":"` + policyUID + `", "success":true}`, `{"uid":"` + policyUID + `", "success":true}`, err
}

func PatchPolicy(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err          error
		policyUID    string
		updateValues map[string]interface{}
	)
	log.Println("PatchPolicy")

	firePolicy := lib.GetDatasetByEnv(r.Header.Get("origin"), "policy")
	policyUID = r.Header.Get("uid")

	b := lib.ErrorByte(io.ReadAll(r.Body))
	err = json.Unmarshal(b, &updateValues)
	if err != nil {
		log.Println("PatchPolicy: unable to unmarshal request body")
		return `{"uid":"` + policyUID + `", "success":false}`, `{"uid":"` + policyUID + `", "success":false}`, nil
	}

	updateValues["updated"] = time.Now().UTC()

	err = lib.UpdateFirestoreErr(firePolicy, policyUID, updateValues)
	if err != nil {
		log.Println("PatchPolicy: error during update policy in firestore ")
		return `{"uid":"` + policyUID + `", "success":false}`, `{"uid":"` + policyUID + `", "success":false}`, nil
	}

	return `{"uid":"` + policyUID + `", "success":true}`, `{"uid":"` + policyUID + `", "success":true}`, err
}
