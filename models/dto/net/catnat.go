package net

import "github.com/wopta/goworkspace/models"

type Contractor struct {
	PersonalDataType          string `json:"tipoAnagrafica"`
	CompanyName               string `json:"ragioneSociale"`
	VatNumber                 string `json:"partitaIva"`
	FiscalCode                string `json:"codiceFiscale"`
	AtecoCode                 string `json:"codiceAteco"`
	PostalCode                string `json:"cap"`
	Address                   string `json:"indirizzo"`
	Locality                  string `json:"comune"`
	CityCode                  string `json:"provincia"`
	Phone                     string `json:"telefonoCellulare"`
	Email                     string `json:"email"`
	PrivacyConsentDate        string `json:"dataConsensoPrivacy"`
	ProcessingConsent         string `json:"consensoTrattamento"`
	GenericMarketingConsent   string `json:"consensoMarketingGenerico"`
	MarketingProfilingConsent string `json:"consensoProfilazioneMarketing"`
	MarketingActivityConsent  string `json:"consensoAttivitaMarketing"`
	DocumentationFormat       int    `json:"formatoDocumentazione"`
}

type LegalRepresentative struct {
	Name       string `json:"nome"`
	Surname    string `json:"cognome"`
	FiscalCode string `json:"codiceFiscale"`
	PostalCode string `json:"cap"`
	Address    string `json:"indirizzo"`
	Locality   string `json:"comune"`
	CityCode   string `json:"provincia"`
	Phone      string `json:"telefonoCellulare"`
	Email      string `json:"email"`
}

type GuaranteeList struct {
	GuaranteeCode string `json:"codGaranzia"`
	CapitalAmount int    `json:"importoCapitale"`
}

type AssetRequest struct {
	ContractorAndTenant  string          `json:"contraenteProprietarioEConduttore"`
	EarthquakeCoverage   string          `json:"presenzaCoperturaTerremoto"`
	FloodCoverage        string          `json:"presenzaCoperturaAlluvione"`
	EarthquakePurchase   string          `json:"acquistoTerremoto"`
	FloodPurchase        string          `json:"acquistoAlluvione"`
	LandSlidePurchase    string          `json:"acquistoFrane"`
	PostalCode           string          `json:"cap"`
	Address              string          `json:"indirizzo"`
	Locality             string          `json:"comune"`
	CityCode             string          `json:"provincia"`
	ConstructionMaterial int             `json:"materialeDiCostruzione"`
	ConstructionYear     int             `json:"annoDiCostruzione"`
	FloorNumber          int             `json:"numeroPianiEdificio"`
	LowestFloor          int             `json:"pianoPiuBassoOccupato"`
	GuaranteeList        []GuaranteeList `json:"elencoGaranzia"`
}

type RequestDTO struct {
	ProductCode         string              `json:"codiceProdotto"`
	Date                string              `json:"dataEffetto"`
	ExternalReference   string              `json:"riferimentoEsterno"`
	DistributorCode     string              `json:"codiceDistributore"`
	SecondLevelCode     string              `json:"codiceSecondoLivello"`
	ThirdLevelCode      string              `json:"codiceTerzoLivello"`
	Splitting           string              `json:"frazionamento"`
	Emission            string              `json:"emissione"`
	SalesChannel        string              `json:"canaleVendita"`
	Contractor          Contractor          `json:"contraente"`
	LegalRepresentative LegalRepresentative `json:"legaleRappresentante"`
	Asset               AssetRequest        `json:"bene"`
}

type ResponseDTO struct {
	PolicyNumber   string          `json:"numeroPolizza,omitempty"`
	ProposalNumber string          `json:"numeroProposta,omitempty"`
	Result         string          `json:"esito,omitempty"`
	AnnualGross    float64         `json:"imp_Lordo_Annuo,omitempty"`
	AnnualNet      float64         `json:"imp_Netto_Annuo,omitempty"`
	AnnualTax      float64         `json:"imp_Tasse_Annuo,omitempty"`
	AssetDetail    []AssetResponse `json:"dettaglioBeni,omitempty"`
	Errors         []Detail        `json:"errori,omitempty"`
	Reports        []Detail        `json:"segnalazioni,omitempty"`
}

type AssetResponse struct {
	ProgressiveNumber string            `json:"progressivoBene,omitempty"`
	GrossAmount       float64           `json:"imp_Lordo_Bene,omitempty"`
	NetAmount         float64           `json:"imp_Netto_Bene,omitempty"`
	TaxAmount         float64           `json:"imp_Tasse_Bene,omitempty"`
	GuaranteeDetail   []GuaranteeDetail `json:"dettaglioGaranzie,omitempty"`
}

type GuaranteeDetail struct {
	GuaranteeCode  string  `json:"codiceGaranzia,omitempty"`
	GuaranteeGross float64 `json:"imp_Lordo_Garanzia,omitempty"`
	GuaranteeNet   float64 `json:"imp_Netto_Garanzia,omitempty"`
	GuaranteeTax   float64 `json:"imp_Tasse_Garanzia,omitempty"`
}

