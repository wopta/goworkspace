package win

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/wopta/goworkspace/models"
)

type ApprovalReq struct {
	IdPratica int    `json:"idPratica"`
	Utente    string `json:"utente"`
}

func approvalCallback(policy models.Policy) (*http.Request, *http.Response, error) {
	log.Println("win wait for approval calback...")

	wp := policyDto(policy)
	payload := ApprovalReq{wp.IdPratica, wp.Utente}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, nil, err
	}

	client := &winClient{
		path: "/restba/extquote/richemissione",
	}
	return client.Post(bytes.NewReader(body))
}
