package win

import (
	"time"

	"github.com/wopta/goworkspace/models"
)

type WinAnagrafica struct {
	Cap              string `json:"cap"`              // min: 0, max: 5
	Cf               string `json:"cf"`               // min:0, max: 16
	CodRagSoc        string `json:"codRagSoc"`        // "CONDOMINIO", "COOPERATIVA", "DITTA INDIVIDUALE", "IMPRESA FAMIGLIARE", "ONLUS", "SAPA", "SAS", "SNC", "SPA", "SRL"
	Cognome          string `json:"cognome"`          // min: 0, max: 40
	Comune           string `json:"comune"`           // min: 0, max: 50
	DataNascita      string `json:"dataNascita"`      // time.DateOnly
	Descrizione      string `json:"descrizione"`      // min: 0, max: 80
	Indirizzo        string `json:"indirizzo"`        // min: 0, max: 80
	LuogoNascita     string `json:"luogoNascita"`     // min: 0, max: 50
	Nazione          string `json:"nazione"`          // min: 0, max: 2
	NazioneNascita   string `json:"nazioneNascita"`   // min: 0, max: 2
	Nome             string `json:"nome"`             // min: 0, max: 40
	Pfg              string `json:"pfg"`              // "F", "G"
	Piva             string `json:"piva"`             // min: 0, max: 11
	Provincia        string `json:"provincia"`        // min: 0, max: 3
	ProvinciaNascita string `json:"provinciaNascita"` // min: 0, max: 3
	Sesso            string `json:"sesso"`            // min: 0, max: 1
}

type WinGarancy struct {
	Garanzia         string `json:"garanzia"`
	Imposte          int    `json:"imposte"`
	PremioImponibile int    `json:"premioImponibile"`
	SommaAssicurare  int    `json:"sommaAssicurata"`
}

type WinPerAss struct {
	BaseAnno          string `json:"baseAnno"`          // "ANNO_SOLARE", "ANNO_COMMERCIALE"
	DataEffetto       string `json:"dataEffetto"`       // time.DateOnly
	DataPrimaScadenza string `json:"dataPrimaScadenza"` // time.DateOnly
	DataScadenza      string `json:"dataScadenza"`      // time.DateOnly
	DurataIniziale    int    `json:"durataIniziale"`
	Frazionamento     string `json:"frazionamento"` // "ANNUALE", "SEMESTRALE", "QUADRIMESTRALE", "TRIMESTRALE", "BIMESTRALE", "MENSILE", "Unica_soluzione"
}

type WinTotalGarantee struct {
	Garanzia         string `json:"garanzia"`
	Imposte          int    `json:"imposte"`
	PremioImponibile int    `json:"premioImponibile"`
	Totale           int    `json:"totale"`
}

type WinTotal struct {
	Imposte          int                `json:"imposte"`
	PremioImponibile int                `json:"premioImponibile"`
	Totale           int                `json:"totale"`
	TotaliGaranzie   []WinTotalGarantee `json:"totaliGaranzie"`
}

type WinPolicy struct {
	Anagrafica     WinAnagrafica `json:"anagrafica"`
	Garanzie       []WinGarancy  `json:"garanzie"`
	IdPratica      int           `json:"idPratica"`
	PerAss         WinPerAss     `json:"perAss"`
	Prodotto       string        `json:"prodotto"`
	TotaleAnnuo    WinTotal      `json:"totaleAnnuo"`
	TotaleFirma    WinTotal      `json:"totaleFirma"`
	TotaleFutura   WinTotal      `json:"totaleFutura"`
	Utente         string        `json:"utente"`
	DtEmissione    string        `json:"dtEmissione"`
	LuogoEmissione string        `json:"luogoEmissione"`
	NumOriginali   int           `json:"numOriginali"`
	NumPol         string        `json:"numPol"`
	NumPolSost     string        `json:"numPolSost"`
	Ramo           string        `json:"ramo"`
}

func policyDto(policy models.Policy) WinPolicy {
	var wp WinPolicy

	wp.DtEmissione = policy.EmitDate.Format(time.DateOnly)
	wp.IdPratica = policy.NumberCompany
	wp.LuogoEmissione = "Italia"
	wp.NumOriginali = 0
	wp.NumPol = policy.CodeCompany
	wp.Ramo = policy.Channel
	wp.Utente = policy.Contractor.Name + " " + policy.Contractor.Surname
	wp.PerAss.BaseAnno = "ANNO_SOLARE" // map paymentSplit
	wp.PerAss.DataEffetto = policy.StartDate.Format(time.DateOnly)
	wp.PerAss.DataPrimaScadenza = policy.StartDate.AddDate(1, 0, 0).Format(time.DateOnly) // base on paymentsplit
	wp.PerAss.DataScadenza = policy.EndDate.Format(time.DateOnly)
	wp.PerAss.Frazionamento = policy.PaymentSplit // map

	return wp
}
