package policy

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

func PatchPolicyFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println("[PatchPolicyFx] Handler start -------------------------------")
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
		log.Printf("[PatchPolicyFx] error unable to unmarshal request body: %s", err.Error())
		response := fmt.Sprintf(responseTemplate, policyUid, false)
		return response, response, nil
	}

	originalPolicy, err := GetPolicy(policyUid, origin)
	if err != nil {
		log.Printf("[PatchPolicyFx] error unable to retrieve original policy: %s", err.Error())
		response := fmt.Sprintf(responseTemplate, policyUid, false)
		return response, response, nil
	}
	originalPolicyBytes, _ := json.Marshal(originalPolicy)
	log.Printf("[PatchPolicyFx] original policy: %s", string(originalPolicyBytes))

	input := PatchPolicy(&policy)
	log.Printf("[PatchPolicyFx] modified policy values: %v", input)

	_, err = lib.FireUpdate(firePolicy, policyUid, input)
	if err != nil {
		log.Printf("[PatchPolicyFx] error updating policy in firestore: %s", err.Error())
		response := fmt.Sprintf(responseTemplate, policyUid, false)
		return response, response, nil
	}

	response := fmt.Sprintf(responseTemplate, policyUid, true)

	return response, response, nil
}

func PatchPolicy(policy *models.Policy) map[string]interface{} {
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
	input["step"] = policy.Step
	input["updated"] = time.Now().UTC()

	isReserved, reservedInfo := reserved.GetReservedInfo(policy)
	input["isReserved"] = isReserved
	input["reservedInfo"] = reservedInfo

	return input
}
