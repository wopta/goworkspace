package models

import (
	"strconv"
	"strings"
	"time"
)

type WiseGuarantee struct {
	Code                       string                    `json:"cdGaranzia"`
	Name                       string                    `json:"txGaranzia"`
	TariffData                 time.Time                 `json:"dtTariffa"`
	TypeCode                   string                    `json:"cdTipologiaGaranzia"`
	Type                       string                    `json:"txTipologiaGaranzia"`
	PriceNet                   float64                   `json:"nImpNetto"`
	PriceGross                 float64                   `json:"nImpLordo"`
	YearlyPrice                float64                   `json:"nImpAnnuo"`
	Tax                        float64                   `json:"nImpTasse"`
	YearlyTax                  float64                   `json:"nImpTasseAnnuo"`
	ExpirationDate             time.Time                 `json:"dtScadenzaGaranzia"`
	NetTransactionAmmount      float64                   `json:"nImpNettoOperazione"`
	NetTransactionTax          float64                   `json:"nImpTasseOperazione"`
	TariffCode                 string                    `json:"cdTariffa"`
	Tariff                     string                    `json:"txTariffa"`
	Deductible                 float64                   `json:"nImpFranchigia"`
	MaxClaimAmmount            float64                   `json:"nImpMaxSinistro"`
	SumInsuredLimitOfIndemnity float64                   `json:"nImpCapitale1"`
	ExcessPercentage           float64                   `json:"nPctScoperto"`
	Parameters                 []WiseGuaranteeParameters `json:"elencoPtfParametri"`
}

type WiseGuaranteeParameters struct {
	Name  string `json:"txParametro"`
	Value string `json:"txValoreParametro"`
}

func (wiseGuarantee *WiseGuarantee) ToDomain() Guarante {
	var guarantee Guarante
	guarantee.Price = wiseGuarantee.PriceGross
	guarantee.PriceGross = wiseGuarantee.PriceGross
	guarantee.PriceNett = wiseGuarantee.PriceNet
	guarantee.Name = wiseGuarantee.Name
	guarantee.CompanyName = wiseGuarantee.Name
	guarantee.Tax = wiseGuarantee.Tax

	if wiseGuarantee.Deductible > 0.001 {
		guarantee.Deductible = strconv.FormatFloat(wiseGuarantee.Deductible, 'f', 2, 64)
	}
	guarantee.SumInsuredLimitOfIndemnity = wiseGuarantee.SumInsuredLimitOfIndemnity

	for _, parameter := range wiseGuarantee.Parameters {
		var guaranteeValue GuaranteValue
		if strings.EqualFold(parameter.Name, "FRANCHIGIA") {
			guaranteeValue.Deductible = parameter.Value
		}
		if strings.EqualFold(parameter.Name, "MASSIMALE") || strings.EqualFold(parameter.Name, "somma assicurata") {
			if value, err := strconv.ParseFloat(strings.ReplaceAll(parameter.Value, ".", ""), 64); err == nil {
				guaranteeValue.SumInsuredLimitOfIndemnity = value
			}
		}
		guarantee.Value = &guaranteeValue
	}

	return guarantee
}

