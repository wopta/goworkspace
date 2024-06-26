package reserved

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

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

	log.SetPrefix("[SetCoverageReservedFx] ")
	defer log.SetPrefix("")

	log.Println("Handler start -----------------------------------------------")

	origin := r.Header.Get("Origin")
	policyUid := chi.URLParam(r, "policyUid")
	firePolicy := lib.GetDatasetByEnv(origin, models.PolicyCollection)

	log.Printf("getting policy %s from firestore...", policyUid)
	originalPolicy, err := plc.GetPolicy(policyUid, origin)
	if err != nil {
		log.Printf("error unable to retrieve original policy: %s", err.Error())
		return "", nil, err
	}

	input, err := UpdatePolicyReservedCoverage(&originalPolicy, origin)
	if err != nil {
		log.Printf("error calculating reserved coverage: %s", err.Error())
		return "", nil, err
	}

	_, err = lib.FireUpdate(firePolicy, policyUid, input)
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

	log.Println("Handler end -------------------------------------------------")

	return string(responseJson), response, err
}
