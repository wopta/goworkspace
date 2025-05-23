package catnat

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models"
)

type QuoteResponse struct {
	PolicyNumber   string          `json:"numeroPolizza,omitempty"`
	ProposalNumber string          `json:"numeroProposta,omitempty"`
	Result         string          `json:"esito,omitempty"`
	AnnualGross    float64         `json:"imp_Lordo_Annuo,omitempty"`
	AnnualNet      float64         `json:"imp_Netto_Annuo,omitempty"`
	AnnualTax      float64         `json:"imp_Tasse_Annuo,omitempty"`
	AssetsDetail   []AssetResponse `json:"dettaglioBeni,omitempty"`
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
	//leave it with "C" otherwise dosnt work
	ConsensoTrattamento string `json:"ConsensoTrattamento,omitempty"`
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
	GuaranteesDetail  []GuaranteeDetail `json:"dettaglioGaranzie,omitempty"`
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

const buildingCode = "/00"
const contentCode = "/01"
const stockCode = "/02"

const earthquakeSlug = "earthquake"
const floodSlug = "flood"
const landslideSlug = "landslides"

const catNatProductCode = "007"
const catNatDistributorCode = "0168"
const catNatLegalPerson = "2"
const catNatSalesChannel = "3"
const catNatSoleProp = "3"

const yes = "si"
const no = "no"

var guaranteeSlugToCode = map[string]string{
	earthquakeSlug: "211",
	floodSlug:      "212",
	landslideSlug:  "209",
}
var guaranteeCodeToSlug = map[string]string{
	"211": earthquakeSlug,
	"212": floodSlug,
	"209": landslideSlug,
}
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

var splittingMap = map[string]string{
	string(models.PaySplitYearly):     "01",
	string(models.PaySplitSemestral):  "02",
	string(models.PaySplitTrimestral): "04",
}
var quoteQuestionMap = map[bool]string{
	true:  "si",
	false: "no",
}

func (d *QuoteRequest) FromPolicy(policy *models.Policy, isEmission bool) error {
	d.ProductCode = catNatProductCode
	d.Date = policy.StartDate.Format("2006-01-02")
	d.ExternalReference = policy.Uid
	d.DistributorCode = catNatDistributorCode
	split, ok := splittingMap[string(policy.PaymentSplit)]
	if ok {
		d.Splitting = split
	} else {
		log.Printf("Use split 'yearly' since 'PaymentSplit' is wrong '%v'", policy.PaymentSplit)
		d.Splitting = splittingMap[string(models.PaySplitYearly)]
	}
	d.Emission = no
	d.SalesChannel = catNatSalesChannel

	var baseAsset models.Asset
	for _, a := range policy.Assets {
		if a.Building != nil {
			baseAsset = a
			break
		}
	}

	if isEmission {
		if policy.PaymentSplit == "" {
			return fmt.Errorf("No valid split selected for NetInsurance")
		}
		d.Emission = yes
		if policy.Contractor.CompanyAddress == nil {
			return errors.New("You need to populate CompanyAddress")
		}
		if policy.Contractors == nil || len(*policy.Contractors) == 0 {
			return errors.New("You need to compile Contractors")
		}
		var dt string

		if policy.Contractor.VatCode == "" {
			return errors.New("You need to compile Contractor.VatCode")
		}
		if policy.Contractor.Type == "legalEntity" { //persona giuridica
			dt = catNatLegalPerson
			if policy.Contractor.CompanyName == "" {
				return errors.New("You need to compile Contractor.CompanyName")
			}
			policy.Contractor.FiscalCode = policy.Contractor.VatCode
		} else { //ditta individuale i need all date
			dt = catNatSoleProp
			if policy.Contractor.Name == "" {
				return errors.New("You need to compile Contractor.Name")
			}
			if policy.Contractor.Surname == "" {
				return errors.New("You need to compile Contractor.Surname")
			}
			if policy.Contractor.FiscalCode == "" {
				return errors.New("You need to compile Contractor.FiscalCode")
			}
		}
		contr := Contractor{
			PersonalDataType:          dt,
			Name:                      policy.Contractor.Name,
			Surname:                   policy.Contractor.Surname,
			CompanyName:               policy.Contractor.CompanyName,
			VatNumber:                 policy.Contractor.VatCode,
			FiscalCode:                policy.Contractor.FiscalCode,
			AtecoCode:                 policy.Contractor.Ateco,
			Phone:                     policy.Contractor.Phone,
			Email:                     policy.Contractor.Mail,
			PrivacyConsentDate:        policy.StartDate.Format("2006-01-02"),
			ProcessingConsent:         no,
			GenericMarketingConsent:   no,
			MarketingProfilingConsent: no,
			MarketingActivityConsent:  no,
			DocumentationFormat:       1,
			ConsensoTrattamento:       "si",
		}
		if policy.Contractor.CompanyAddress != nil {
			contr.Address = formatAddress(policy.Contractor.CompanyAddress)
			contr.Locality = policy.Contractor.CompanyAddress.Locality
			contr.CityCode = policy.Contractor.CompanyAddress.CityCode
			contr.PostalCode = policy.Contractor.CompanyAddress.PostalCode
		}

		d.Contractor = contr

		var legalRep LegalRepresentative
		if policy.Contractors != nil {
			for _, v := range *policy.Contractors {
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

	alreadyEarthquake := policy.QuoteQuestions["alreadyEarthquake"]
	if alreadyEarthquake == nil {
		return errors.New("missing field alreadyEarthquake")
	}
	wantEarthquake := policy.QuoteQuestions["wantEarthquake"]
	if wantEarthquake == nil {
		wantEarthquake = false
	}
	alreadyFlood := policy.QuoteQuestions["alreadyFlood"]
	if alreadyFlood == nil {
		return errors.New("missing field alreadyFlood")
	}
	wantFlood := policy.QuoteQuestions["wantFlood"]
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
			setGuaranteeValue(&asset, g, guaranteeSlugToCode[g.Slug])
		}
	}
	d.Asset = asset
	return nil
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

func getGuarantee(policy *models.Policy, codeGuarantees string) *models.Guarante {
	slug := guaranteeCodeToSlug[codeGuarantees]
	for i := range policy.Assets[0].Guarantees {
		if policy.Assets[0].Guarantees[i].Slug == slug {
			return &policy.Assets[0].Guarantees[i]
		}
	}
	log.Println("No guarantees found")
	return nil
}
func MappingQuoteResponseToGuarantee(quoteResponse QuoteResponse, policy *models.Policy) {
	var currentGuaranteeCode string
	var currentGuaranteeValueCode string

	for _, assetDetailCatnat := range quoteResponse.AssetsDetail {
		for _, guaranteeDetailCatnat := range assetDetailCatnat.GuaranteesDetail {
			guaranteeCodes := strings.Split(guaranteeDetailCatnat.GuaranteeCode, "/")
			currentGuaranteeCode, currentGuaranteeValueCode = guaranteeCodes[0], "/"+guaranteeCodes[1]
			guarantee := getGuarantee(policy, currentGuaranteeCode)
			switch currentGuaranteeValueCode {
			case buildingCode:
				guarantee.Value.SumInsuredLimitOfIndemnity = guaranteeDetailCatnat.GuaranteeGross
			case contentCode:
				guarantee.Value.SumInsured = guaranteeDetailCatnat.GuaranteeGross
			case stockCode:
				guarantee.Value.LimitOfIndemnity = guaranteeDetailCatnat.GuaranteeGross

			}
		}
	}
}

func MappingQuoteResponseToPolicy(quoteResponse QuoteResponse, policy *models.Policy) {
	eOffer := make(map[string]*models.GuaranteValue)
	fOffer := make(map[string]*models.GuaranteValue)
	lOffer := make(map[string]*models.GuaranteValue)

	eOffer["default"] = new(models.GuaranteValue)
	fOffer["default"] = new(models.GuaranteValue)
	lOffer["default"] = new(models.GuaranteValue)
	for _, asset := range policy.Assets {
		for _, guarantee := range asset.Guarantees {
			if guarantee.Slug == earthquakeSlug {
				guarantee.Offer = eOffer
			}
			if guarantee.Slug == floodSlug {
				guarantee.Offer = fOffer
			}
			if guarantee.Slug == landslideSlug {
				guarantee.Offer = lOffer
			}
		}
	}
	policy.PriceGross = quoteResponse.AnnualGross
	policy.PriceNett = quoteResponse.AnnualNet
	policy.TaxAmount = quoteResponse.AnnualTax
	split := string(models.PaySplitYearly)
	policy.OffersPrices = map[string]map[string]*models.Price{
		"default": {
			split: &models.Price{},
		},
	}
	policy.OffersPrices["default"][split].Gross = policy.PriceGross
	policy.OffersPrices["default"][split].Net = policy.PriceNett
	policy.OffersPrices["default"][split].Tax = policy.TaxAmount
	//TODO: populate in policy->guaranteesMap:>sumInsured sunInsuredLimitOfIndemnity ecc with the different value of response,ex:{"codiceGaranzia":"212/02","imp_Lordo_Garanzia":48.9,"imp_Netto_Garanzia":40,"imp_Tasse_Garanzia":8.9}
	policy.OfferlName = "default"
}

func formatAddress(addr *models.Address) string {
	res := addr.StreetName + "," + addr.StreetNumber
	return res
}
