package models

type WiseAsset struct {
	AssetTypeCode string                 `json:"cdTipoBene"`
	AssetType     string                 `json:"txTipoBene"`
	Guarantees    []WiseGuarantee `json:"elencoGaranzie"`
	Location      *WiseLocation          `json:"ubicazione"`
	Person        *struct {
		Registry WiseUserRegistryDto `json:"anagrafica"`
	} `json:"assicurato"`
}

func (wiseAsset *WiseAsset) ToDomain() Asset {
	var (
		asset Asset
	)

	if wiseAsset.Person != nil {
		asset.Person = wiseAsset.Person.Registry.ToDomain()
	}

	switch wiseAsset.AssetType {
	case "GENERICO":
	case "UBICAZIONE":
	case "ASSICURATO":
	}

	for _, wiseGuarantee := range wiseAsset.Guarantees {
		asset.Guarantees = append(asset.Guarantees, wiseGuarantee.ToDomain())
	}

	return asset
}