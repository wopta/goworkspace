package win

import (
	"bytes"
	"encoding/json"
	"log"
	"time"

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

func emitCallback(policy models.Policy) error {
	log.Println("win emit calback...")

	var payload EmitReq
	payload.DtEmissione = policy.EmitDate.Format(time.DateOnly)
	payload.IdPratica = policy.NumberCompany
	payload.LuogoEmissione = "Italia"
	payload.NumOriginali = 0
	payload.NumPol = policy.CodeCompany
	payload.Ramo = policy.Channel
	payload.Utente = policy.Contractor.Name + " " + policy.Contractor.Surname
	payload.PerAss.BaseAnno = "ANNO_SOLARE" // map paymentSplit
	payload.PerAss.DataEffetto = policy.StartDate.Format(time.DateOnly)
	payload.PerAss.DataPrimaScadenza = policy.StartDate.AddDate(1, 0, 0).Format(time.DateOnly) // base on paymentsplit
	payload.PerAss.DataScadenza = policy.EndDate.Format(time.DateOnly)
	payload.PerAss.Frazionamento = policy.PaymentSplit // map

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	client := &winClient{
		path: "/restba/extquote/emissione",
	}
	res, err := client.Post(bytes.NewReader(body))

	// TODO: should we do somethoing with the response?

	log.Println(res)

	return err
}
