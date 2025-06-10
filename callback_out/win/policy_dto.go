package win

import (
	"fmt"
	"slices"
	"time"

	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/models"
)

type anagrafica struct {
	Cap              string `json:"cap"`                   // min: 0, max: 5
	Cf               string `json:"cf"`                    // min:0, max: 16
	CodRagSoc        string `json:"codRagSoc,omitempty"`   // "CONDOMINIO", "COOPERATIVA", "DITTA INDIVIDUALE", "IMPRESA FAMIGLIARE", "ONLUS", "SAPA", "SAS", "SNC", "SPA", "SRL"
	Cognome          string `json:"cognome"`               // min: 0, max: 40
	Comune           string `json:"comune"`                // min: 0, max: 50
	DataNascita      string `json:"dataNascita"`           // time.DateOnly
	Descrizione      string `json:"descrizione,omitempty"` // min: 0, max: 80
	Indirizzo        string `json:"indirizzo"`             // min: 0, max: 80
	LuogoNascita     string `json:"luogoNascita"`          // min: 0, max: 50
	Nazione          string `json:"nazione"`               // min: 0, max: 2
	NazioneNascita   string `json:"nazioneNascita"`        // min: 0, max: 2
	Nome             string `json:"nome"`                  // min: 0, max: 40
	Pfg              string `json:"pfg"`                   // "F", "G"
	Piva             string `json:"piva,omitempty"`        // min: 0, max: 11
	Provincia        string `json:"provincia"`             // min: 0, max: 3
	ProvinciaNascita string `json:"provinciaNascita"`      // min: 0, max: 3
	Sesso            string `json:"sesso"`                 // min: 0, max: 1
}

type garanzia struct {
	Garanzia         string  `json:"garanzia"`
	Imposte          float64 `json:"imposte"`
	PremioImponibile float64 `json:"premioImponibile"`
	SommaAssicurata  float64 `json:"sommaAssicurare"`
}

type perAss struct {
	BaseAnno          string `json:"baseAnno"`          // "ANNO_SOLARE", "ANNO_COMMERCIALE"
	DataEffetto       string `json:"dataEffetto"`       // time.DateOnly
	DataPrimaScadenza string `json:"dataPrimaScadenza"` // time.DateOnly
	DataScadenza      string `json:"dataScadenza"`      // time.DateOnly
	DurataIniziale    int    `json:"durataIniziale"`    // DataPrimaScadenza - DataEffetto in days
	Frazionamento     string `json:"frazionamento"`     // "ANNUALE", "SEMESTRALE", "QUADRIMESTRALE", "TRIMESTRALE", "BIMESTRALE", "MENSILE", "Unica_soluzione"
}

type totaleGaranzia struct {
	Garanzia         string  `json:"garanzia"`
	Imposte          float64 `json:"imposte"`
	PremioImponibile float64 `json:"premioImponibile"`
	Totale           float64 `json:"totale"`
}

type totale struct {
	Imposte          float64          `json:"imposte"`
	PremioImponibile float64          `json:"premioImponibile"`
	Totale           float64          `json:"totale"`
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
	NumPolSost     string     `json:"numPolSost,omitempty"`
	Ramo           string     `json:"ramo"`
}

// var paymentSplitMap map[string]string = map[string]string{
// 	string(models.PaySplitYearly):  "ANNUALE",
// 	string(models.PaySplitMonthly): "MENSILE",
// }

var guaranteeMap = map[string]string{
	"death":                "VITA",
	"permanent-disability": "INVALIDITA",
	"temporary-disability": "INABILITA",
	"serious-ill":          "MALATTIA",
}

