package broker

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	plc "github.com/wopta/goworkspace/policy"
	"github.com/wopta/goworkspace/reserved"
)

type UpdatePolicyResponse struct {
	Policy *models.Policy `json:"policy"`
}

func UpdatePolicyFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err         error
		inputPolicy models.Policy
		response    UpdatePolicyResponse
	)

	log.SetPrefix("[UpdatePolicyFx]")
	defer log.SetPrefix("")

	log.Println("Handler start ------------------------------")

	origin := r.Header.Get("Origin")
	policyUid := r.Header.Get("uid")
	firePolicy := lib.GetDatasetByEnv(origin, models.PolicyCollection)

	body := lib.ErrorByte(io.ReadAll(r.Body))
	err = json.Unmarshal(body, &inputPolicy)
	if err != nil {
		log.Printf("error unable to unmarshal request body: %s", err.Error())
		return "", nil, err
	}

	log.Printf("request body: %s", string(body))

	inputPolicy.Normalize()

	originalPolicy, err := plc.GetPolicy(policyUid, origin)
	if err != nil {
		log.Printf("error unable to retrieve original policy: %s", err.Error())
		return "", nil, err
	}
	originalPolicyBytes, _ := json.Marshal(originalPolicy)
	log.Printf("original policy: %s", string(originalPolicyBytes))

	mergedInput := make(map[string]interface{})
	input := plc.UpdatePolicy(&inputPolicy)
	for k, v := range input {
		mergedInput[k] = v
	}
	inputReserved := reserved.UpdatePolicyReserved(&inputPolicy)
	for k, v := range inputReserved {
		mergedInput[k] = v
	}

	inputJson, err := json.Marshal(mergedInput)
	if err != nil {
		log.Printf("error unable to marshal input result: %s", err.Error())
		return "", nil, err
	}
	log.Printf("modified policy values: %v", string(inputJson))

	_, err = lib.FireUpdate(firePolicy, policyUid, mergedInput)
	if err != nil {
		log.Printf("error updating policy in firestore: %s", err.Error())
		return "", nil, err
	}

	// TODO: improve me
	updatedPolicy, err := plc.GetPolicy(policyUid, origin)
	if err != nil {
		log.Printf("error unable to retrieve updated policy: %s", err.Error())
		return "", nil, err
	}
	response.Policy = &updatedPolicy
	responseJson, err := json.Marshal(&response)

	updatedPolicy.BigquerySave(origin)

	log.Printf("response: %s", string(responseJson))

	return string(responseJson), response, err
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
		return "", nil, err
	}

	updateValues["updated"] = time.Now().UTC()

	err = lib.UpdateFirestoreErr(firePolicy, policyUID, updateValues)
	if err != nil {
		log.Println("PatchPolicy: error during policy update in firestore")
		return "", nil, err
	}

	return `{"uid":"` + policyUID + `"}`, `{"uid":"` + policyUID + `"}`, err
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
		return "", nil, err
	}
	firePolicy := lib.GetDatasetByEnv(r.Header.Get("origin"), "policy")
	docsnap := lib.GetFirestore(firePolicy, policyUID)
	docsnap.DataTo(&policy)
	if policy.IsDeleted || !policy.CompanyEmit {
		log.Printf("DeletePolicy: can't delete policy %s", policyUID)
		return "", nil, err

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
	return `{"uid":"` + policyUID + `"}`, `{"uid":"` + policyUID + `"}`, err
}

type PolicyDeleteReq struct {
	Code        string    `json:"code,omitempty"`
	Description string    `json:"description,omitempty"`
	Date        time.Time `json:"date"`
	RefundType  string    `json:"refundType,omitempty"`
}
