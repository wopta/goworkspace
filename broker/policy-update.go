package broker

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"gitlab.dev.wopta.it/goworkspace/lib/log"

	"github.com/go-chi/chi/v5"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
	plc "gitlab.dev.wopta.it/goworkspace/policy"
	"gitlab.dev.wopta.it/goworkspace/reserved"
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

	log.AddPrefix("UpdatePolicyFx")
	defer log.PopPrefix()

	log.Println("Handler start -----------------------------------------------")

	policyUid := chi.URLParam(r, "uid")
	firePolicy := lib.PolicyCollection

	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()
	err = json.Unmarshal(body, &inputPolicy)
	if err != nil {
		log.ErrorF("error unable to unmarshal request body: %s", err.Error())
		return "", nil, err
	}

	inputPolicy.Normalize()

	originalPolicy, err := plc.GetPolicy(policyUid)
	if err != nil {
		log.ErrorF("error unable to retrieve original policy: %s", err.Error())
		return "", nil, err
	}
	originalPolicyBytes, _ := json.Marshal(originalPolicy)
	log.Printf("original policy: %s", string(originalPolicyBytes))

	mergedInput := make(map[string]interface{})
	input, err := plc.UpdatePolicy(&inputPolicy)
	if err != nil {
		log.Printf("unable to update policy: %s", err.Error())
		return "", nil, err
	}
	for k, v := range input {
		mergedInput[k] = v
	}
	inputReserved := reserved.UpdatePolicyReserved(&inputPolicy)
	for k, v := range inputReserved {
		mergedInput[k] = v
	}

	inputJson, err := json.Marshal(mergedInput)
	if err != nil {
		log.ErrorF("error unable to marshal input result: %s", err.Error())
		return "", nil, err
	}
	log.Printf("modified policy values: %v", string(inputJson))

	_, err = lib.FireUpdate(firePolicy, policyUid, mergedInput)
	if err != nil {
		log.ErrorF("error updating policy in firestore: %s", err.Error())
		return "", nil, err
	}

	// TODO: improve me
	updatedPolicy, err := plc.GetPolicy(policyUid)
	if err != nil {
		log.ErrorF("error unable to retrieve updated policy: %s", err.Error())
		return "", nil, err
	}
	response.Policy = &updatedPolicy
	responseJson, err := json.Marshal(&response)

	updatedPolicy.BigquerySave()

	log.Println("Handler end -------------------------------------------------")

	return string(responseJson), response, err
}

func PatchPolicyFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err          error
		policyUID    string
		updateValues map[string]interface{}
	)

	log.AddPrefix("PatchPolicyFx")
	defer log.PopPrefix()

	log.Println("Handler start -----------------------------------------------")

	firePolicy := lib.PolicyCollection
	policyUID = chi.URLParam(r, "uid")

	b := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	err = json.Unmarshal(b, &updateValues)
	if err != nil {
		log.Println("unable to unmarshal request body")
		return "", nil, err
	}

	updateValues["updated"] = time.Now().UTC()

	err = lib.UpdateFirestoreErr(firePolicy, policyUID, updateValues)
	if err != nil {
		log.ErrorF("error during policy update in firestore")
		return "", nil, err
	}

	log.Println("Handler end -------------------------------------------------")

	return `{"uid":"` + policyUID + `"}`, `{"uid":"` + policyUID + `"}`, err
}

func DeletePolicyFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err       error
		policy    models.Policy
		policyUID string
		request   PolicyDeleteReq
	)

	log.AddPrefix("DeletePolicyFx")
	defer log.PopPrefix()

	log.Println("Handler start -----------------------------------------------")

	policyUID = chi.URLParam(r, "uid")
	guaranteFire := lib.GuaranteeCollection
	req := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	err = json.Unmarshal(req, &request)
	if err != nil {
		log.Printf("DeletePolicy: unable to delete policy %s", policyUID)
		return "", nil, err
	}
	firePolicy := lib.PolicyCollection
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
	policy.BigquerySave()
	models.SetGuaranteBigquery(policy, "delete", guaranteFire)

	log.Println("Handler end -------------------------------------------------")

	return `{"uid":"` + policyUID + `"}`, `{"uid":"` + policyUID + `"}`, err
}

type PolicyDeleteReq struct {
	Code        string    `json:"code,omitempty"`
	Description string    `json:"description,omitempty"`
	Date        time.Time `json:"date"`
	RefundType  string    `json:"refundType,omitempty"`
}