func policyDto(p models.Policy, producer string) policy {
	var (
		wp                  policy
		an                  anagrafica
		pa                  perAss
		totale              totale
		contractorBirthDate time.Time
		err                 error
	)

	// Map contractor data
	if p.Contractor.Residence != nil {
		an.Cap = p.Contractor.Residence.PostalCode
		an.Comune = p.Contractor.Residence.City
	}
	an.Cf = p.Contractor.FiscalCode
	an.Cognome = p.Contractor.Surname
	if contractorBirthDate, err = time.Parse(time.RFC3339, p.Contractor.BirthDate); err != nil {
		return policy{}
	}
	an.DataNascita = contractorBirthDate.Format(time.DateOnly)
	an.Indirizzo = lib.TrimSpace(fmt.Sprintf("%s, %s", p.Contractor.Residence.StreetName, p.Contractor.Residence.StreetNumber))
	an.LuogoNascita = p.Contractor.BirthCity
	an.Nazione = "IT"
	an.NazioneNascita = "IT" // for now we'll hardcode it. To enrich based on fiscalCode
	an.Nome = p.Contractor.Name
	an.Pfg = "F"
	an.Provincia = p.Contractor.Residence.CityCode
	an.ProvinciaNascita = p.Contractor.BirthProvince
	an.Sesso = p.Contractor.Gender

	wp.Anagrafica = an

	// Map guarantees
	wp.Garanzie = make([]garanzia, 0)
	for woptaSlug, winSlug := range guaranteeMap {
		var g garanzia
		g.Garanzia = winSlug
		index := slices.IndexFunc(p.Assets[0].Guarantees, func(item models.Guarante) bool {
			return item.Slug == woptaSlug
		})
		if index != -1 {
			guarantee := p.Assets[0].Guarantees[index]
			g.SommaAssicurata = guarantee.Value.SumInsuredLimitOfIndemnity
			g.PremioImponibile = guarantee.Value.PremiumNetYearly
			g.Imposte = guarantee.Value.PremiumTaxAmountYearly
		}
		wp.Garanzie = append(wp.Garanzie, g)
	}

	wp.IdPratica = p.ProposalNumber

	firstExpiryDate := lib.AddMonths(p.StartDate, 12)
	if p.PaymentSplit == string(models.PaySplitMonthly) {
		firstExpiryDate = lib.AddMonths(p.StartDate, 1)
	}

	pa.BaseAnno = "ANNO_SOLARE"
	pa.DataEffetto = p.StartDate.Format(time.DateOnly)
	pa.DataPrimaScadenza = firstExpiryDate.Format(time.DateOnly)
	pa.DataScadenza = p.EndDate.Format(time.DateOnly)
	pa.DurataIniziale = int(firstExpiryDate.Sub(p.StartDate).Hours() / 24)
	pa.Frazionamento = "ANNUALE"

	wp.PerAss = pa

	wp.Prodotto = "WOPTA_VITA"

	// Map totals
	totale.Imposte = p.TaxAmount
	totale.PremioImponibile = p.PriceNett
	totale.Totale = p.PriceGross
	totale.TotaliGaranzie = make([]totaleGaranzia, 0)
	for _, g := range wp.Garanzie {
		totale.TotaliGaranzie = append(totale.TotaliGaranzie, totaleGaranzia{
			Garanzia:         g.Garanzia,
			Imposte:          g.Imposte,
			PremioImponibile: g.PremioImponibile,
			Totale:           g.Imposte + g.PremioImponibile,
		})
	}
	wp.TotaleAnnuo = totale
	wp.TotaleFirma = totale
	wp.TotaleFutura = totale

	if p.PaymentSplit == string(models.PaySplitMonthly) {
		wp.TotaleFirma.Imposte = p.TaxAmountMonthly
		wp.TotaleFirma.PremioImponibile = p.PriceNettMonthly
		wp.TotaleFirma.Totale = p.PriceGrossMonthly
		for i := range wp.TotaleFirma.TotaliGaranzie {
			wp.TotaleFirma.TotaliGaranzie[i].Imposte /= 12
			wp.TotaleFirma.TotaliGaranzie[i].PremioImponibile /= 12
			wp.TotaleFirma.TotaliGaranzie[i].Totale = wp.TotaleFirma.TotaliGaranzie[i].Imposte + wp.TotaleFirma.TotaliGaranzie[i].PremioImponibile
		}

		wp.TotaleFutura.Imposte = wp.TotaleAnnuo.Imposte - wp.TotaleFirma.Imposte
		wp.TotaleFutura.PremioImponibile = wp.TotaleAnnuo.PremioImponibile - wp.TotaleFirma.PremioImponibile
		wp.TotaleFutura.Totale = wp.TotaleAnnuo.Totale - wp.TotaleFirma.Totale
		for i := range wp.TotaleFutura.TotaliGaranzie {
			wp.TotaleFutura.TotaliGaranzie[i].Imposte = wp.TotaleAnnuo.TotaliGaranzie[i].Imposte - wp.TotaleFirma.TotaliGaranzie[i].Imposte
			wp.TotaleFutura.TotaliGaranzie[i].PremioImponibile = wp.TotaleAnnuo.TotaliGaranzie[i].PremioImponibile - wp.TotaleFirma.TotaliGaranzie[i].PremioImponibile
			wp.TotaleFutura.TotaliGaranzie[i].Totale = wp.TotaleAnnuo.TotaliGaranzie[i].Totale - wp.TotaleFirma.TotaliGaranzie[i].Totale
		}
	}

	wp.Utente = producer

	wp.DtEmissione = p.EmitDate.Format(time.DateOnly)

	wp.LuogoEmissione = "Italia"

	wp.NumOriginali = 1

	wp.NumPol = p.CodeCompany

	wp.Ramo = "01"

	return wp
}
