package win

import (
	"encoding/json"

	"github.com/wopta/goworkspace/models"
)

type emissioneReq struct {
	DtEmissione    string `json:"dtEmissione"`
	IdPratica      int    `json:"idPratica"`
	LuogoEmissione string `json:"luogoEmissione"`
	NumOriginali   int    `json:"numOriginali"`
	NumPol         string `json:"numPol"`
	NumPolSost     string `json:"numPolSost"`
	Ramo           string `json:"ramo"`
	PerAss         perAss `json:"perAss"`
	Utente         string `json:"utente"`
}

func emissione(policy models.Policy) ([]byte, error) {
	winPolicy := policyDto(policy)
	payload := emissioneReq{
		DtEmissione:    winPolicy.DtEmissione,
		IdPratica:      winPolicy.IdPratica,
		LuogoEmissione: winPolicy.LuogoEmissione,
		NumOriginali:   winPolicy.NumOriginali,
		NumPol:         winPolicy.NumPol,
		NumPolSost:     winPolicy.NumPolSost,
		PerAss:         winPolicy.PerAss,
		Ramo:           winPolicy.Ramo,
		Utente:         winPolicy.Utente,
	}

	return json.Marshal(payload)
}
