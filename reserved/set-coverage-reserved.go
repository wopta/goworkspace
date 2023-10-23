package reserved

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	plc "github.com/wopta/goworkspace/policy"
)

type SetCoverageReservedResponse struct {
	Policy *models.Policy `json:"policy"`
}

func SetCoverageReservedFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		response SetCoverageReservedResponse
		err      error
	)

	log.Println("[SetCoverageReservedFx] Handler start -----------------------")

	origin := r.Header.Get("Origin")
	policyUid := r.Header.Get("policyUid")
	firePolicy := lib.GetDatasetByEnv(origin, models.PolicyCollection)

	log.Printf("[SetCoverageReservedFx] getting policy %s from firestore...", policyUid)
	originalPolicy, err := plc.GetPolicy(policyUid, origin)
	if err != nil {
		log.Printf("[SetCoverageReservedFx] error unable to retrieve original policy: %s", err.Error())
		return "", nil, err
	}

	input := UpdatePolicyReservedCoverage(&originalPolicy, origin)

	_, err = lib.FireUpdate(firePolicy, policyUid, input)
	if err != nil {
		log.Printf("[SetCoverageReservedFx] error updating policy in firestore: %s", err.Error())
		return "", nil, err
	}

	// TODO: improve me
	updatedPolicy, err := plc.GetPolicy(policyUid, origin)
	if err != nil {
		log.Printf("[SetCoverageReservedFx] error unable to retrieve updated policy: %s", err.Error())
		return "", nil, err
	}
	response.Policy = &updatedPolicy

	responseJson, err := json.Marshal(&response)

	log.Printf("[SetCoverageReservedFx] response: %s", string(responseJson))
	log.Println("[SetCoverageReservedFx] Handler end -------------------------")

	return string(responseJson), response, err
}
