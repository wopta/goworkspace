package renew

import (
	"encoding/json"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"io"
	"log"
	"net/http"
)

type RenewReport struct {
	Policy       models.Policy        `json:"policy"`
	Transactions []models.Transaction `json:"transactions"`
	Error        string               `json:"error,omitempty"`
}

type RenewPolicyReq struct {
	PolicyUid string `json:"policyUid"`
}

type RenewPolicyResp struct {
	Success []RenewReport `json:"success"`
	Failure []RenewReport `json:"failure"`
}

func RenewPolicyFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err  error
		req  RenewPolicyReq
		resp = RenewPolicyResp{
			Success: make([]RenewReport, 0),
			Failure: make([]RenewReport, 0),
		}
	)

	log.SetPrefix("[RenewPolicyFx] ")
	go func() {
		defer log.SetPrefix("")
		defer log.Println("Handler end -------------------------------------------------")
	}()

	log.Println("Handler start -----------------------------------------------")

	body := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Printf("error unmarshalling body: %v", err)
		return "", nil, err
	}

	rawResp, err := json.Marshal(resp)

	return string(rawResp), resp, err
}
