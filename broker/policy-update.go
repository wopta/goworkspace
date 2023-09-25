package broker

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/reserved"
)

func UpdatePolicyFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println("[UpdatePolicyFx] Handler start ------------------------------")
	var (
		err              error
		responseTemplate = `{"uid":"%s","success":%t}`
		policy           models.Policy
	)

	origin := r.Header.Get("Origin")
	policyUid := r.Header.Get("uid")
	firePolicy := lib.GetDatasetByEnv(origin, models.PolicyCollection)

	body := lib.ErrorByte(io.ReadAll(r.Body))
	err = json.Unmarshal(body, &policy)
	if err != nil {
		log.Printf("[UpdatePolicyFx] error unable to unmarshal request body: %s", err.Error())
		response := fmt.Sprintf(responseTemplate, policyUid, false)
		return response, response, nil
	}

	originalPolicy, err := GetPolicy(policyUid, origin)
	if err != nil {
		log.Printf("[UpdatePolicyFx] error unable to retrieve original policy: %s", err.Error())
		response := fmt.Sprintf(responseTemplate, policyUid, false)
		return response, response, nil
	}
	originalPolicyBytes, _ := json.Marshal(originalPolicy)
	log.Printf("[UpdatePolicyFx] original policy: %s", string(originalPolicyBytes))

	input := UpdatePolicy(&policy)
	inputJson, err := json.Marshal(input)
	if err != nil {
		log.Printf("[UpdatePolicyFx] error unable to marshal input result: %s", err.Error())
		response := fmt.Sprintf(responseTemplate, policyUid, false)
		return response, response, nil
	}
	log.Printf("[UpdatePolicyFx] modified policy values: %v", string(inputJson))

	_, err = lib.FireUpdate(firePolicy, policyUid, input)
	if err != nil {
		log.Printf("[UpdatePolicyFx] error updating policy in firestore: %s", err.Error())
		response := fmt.Sprintf(responseTemplate, policyUid, false)
		return response, response, nil
	}

	response := fmt.Sprintf(responseTemplate, policyUid, true)

	return response, response, nil
}

func UpdatePolicy(policy *models.Policy) map[string]interface{} {
	input := make(map[string]interface{}, 0)

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

	isReserved, reservedInfo := reserved.GetReservedInfo(policy)
	input["isReserved"] = isReserved
	input["reservedInfo"] = reservedInfo

	return input
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
		log.Println("PatchPolicy: unable to unmarshal request body")
		return `{"uid":"` + policyUID + `", "success":false}`, `{"uid":"` + policyUID + `", "success":false}`, nil
	}

	updateValues["updated"] = time.Now().UTC()

	err = lib.UpdateFirestoreErr(firePolicy, policyUID, updateValues)
	if err != nil {
		log.Println("PatchPolicy: error during policy update in firestore")
		return `{"uid":"` + policyUID + `", "success":false}`, `{"uid":"` + policyUID + `", "success":false}`, nil
	}

	return `{"uid":"` + policyUID + `", "success":true}`, `{"uid":"` + policyUID + `", "success":true}`, err
}

func DeletePolicy(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err       error
		policy    models.Policy
		policyUID string
		request   PolicyDeleteReq
	)
	log.Println("DeletePolicy")
	policyUID = r.Header.Get("uid")
	guaranteFire := lib.GetDatasetByEnv(r.Header.Get("origin"), "guarante")
	req := lib.ErrorByte(io.ReadAll(r.Body))
	err = json.Unmarshal(req, &request)
	if err != nil {
		log.Printf("DeletePolicy: unable to delete policy %s", policyUID)
		return `{"uid":"` + policyUID + `", "success":false}`, `{"uid":"` + policyUID + `", "success":false}`, nil
	}
	firePolicy := lib.GetDatasetByEnv(r.Header.Get("origin"), "policy")
	docsnap := lib.GetFirestore(firePolicy, policyUID)
	docsnap.DataTo(&policy)
	if policy.IsDeleted || !policy.CompanyEmit {
		log.Printf("DeletePolicy: can't delete policy %s", policyUID)
		return `{"uid":"` + policyUID + `", "success":false}`, `{"uid":"` + policyUID + `", "success":false}`, nil

	}
	policy.IsDeleted = true
	policy.DeleteCode = request.Code
	policy.DeleteDesc = request.Description
	policy.DeleteDate = request.Date
	policy.RefundType = request.RefundType
	policy.Status = models.PolicyStatusDeleted
	policy.StatusHistory = append(policy.StatusHistory, models.PolicyStatusDeleted)
	lib.SetFirestore(firePolicy, policyUID, policy)
	policy.BigquerySave(r.Header.Get("origin"))
	models.SetGuaranteBigquery(policy, "delete", guaranteFire)
	return `{"uid":"` + policyUID + `", "success":true}`, `{"uid":"` + policyUID + `", "success":true}`, err
}

type PolicyDeleteReq struct {
	Code        string    `json:"code,omitempty"`
	Description string    `json:"description,omitempty"`
	Date        time.Time `json:"date"`
	RefundType  string    `json:"refundType,omitempty"`
}
