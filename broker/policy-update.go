package broker

import (
	"encoding/json"
	"github.com/wopta/goworkspace/lib"
	"io"
	"log"
	"net/http"
	"time"
)

func UpdatePolicy(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
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
