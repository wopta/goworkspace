package renew

import (
	"encoding/json"
	"github.com/wopta/goworkspace/lib"
	"io"
	"log"
	"net/http"
)

type RenewPolicyReq struct {
	PolicyUid string `json:"policyUid"`
}

func RenewPolicyFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err error
		req RenewPolicyReq
	)

	log.SetPrefix("[RenewPolicyFx] ")
	defer log.SetPrefix("")

	log.Println("Handler start -----------------------------------------------")

	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Printf("error unmarshalling body: %v", err)
		return "", nil, err
	}

	log.Println("Handler end -------------------------------------------------")

	return "", nil, err
}
