package win

import (
	"encoding/json"

	"github.com/wopta/goworkspace/models"
)

type inspraticaReq struct {
	Anagrafica   anagrafica `json:"anagrafica"`
	Garanzie     []garanzia `json:"garanzie"`
	IdPratica    int        `json:"idPratica"`
	PerAss       perAss     `json:"perAss"`
	Prodotto     string     `json:"prodotto"`
	TotaleAnnuo  totale     `json:"totaleAnnuo"`
	TotaleFirma  totale     `json:"totaleFirma"`
	TotaleFutura totale     `json:"totaleFutura"`
	Utente       string     `json:"utente"`
	StatoPratica string     `json:"statoPratica"` // "QUOTAZIONE_ACCETTATA", "RICHIESTA_QUOTAZIONE"
}

func inspratica(policy models.Policy, state string) ([]byte, error) {
	winPolicy := policyDto(policy)
	payload := inspraticaReq{
		Anagrafica:   winPolicy.Anagrafica,
		Garanzie:     winPolicy.Garanzie,
		IdPratica:    winPolicy.IdPratica,
		PerAss:       winPolicy.PerAss,
		Prodotto:     winPolicy.Prodotto,
		TotaleAnnuo:  winPolicy.TotaleAnnuo,
		TotaleFirma:  winPolicy.TotaleFirma,
		TotaleFutura: winPolicy.TotaleFutura,
		Utente:       winPolicy.Utente,
		StatoPratica: state,
	}

	return json.Marshal(payload)
}
