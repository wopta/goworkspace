package models

import (
	"fmt"
	"time"

	"github.com/wopta/goworkspace/lib"
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

	updatePolicyEndDate(&policy, wisePolicy)

	policy.Uid = fmt.Sprintf("wise:%d", wisePolicy.Id)
	policy.Contractor = *wisePolicy.Contractors[0].Registry.ToDomain()
	policy.CodeCompany = wisePolicy.PolicyNumber

	updatePolicyPrice(&policy, wisePolicy)

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

func updatePolicyEndDate(policy *Policy, wisePolicy *WiseCompletePolicy) {
	location, err := time.LoadLocation("Europe/Rome")
	lib.CheckError(err)
	policyEndDate := wisePolicy.Contract.PolicyExpirationDate.In(location)
	_, offset := policyEndDate.Zone()
	policy.EndDate = policyEndDate.Add(time.Duration(time.Second * time.Duration(offset)))
}

func updatePolicyPrice(policy *Policy, wisePolicy *WiseCompletePolicy) {
	policy.PriceGross = wisePolicy.Contract.AnnualGrossPrice
	policy.PriceNett = wisePolicy.Contract.NetAmount
	policy.TaxAmount = policy.PriceGross - policy.PriceNett

	if wisePolicy.Contract.InstalmentTypeCode == WiseYearlyPaymentSplitCode {
		policy.PaymentSplit = string(PaySplitYear)
		return
	}

	policy.PriceGrossMonthly = wisePolicy.Contract.GrossAmount
	policy.PriceNettMonthly = wisePolicy.Contract.NetAmount
	policy.TaxAmountMonthly = policy.PriceGrossMonthly - policy.PriceNettMonthly

	if wisePolicy.Contract.InstalmentTypeCode == WiseMonthlyPaymentSplitCode {
		policy.PaymentSplit = string(PaySplitMonthly)
		policy.PriceGross = policy.PriceGrossMonthly * 12
		policy.PriceNett = policy.PriceNett * 12
		policy.TaxAmount = policy.TaxAmountMonthly * 12
		return
	}
	if wisePolicy.Contract.InstalmentTypeCode == WiseSemestralPaymentSplitCode {
		policy.PaymentSplit = string(PaySplitSemestral)
		policy.PriceGross = policy.PriceGrossMonthly * 2
		policy.PriceNett = policy.PriceNett * 2
		policy.TaxAmount = policy.TaxAmountMonthly * 2
		return
	}
}

type WiseBase64Annex struct {
	Bytes string `json:"fileAllegato"`
}
