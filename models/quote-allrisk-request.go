package models

type QuoteAllriskRequestMunichRe struct {
	Sme struct {
		Ateco        string  `json:"ateco"`
		SubproductID float64 `json:"subproductId"`
		UWRole       string  `json:"UW_role"`
		Company      struct {
			Vatnumber struct {
				Value string `json:"value"`
			} `json:"vatnumber"`
			OpreEur struct {
				Value float64 `json:"value"`
			} `json:"opre_eur"`
			Employees struct {
				Value float64 `json:"value"`
			} `json:"employees"`
		} `json:"company"`
		Answers struct {
			Step1 []struct {
				Slug  string `json:"slug"`
				Value struct {
					TypeOfSumInsured           string  `json:"typeOfSumInsured"`
					Deductible                 string  `json:"deductible"`
					SumInsuredLimitOfIndemnity float64 `json:"sumInsuredLimitOfIndemnity"`
				} `json:"value"`
			} `json:"step1"`
			Step2 []struct {
				BuildingID string `json:"buildingId"`
				Value      struct {
					BuildingType     string `json:"buildingType"`
					NumberOfFloors   string `json:"numberOfFloors"`
					ConstructionYear string `json:"constructionYear"`
					Alarm            string `json:"alarm"`
					TypeOfInsurance  string `json:"typeOfInsurance"`
					Ateco            string `json:"ateco"`
					Postcode         string `json:"postcode"`
					Province         string `json:"province"`
					Answer           []struct {
						Slug  string `json:"slug"`
						Value struct {
							TypeOfSumInsured           string  `json:"typeOfSumInsured"`
							Deductible                 string  `json:"deductible"`
							SumInsuredLimitOfIndemnity float64 `json:"sumInsuredLimitOfIndemnity"`
						} `json:"value"`
					} `json:"answer"`
				} `json:"value"`
			} `json:"step2"`
		} `json:"answers"`
	} `json:"sme"`
}