type Detail struct {
	Code        string `json:"codice,omitempty"`
	Description string `json:"descrizione,omitempty"`
}

type ErrorResponse struct {
	Type     string         `json:"type"`
	Title    string         `json:"title"`
	Status   int            `json:"status"`
	Detail   string         `json:"detail"`
	Instance string         `json:"instance"`
	Errors   map[string]any `json:"errors"`
}

const earthquakeCode = "211"
const floodCode = "212"
const landslideCode = "209"
const buildingCode = "/00"
const contentCode = "/01"
const stockCode = "/02"

const earthquakeBuildingCode = earthquakeCode + buildingCode
const earthquakeContentCode = earthquakeCode + contentCode
const earthquakeStockCode = earthquakeCode + stockCode
const floodBuildingCode = floodCode + buildingCode
const floodContentCode = floodCode + contentCode
const floodStockCode = floodCode + stockCode
const landslideBuildingCode = landslideCode + buildingCode
const landslideContentCode = landslideCode + contentCode
const landslideStockCode = landslideCode + stockCode

const earthquakeSlug = "earthquake"
const floodSlug = "flood"
const landslideSlug = "landslide"

const yes = "si"
const no = "no"

func (d *RequestDTO) FromPolicy(p *models.Policy) error {

	const catNatProductCode = "007"
	const catNatDistributorCode = "0155"
	const catNatSecondLevelCode = "0001"
	const catNatThirdLevelCode = "00180"
	const catNatSplitting = "01"
	const catNatSalesChannel = "3"
	const catNatPersonalDataType = "2"

	d.ProductCode = catNatProductCode
	d.Date = p.StartDate.Format("2006-01-02")
	d.ExternalReference = p.Uid
	d.DistributorCode = catNatDistributorCode
	d.SecondLevelCode = catNatSecondLevelCode
	d.ThirdLevelCode = catNatThirdLevelCode
	d.Splitting = catNatSplitting
	d.Emission = no
	d.SalesChannel = catNatSalesChannel

	var atecoCode string
	for _, v := range p.Assets {
		if v.Building != nil {
			atecoCode = v.Building.Ateco
		}
	}
	contr := Contractor{
		PersonalDataType:          catNatPersonalDataType,
		CompanyName:               p.Contractor.Name,
		VatNumber:                 p.Contractor.VatCode,
		FiscalCode:                p.Contractor.FiscalCode,
		AtecoCode:                 atecoCode,
		Phone:                     p.Contractor.Phone,
		Email:                     p.Contractor.Mail,
		PrivacyConsentDate:        p.StartDate.Format("2006-01-02"),
		ProcessingConsent:         no,
		GenericMarketingConsent:   no,
		MarketingProfilingConsent: no,
		MarketingActivityConsent:  no,
		DocumentationFormat:       1,
	}
	if p.Contractor.Residence != nil {
		contr.Address = formatAddress(p.Contractor.Residence)
		contr.Locality = p.Contractor.Residence.Locality
		contr.CityCode = p.Contractor.Residence.CityCode
	}

	d.Contractor = contr

	var legalRep LegalRepresentative
	if p.Contractors != nil {
		for _, v := range *p.Contractors {
			if v.IsSignatory {
				legalRep.Name = v.Name
				legalRep.Surname = v.Surname
				legalRep.FiscalCode = v.FiscalCode
				legalRep.Phone = v.Phone
				legalRep.Email = v.Mail
				if v.Residence != nil {
					legalRep.Address = formatAddress(v.Residence)
					legalRep.PostalCode = v.Residence.PostalCode
					legalRep.Locality = v.Residence.Locality
					legalRep.CityCode = v.Residence.CityCode
				}
				break
			}
		}
	}

	d.LegalRepresentative = legalRep

	asset := AssetRequest{
		ContractorAndTenant:  yes, // TODO
		EarthquakeCoverage:   no,  // TODO
		FloodCoverage:        no,  // TODO
		EarthquakePurchase:   no,
		FloodPurchase:        no,
		LandSlidePurchase:    no,
		ConstructionMaterial: 0, // TODO
		ConstructionYear:     0, // TODO
		FloorNumber:          0, // TODO
		LowestFloor:          0, // TODO
		GuaranteeList:        make([]GuaranteeList, 0),
	}

	for _, v := range p.Assets {
		if v.Building != nil {
			if v.Building.BuildingAddress != nil {
				asset.PostalCode = v.Building.BuildingAddress.PostalCode
				asset.Address = formatAddress(v.Building.BuildingAddress)
				asset.Locality = v.Building.BuildingAddress.Locality
				asset.CityCode = v.Building.BuildingAddress.CityCode
			}
		}
		for _, g := range v.Guarantees {
			if g.Slug == earthquakeSlug { // TODO check slug
				asset.EarthquakePurchase = yes
				if g.Value != nil {
					if g.Value.SumInsuredLimitOfIndemnity != 0 {
						var gL GuaranteeList
						gL.GuaranteeCode = earthquakeBuildingCode
						gL.CapitalAmount = int(g.Value.SumInsuredLimitOfIndemnity)
						asset.GuaranteeList = append(asset.GuaranteeList, gL)
					}
					if g.Value.SumInsured != 0 {
						var gL GuaranteeList
						gL.GuaranteeCode = earthquakeContentCode
						gL.CapitalAmount = int(g.Value.SumInsured)
						asset.GuaranteeList = append(asset.GuaranteeList, gL)
					}
					if g.Value.LimitOfIndemnity != 0 {
						var gL GuaranteeList
						gL.GuaranteeCode = earthquakeStockCode
						gL.CapitalAmount = int(g.Value.LimitOfIndemnity)
						asset.GuaranteeList = append(asset.GuaranteeList, gL)
					}
				}
			}
			if g.Slug == floodSlug { // TODO check slug
				asset.FloodPurchase = yes
				if g.Value != nil {
					if g.Value.SumInsuredLimitOfIndemnity != 0 {
						var gL GuaranteeList
						gL.GuaranteeCode = floodBuildingCode
						gL.CapitalAmount = int(g.Value.SumInsuredLimitOfIndemnity)
						asset.GuaranteeList = append(asset.GuaranteeList, gL)
					}
					if g.Value.SumInsured != 0 {
						var gL GuaranteeList
						gL.GuaranteeCode = floodContentCode
						gL.CapitalAmount = int(g.Value.SumInsured)
						asset.GuaranteeList = append(asset.GuaranteeList, gL)
					}
					if g.Value.LimitOfIndemnity != 0 {
						var gL GuaranteeList
						gL.GuaranteeCode = floodStockCode
						gL.CapitalAmount = int(g.Value.LimitOfIndemnity)
						asset.GuaranteeList = append(asset.GuaranteeList, gL)
					}
				}
			}
			if g.Slug == landslideSlug { // TODO check slug
				asset.LandSlidePurchase = yes
				if g.Value != nil {
					if g.Value.SumInsuredLimitOfIndemnity != 0 {
						var gL GuaranteeList
						gL.GuaranteeCode = landslideBuildingCode
						gL.CapitalAmount = int(g.Value.SumInsuredLimitOfIndemnity)
						asset.GuaranteeList = append(asset.GuaranteeList, gL)
					}
					if g.Value.SumInsured != 0 {
						var gL GuaranteeList
						gL.GuaranteeCode = landslideContentCode
						gL.CapitalAmount = int(g.Value.SumInsured)
						asset.GuaranteeList = append(asset.GuaranteeList, gL)
					}
					if g.Value.LimitOfIndemnity != 0 {
						var gL GuaranteeList
						gL.GuaranteeCode = landslideStockCode
						gL.CapitalAmount = int(g.Value.LimitOfIndemnity)
						asset.GuaranteeList = append(asset.GuaranteeList, gL)
					}
				}
			}
		}
	}
	d.Asset = asset

	return nil
}

