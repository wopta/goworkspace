package win

import (
	"bytes"
	"encoding/json"
	"log"

	"github.com/wopta/goworkspace/models"
)

type ApprovalReq struct {
	IdPratica int    `json:"idPratica"`
	Utente    string `json:"utente"`
}

func approvalCallback(policy models.Policy) error {
	log.Println("win wait for approval calback...")

	wp := policyDto(policy)
	payload := ApprovalReq{wp.IdPratica, wp.Utente}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	client := &winClient{
		path: "/restba/extquote/richemissione",
	}
	res, err := client.Post(bytes.NewReader(body))

	// TODO: should we do somethoing with the response?

	log.Println(res)

	return err
}