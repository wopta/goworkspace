package reserved

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models"
	plc "gitlab.dev.wopta.it/goworkspace/policy"
)

type SetCoverageReservedResponse struct {
	Policy *models.Policy `json:"policy"`
}

func setCoverageReservedFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		response SetCoverageReservedResponse
		err      error
	)

	log.AddPrefix("SetCoverageReservedFx")
	defer log.PopPrefix()

	log.Println("Handler start -----------------------------------------------")

	policyUid := chi.URLParam(r, "policyUid")
	firePolicy := models.PolicyCollection

	log.Printf("getting policy %s from firestore...", policyUid)
	originalPolicy, err := plc.GetPolicy(policyUid)
	if err != nil {
		log.ErrorF("error unable to retrieve original policy: %s", err.Error())
		return "", nil, err
	}

	input, err := UpdatePolicyReservedCoverage(&originalPolicy)
	if err != nil {
		log.ErrorF("error calculating reserved coverage: %s", err.Error())
		return "", nil, err
	}

	_, err = lib.FireUpdate(firePolicy, policyUid, input)
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

	log.Println("Handler end -------------------------------------------------")

	return string(responseJson), response, err
}
