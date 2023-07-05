package models

import (
	"fmt"
	"time"
)

type WiseCompletePolicy struct {
	Id              int              `json:"idPolizza"`
	Attachments     []WiseAnnex      `json:"elencoAllegati"`
	Assets          *[]WiseAsset     `json:"elencoBeni"`
	ExpiryDate      time.Time        `json:"dtScadenza"`
	PolicyNumber    string           `json:"txNPolizza"`
	Name            string           `json:"descProdotto"`
	Contractors     []WiseContractor `json:"elencoContraenti"`
	Contract        WiseContract     `json:"contratto"`
	ProductTypeCode string           `json:"cdProdotto"`
	Events          WisePolicyEvents `json:"elencoEventi"`
}

func (wisePolicy *WiseCompletePolicy) ToDomain() Policy {
	var (
		policy     Policy
		wiseToken  *string = nil
		attachment Attachment
	)

	policy.Uid = fmt.Sprintf("wise:%d", wisePolicy.Id)
	policy.Contractor = *wisePolicy.Contractors[0].Registry.ToDomain()
	policy.EndDate = wisePolicy.Contract.PolicyExpirationDate
	policy.CodeCompany = wisePolicy.PolicyNumber
	policy.PriceGross = wisePolicy.Contract.GrossAmount
	policy.PriceNett = wisePolicy.Contract.NetAmount
	policy.TaxAmount = wisePolicy.Contract.TaxesAmount
	policy.ProductVersion = "v1"

	switch wisePolicy.ProductTypeCode {
	case "PMIW":
		policy.Name = "pmi"
		policy.NameDesc = "Wopta per Artigiani & Imprese"
		policy.Company = "global"
	case "WPIN":
		policy.Name = "persona"
		policy.NameDesc = "Wopta per te Persona"
		policy.Company = "global"
	default:
	}

	for _, wiseAsset := range *wisePolicy.Assets {
		policy.Assets = append(policy.Assets, wiseAsset.ToDomain(wisePolicy))
	}

	policy.Attachments = &[]Attachment{}
	for _, wiseAttachment := range wisePolicy.Attachments {
		attachment, wiseToken = wiseAttachment.ToDomain(wiseToken)
		*policy.Attachments = append(*policy.Attachments, attachment)
	}

	return policy
}

type WiseBase64Annex struct {
	Bytes string `json:"fileAllegato"`
}
