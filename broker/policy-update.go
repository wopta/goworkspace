package broker

import (
	"encoding/json"
	"github.com/wopta/goworkspace/lib"
	"io"
	"log"
	"net/http"
)

func UpdatePolicy(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err          error
		updateValues map[string]interface{}
	)
	log.Println("UpdatePolicy")

	firePolicy := lib.GetDatasetByEnv(r.Header.Get("origin"), "policy")

	b := lib.ErrorByte(io.ReadAll(r.Body))
	err = json.Unmarshal(b, &updateValues)
	lib.CheckError(err)

	err = lib.UpdateFirestoreErr(firePolicy, updateValues["uid"].(string), updateValues)

	return `{"uid":"` + updateValues["uid"].(string) + `"}`, `{"uid":"` + updateValues["uid"].(string) + `"}`, err
}
