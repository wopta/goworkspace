package models

import (
	"encoding/json"
)

type FxCatnat struct{}

func (FxCatnat) GetFiscalCodeAteco(input map[string]any) (res string) {
	j, _ := json.Marshal(input)
	var p Policy
	_ = json.Unmarshal(j, &p)
	if p.Contractor.FiscalCode != "" {
		res = p.Contractor.FiscalCode + "/"
	}
	if p.Contractor.VatCode != "" {
		res += p.Contractor.VatCode
	}
	return res
}
