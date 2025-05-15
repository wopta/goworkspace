package catnat

import (
	"errors"
	"time"

	"github.com/wopta/goworkspace/lib/log"
	"github.com/wopta/goworkspace/models"
)

type QuoteResponse struct {
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

type QuoteRequest struct {
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

type DownloadRequest struct {
	CodiceCompagnia string    `json:"codiceCompagnia"`
	NumeroPolizza   string    `json:"NumeroPolizza"`
	TipoOperazione  string    `json:"TipoOperazione"`
	DataOperazione  time.Time `json:"DataOperazione"`
}

type DownloadResponse struct {
	Result        string      `json:"esito"`
	NumeroPolizza string      `json:"numeroPolizza"`
	Documento     []Documento `json:"documento"`
	Errors        interface{} `json:"errori"` // or *string if it's always null or a string
}
type Contractor struct {
	Name                      string `json:"nome,omitempty"`
	Surname                   string `json:"cognome,omitempty"`
	PersonalDataType          string `json:"tipoAnagrafica,omitempty"`
	CompanyName               string `json:"ragioneSociale,omitempty"`
	VatNumber                 string `json:"partitaIva,omitempty"`
	FiscalCode                string `json:"codiceFiscale,omitempty"`
	AtecoCode                 string `json:"codiceAteco,omitempty"`
	PostalCode                string `json:"cap,omitempty"`
	Address                   string `json:"indirizzo,omitempty"`
	Locality                  string `json:"comune,omitempty"`
	CityCode                  string `json:"provincia,omitempty"`
	Phone                     string `json:"telefonoCellulare,omitempty"`
	Email                     string `json:"email,omitempty"`
	PrivacyConsentDate        string `json:"dataConsensoPrivacy,omitempty"`
	ProcessingConsent         string `json:"consensoTrattamento,omitempty"`
	GenericMarketingConsent   string `json:"consensoMarketingGenerico,omitempty"`
	MarketingProfilingConsent string `json:"consensoProfilazioneMarketing,omitempty"`
	MarketingActivityConsent  string `json:"consensoAttivitaMarketing,omitempty"`
	DocumentationFormat       int    `json:"formatoDocumentazione"`
	ConsensoTrattamento       string `json:"ConsensoTrattamento,omitempty"`
}

type LegalRepresentative struct {
	Name       string `json:"nome,omitempty"`
	Surname    string `json:"cognome,omitempty"`
	FiscalCode string `json:"codiceFiscale,omitempty"`
	PostalCode string `json:"cap,omitempty"`
	Address    string `json:"indirizzo,omitempty"`
	Locality   string `json:"comune,omitempty"`
	CityCode   string `json:"provincia,omitempty"`
	Phone      string `json:"telefonoCellulare,omitempty"`
	Email      string `json:"email,omitempty"`
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

type Documento struct {
	DescrizioneDocumento string `json:"descrizioneDocumento"`
	DatiDocumento        string `json:"datiDocumento"`
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
const landslideSlug = "landslides"

const catNatProductCode = "007"
const catNatDistributorCode = "0168"
const catNatSplitting = "01"
const catNatLegalPerson = "2"
const catNatSalesChannel = "3"
const catNatSoleProp = "3"

const yes = "si"
const no = "no"

var useTypeMap = map[string]string{
	"owner-tenant": "si",
	"tenant":       "no",
}
var buildingYearMap = map[string]int{
	"before_1950":       1,
	"from_1950_to_1990": 2,
	"after_1990":        3,
	"unknown":           4,
}
var floorMap = map[string]int{
	"up_to_2":     2,
	"more_than_3": 1,
}
var lowestFloorMap = map[string]int{
	"first_floor":  1,
	"upper_floor":  2,
	"ground_floor": 3,
	"underground":  4,
}
var buildingMaterialMap = map[string]int{
	"brick":    1,
	"concrete": 2,
	"steel":    3,
	"unknown":  4,
}
var quoteQuestionMap = map[bool]string{
	true:  "si",
	false: "no",
}

func (d *QuoteRequest) FromPolicy(p *models.Policy, isEmission bool) error {

	d.ProductCode = catNatProductCode
	d.Date = p.StartDate.Format("2006-01-02")
	d.ExternalReference = p.Uid
	d.DistributorCode = catNatDistributorCode
	d.Splitting = catNatSplitting
	d.Emission = no
	d.SalesChannel = catNatSalesChannel

	var baseAsset models.Asset
	for _, a := range p.Assets {
		if a.Building != nil {
			baseAsset = a
			break
		}
	}

	if isEmission {
		d.Emission = yes
		var dt string
		if p.Contractor.Type == "legalEntity" && p.Contractor.FiscalCode == "" {
			dt = catNatLegalPerson
		} else {
			dt = catNatSoleProp
		}
		contr := Contractor{
			PersonalDataType:          dt,
			Name:                      p.Contractor.Name,
			Surname:                   p.Contractor.Surname,
			PostalCode:                p.Contractor.Residence.PostalCode,
			CompanyName:               p.Contractor.Name,
			VatNumber:                 p.Contractor.VatCode,
			FiscalCode:                p.Contractor.FiscalCode,
			AtecoCode:                 baseAsset.Building.Ateco,
			Phone:                     p.Contractor.Phone,
			Email:                     p.Contractor.Mail,
			PrivacyConsentDate:        p.StartDate.Format("2006-01-02"),
			ProcessingConsent:         no,
			GenericMarketingConsent:   no,
			MarketingProfilingConsent: no,
			MarketingActivityConsent:  no,
			DocumentationFormat:       1,
			ConsensoTrattamento:       "si",
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
	}

	alreadyEarthquake := p.QuoteQuestions["alreadyEarthquake"]
	if alreadyEarthquake == nil {
		return errors.New("missing field alreadyEarthquake")
	}
	wantEarthquake := p.QuoteQuestions["wantEarthquake"]
	if wantEarthquake == nil {
		wantEarthquake = false
	}
	alreadyFlood := p.QuoteQuestions["alreadyFlood"]
	if alreadyFlood == nil {
		return errors.New("missing field alreadyFlood")
	}
	wantFlood := p.QuoteQuestions["wantFlood"]
	if wantFlood == nil {
		wantFlood = false
	}
	asset := AssetRequest{
		ContractorAndTenant:  useTypeMap[baseAsset.Building.UseType],
		EarthquakeCoverage:   quoteQuestionMap[alreadyEarthquake.(bool)],
		FloodCoverage:        quoteQuestionMap[alreadyFlood.(bool)],
		EarthquakePurchase:   quoteQuestionMap[(alreadyEarthquake.(bool) && wantEarthquake.(bool)) || !alreadyEarthquake.(bool)],
		FloodPurchase:        quoteQuestionMap[(alreadyFlood.(bool) && wantFlood.(bool)) || !alreadyFlood.(bool)],
		LandSlidePurchase:    no,
		ConstructionMaterial: buildingMaterialMap[baseAsset.Building.BuildingMaterial],
		ConstructionYear:     buildingYearMap[baseAsset.Building.BuildingYear],
		FloorNumber:          floorMap[baseAsset.Building.Floor],
		LowestFloor:          lowestFloorMap[baseAsset.Building.LowestFloor],
		GuaranteeList:        make([]GuaranteeList, 0),
	}
	if baseAsset.Building.BuildingAddress != nil {
		asset.PostalCode = baseAsset.Building.BuildingAddress.PostalCode
		asset.Address = formatAddress(baseAsset.Building.BuildingAddress)
		asset.Locality = baseAsset.Building.BuildingAddress.Locality
		asset.CityCode = baseAsset.Building.BuildingAddress.CityCode
	}
	log.Println("Managing slug guarantees")
	for _, g := range baseAsset.Guarantees {
		if g.IsSelected {
			setGuaranteeValue(&asset, g, mapCodeFromSlug(g.Slug))
		}
	}
	d.Asset = asset

	return nil
}

func mapCodeFromSlug(slug string) string {
	switch slug {
	case earthquakeSlug:
		log.Println("adding earthquake ")
		return earthquakeCode
	case floodSlug:
		log.Println("adding flood")
		return floodCode
	case landslideSlug:
		log.Println("adding landslides")
		return landslideCode
	}
	return ""
}

func setGuaranteeValue(asset *AssetRequest, guarantee models.Guarante, code string) {
	var gL GuaranteeList
	if guarantee.Value.SumInsuredLimitOfIndemnity != 0 {
		gL.GuaranteeCode = code + buildingCode
		gL.CapitalAmount = int(guarantee.Value.SumInsuredLimitOfIndemnity)
		asset.GuaranteeList = append(asset.GuaranteeList, gL)
	}
	if guarantee.Value.SumInsured != 0 {
		gL.GuaranteeCode = code + contentCode
		gL.CapitalAmount = int(guarantee.Value.SumInsured)
		asset.GuaranteeList = append(asset.GuaranteeList, gL)
	}
	if guarantee.Value.LimitOfIndemnity != 0 {
		gL.GuaranteeCode = code + stockCode
		gL.CapitalAmount = int(guarantee.Value.LimitOfIndemnity)
		asset.GuaranteeList = append(asset.GuaranteeList, gL)
	}
}

func (d *QuoteResponse) ToPolicy(p *models.Policy) {
	eOffer := make(map[string]*models.GuaranteValue)
	fOffer := make(map[string]*models.GuaranteValue)
	lOffer := make(map[string]*models.GuaranteValue)

	eOffer["default"] = new(models.GuaranteValue)
	fOffer["default"] = new(models.GuaranteValue)
	lOffer["default"] = new(models.GuaranteValue)
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
	p.PriceGross = d.AnnualGross
	p.PriceNett = d.AnnualNet
	p.TaxAmount = d.AnnualTax
	p.OffersPrices = map[string]map[string]*models.Price{
		"default": {
			"yearly": &models.Price{},
		},
	}
	p.OffersPrices["default"]["yearly"].Gross = p.PriceGross
	p.OffersPrices["default"]["yearly"].Net = p.PriceNett
	p.OffersPrices["default"]["yearly"].Tax = p.TaxAmount

}

func formatAddress(addr *models.Address) string {
	res := addr.StreetName + "," + addr.StreetNumber

	return res
}
