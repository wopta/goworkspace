package models

import (
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
		policy Policy
	)

	for _, wiseAsset := range *wisePolicy.Assets {
		policy.Assets = append(policy.Assets, wiseAsset.ToDomain())
	}
	policy.Contractor = *wisePolicy.Contractors[0].Registry.ToDomain()
	policy.EndDate = wisePolicy.Contract.PolicyExpirationDate
	// policy.EmitDate = wisePolicy.Events.f

	return policy
}
