package models

type WiseAsset struct {
	AssetTypeCode       string          `json:"cdTipoBene"`
	AssetType           string          `json:"txTipoBene"`
	Guarantees          []WiseGuarantee `json:"elencoGaranzie"`
	Location            *WiseLocation   `json:"ubicazione"`
	ActivityCode        string          `json:"cdAteco"`
	ActivityDescription string          `json:"txAteco"`
	Person              *struct {
		Registry WiseUserRegistryDto `json:"anagrafica"`
	} `json:"assicurato"`
}

func (wiseAsset *WiseAsset) ToDomain(wisePolicy *WiseCompletePolicy) Asset {
	var (
		asset Asset
	)

	if wiseAsset.Person != nil {
		asset.Person = wiseAsset.Person.Registry.ToDomain()
	}

	switch wiseAsset.AssetType {
	case "GENERICO":
		// enterprise
		if wisePolicy != nil {
			var enterprise Enterprise
			enterprise.Name = wisePolicy.Contractors[0].Registry.BusinessName
			enterprise.Ateco = wiseAsset.ActivityCode
			enterprise.AtecoDesc = wiseAsset.ActivityDescription
			asset.Enterprise = &enterprise
		}
	case "UBICAZIONE":
		// fabbricato
		var building Building
		building.Address = wiseAsset.Location.Address.Description
		building.Ateco = wiseAsset.Location.ActivityCode
		building.AtecoDesc = wiseAsset.Location.Activity
		building.City = wiseAsset.Location.Address.Municipality
		building.Name = wiseAsset.Location.Address.Description
		// building.Type =
		building.StreetNumber = wiseAsset.Location.Address.HouseNumber
		building.CityCode = wiseAsset.Location.Address.Province
		building.PostalCode = wiseAsset.Location.Address.PostalCode
		building.Locality = wiseAsset.Location.Address.Locality
		// building.Location
		// building.BuildingType =
		// building.BuildingMateria =
		// building.BuildingYear =
		// building.Employer =
		// building.IsAllarm =
		// building.Floor =
		// building.AtecoMacro =
		// building.AtecoSub =
		// building.Costruction =
		// building.IsHolder =
		asset.Building = &building
	case "ASSICURATO":
	}

	for _, wiseGuarantee := range wiseAsset.Guarantees {
		asset.Guarantees = append(asset.Guarantees, wiseGuarantee.ToDomain())
	}

	return asset
}
