package win

import (
	"time"

	"github.com/wopta/goworkspace/models"
)

type WinPolicy struct {
	DtEmissione    string
	IdPratica      int
	LuogoEmissione string
	NumOriginali   int
	NumPol         string
	NumPolSost     string
	PerAss         struct {
		BaseAnno          string
		DataEffetto       string
		DataPrimaScadenza string
		DataScadenza      string
		DurataIniziale    int
		Frazionamento     string
	}
	Ramo   string
	Utente string
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