func (d *ResponseDTO) ToPolicy(p *models.Policy) error {
	eOffer := make(map[string]*models.GuaranteValue)
	fOffer := make(map[string]*models.GuaranteValue)
	lOffer := make(map[string]*models.GuaranteValue)

	for _, a := range d.AssetDetail {
		for _, g := range a.GuaranteeDetail {
			if g.GuaranteeCode == earthquakeBuildingCode {
				eOffer["default"].SumInsuredLimitOfIndemnity = g.GuaranteeGross
			}
			if g.GuaranteeCode == earthquakeContentCode {
				eOffer["default"].SumInsured = g.GuaranteeGross
			}
			if g.GuaranteeCode == earthquakeStockCode {
				eOffer["default"].LimitOfIndemnity = g.GuaranteeGross
			}
			if g.GuaranteeCode == floodBuildingCode {
				fOffer["default"].SumInsuredLimitOfIndemnity = g.GuaranteeGross
			}
			if g.GuaranteeCode == floodContentCode {
				fOffer["default"].SumInsured = g.GuaranteeGross
			}
			if g.GuaranteeCode == floodStockCode {
				fOffer["default"].LimitOfIndemnity = g.GuaranteeGross
			}
			if g.GuaranteeCode == landslideBuildingCode {
				lOffer["default"].SumInsuredLimitOfIndemnity = g.GuaranteeGross
			}
			if g.GuaranteeCode == landslideContentCode {
				lOffer["default"].SumInsured = g.GuaranteeGross
			}
			if g.GuaranteeCode == landslideStockCode {
				lOffer["default"].LimitOfIndemnity = g.GuaranteeGross
			}
		}
	}

	for _, a := range p.Assets {
		for _, g := range a.Guarantees {
			if g.Slug == earthquakeSlug {
				g.Offer = eOffer
			}
			if g.Slug == floodSlug {
				g.Offer = fOffer
			}
			if g.Slug == landslideSlug {
				g.Offer = lOffer
			}
		}
	}

	return nil
}

func formatAddress(addr *models.Address) string {
	res := addr.StreetName + "," + addr.StreetNumber

	return res
}
