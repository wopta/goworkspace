package win

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/wopta/goworkspace/models"
)

type EmitReq struct {
	DtEmissione    string `json:"dtEmissione"`
	IdPratica      int    `json:"idPratica"`
	LuogoEmissione string `json:"luogoEmissione"`
	NumOriginali   int    `json:"numOriginali"`
	NumPol         string `json:"numPol"`
	NumPolSost     string `json:"numPolSost"`
	PerAss         struct {
		BaseAnno          string `json:"baseAnno"`
		DataEffetto       string `json:"dataEffetto"`
		DataPrimaScadenza string `json:"dataPrimaScadenza"`
		DataScadenza      string `json:"dataScadenza"`
		DurataIniziale    int    `json:"durataIniziale"`
		Frazionamento     string `json:"frazionamento"`
	} `json:"perAss"`
	Ramo   string `json:"ramo"`
	Utente string `json:"utente"`
}

func emitCallback(policy models.Policy) (*http.Request, *http.Response, error) {
	log.Println("win emit calback...")

	winPolicy := policyDto(policy)
	payload := EmitReq(winPolicy)

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, nil, err
	}

	client := &winClient{
		path: "/restba/extquote/emissione",
	}
	return client.Post(bytes.NewReader(body))
}
