package win

import (
	"fmt"
	"time"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

type anagrafica struct {
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

type garanzia struct {
	Garanzia         string `json:"garanzia"`
	Imposte          int    `json:"imposte"`
	PremioImponibile int    `json:"premioImponibile"`
	SommaAssicurata  int    `json:"sommaAssicurata"`
}

type perAss struct {
	BaseAnno          string `json:"baseAnno"`          // "ANNO_SOLARE", "ANNO_COMMERCIALE"
	DataEffetto       string `json:"dataEffetto"`       // time.DateOnly
	DataPrimaScadenza string `json:"dataPrimaScadenza"` // time.DateOnly
	DataScadenza      string `json:"dataScadenza"`      // time.DateOnly
	DurataIniziale    int    `json:"durataIniziale"`
	Frazionamento     string `json:"frazionamento"` // "ANNUALE", "SEMESTRALE", "QUADRIMESTRALE", "TRIMESTRALE", "BIMESTRALE", "MENSILE", "Unica_soluzione"
}

type totaleGaranzia struct {
	Garanzia         string `json:"garanzia"`
	Imposte          int    `json:"imposte"`
	PremioImponibile int    `json:"premioImponibile"`
	Totale           int    `json:"totale"`
}

type totale struct {
	Imposte          int              `json:"imposte"`
	PremioImponibile int              `json:"premioImponibile"`
	Totale           int              `json:"totale"`
	TotaliGaranzie   []totaleGaranzia `json:"totaliGaranzie"`
}

type policy struct {
	Anagrafica     anagrafica `json:"anagrafica"`
	Garanzie       []garanzia `json:"garanzie"`
	IdPratica      int        `json:"idPratica"`
	PerAss         perAss     `json:"perAss"`
	Prodotto       string     `json:"prodotto"`
	TotaleAnnuo    totale     `json:"totaleAnnuo"`
	TotaleFirma    totale     `json:"totaleFirma"`
	TotaleFutura   totale     `json:"totaleFutura"`
	Utente         string     `json:"utente"`
	DtEmissione    string     `json:"dtEmissione"`
	LuogoEmissione string     `json:"luogoEmissione"`
	NumOriginali   int        `json:"numOriginali"`
	NumPol         string     `json:"numPol"`
	NumPolSost     string     `json:"numPolSost"`
	Ramo           string     `json:"ramo"`
}

var paymentSplitMap map[string]string = map[string]string{
	string(models.PaySplitYearly):  "ANNUALE",
	string(models.PaySplitMonthly): "MENSILE",
}

func policyDto(p models.Policy) policy {
	var (
		wp                  policy
		an                  anagrafica
		pa                  perAss
		totale              totale
		contractorBirthDate time.Time
		err                 error
	)

	// Map contractor data
	an.Cap = p.Contractor.PostalCode
	an.Cf = p.Contractor.FiscalCode
	an.CodRagSoc = "" // ?
	an.Cognome = p.Contractor.Surname
	an.Comune = p.Contractor.Residence.City
	if contractorBirthDate, err = time.Parse(time.RFC3339, p.Contractor.BirthDate); err != nil {
		return policy{}
	}
	an.DataNascita = contractorBirthDate.Format(time.DateOnly)
	an.Descrizione = ""
	an.Indirizzo = lib.TrimSpace(fmt.Sprintf("%s, %s", p.Contractor.Residence.StreetName, p.Contractor.Residence.StreetNumber))
	an.LuogoNascita = p.Contractor.BirthCity
	an.Nazione = "IT"
	an.NazioneNascita = ""
	an.Nome = p.Contractor.Name
	an.Pfg = "F"
	an.Piva = ""
	an.Provincia = p.Contractor.Residence.CityCode
	an.ProvinciaNascita = p.Contractor.BirthProvince
	an.Sesso = p.Contractor.Gender

	wp.Anagrafica = an

	// Map guarantees
	wp.Garanzie = make([]garanzia, 0)
	for _, guarantee := range p.Assets[0].Guarantees {
		var g garanzia
		g.Garanzia = guarantee.Slug
		g.Imposte = int(guarantee.Value.Tax)
		g.SommaAssicurata = int(guarantee.Value.SumInsuredLimitOfIndemnity)
		if p.PaymentSplit == string(models.PaySplitMonthly) {
			g.PremioImponibile = int(guarantee.Value.PremiumNetMonthly)
		} else {
			g.PremioImponibile = int(guarantee.Value.PremiumNetYearly)
		}
		wp.Garanzie = append(wp.Garanzie, g)
	}

	wp.IdPratica = p.NumberCompany

	pa.BaseAnno = "ANNO_SOLARE"
	pa.DataEffetto = p.StartDate.Format(time.DateOnly)
	if p.PaymentSplit == string(models.PaySplitMonthly) {
		pa.DataPrimaScadenza = lib.AddMonths(p.StartDate, 1).Format(time.DateOnly)
	} else {
		pa.DataPrimaScadenza = lib.AddMonths(p.StartDate, 12).Format(time.DateOnly)
	}
	pa.DataScadenza = p.EndDate.Format(time.DateOnly)
	pa.DurataIniziale = p.GetDurationInYears()
	pa.Frazionamento = paymentSplitMap[p.PaymentSplit]

	wp.PerAss = pa

	wp.Prodotto = p.Name

	// Map totals
	totale.Imposte = int(p.TaxAmount)
	totale.PremioImponibile = int(p.PriceNett)
	totale.Totale = int(p.PriceGross)
	totale.TotaliGaranzie = make([]totaleGaranzia, 0)
	for _, g := range wp.Garanzie {
		totale.TotaliGaranzie = append(totale.TotaliGaranzie, totaleGaranzia{
			Garanzia:         g.Garanzia,
			Imposte:          g.Imposte,
			PremioImponibile: g.PremioImponibile,
			Totale:           g.SommaAssicurata,
		})
	}
	wp.TotaleAnnuo = totale
	wp.TotaleFirma = totale
	wp.TotaleFutura = totale

	wp.Utente = lib.TrimSpace(fmt.Sprintf("%s %s", p.Contractor.Name, p.Contractor.Surname))

	wp.DtEmissione = p.EmitDate.Format(time.DateOnly)

	wp.LuogoEmissione = "Italia"

	wp.NumOriginali = 0

	wp.NumPol = p.CodeCompany

	wp.Ramo = p.Channel

	return wp
}
