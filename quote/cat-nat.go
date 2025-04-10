package quote

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

type contractor struct {
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

type legalRepresentative struct {
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

type guaranteeList struct {
	GuaranteeCode string `json:"codGaranzia"`
	CapitalAmount int    `json:"importoCapitale"`
}

type assetRequest struct {
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
	GuaranteeList        []guaranteeList `json:"elencoGaranzia"`
}

type catNatRequestDTO struct {
	ProductCode         string              `json:"codiceProdotto"`
	Date                string              `json:"dataEffetto"`
	ExternalReference   string              `json:"riferimentoEsterno"`
	DistributorCode     string              `json:"codiceDistributore"`
	SecondLevelCode     string              `json:"codiceSecondoLivello"`
	ThirdLevelCode      string              `json:"codiceTerzoLivello"`
	Splitting           string              `json:"frazionamento"`
	Emission            string              `json:"emissione"`
	SalesChannel        string              `json:"canaleVendita"`
	Contractor          contractor          `json:"contraente"`
	LegalRepresentative legalRepresentative `json:"legaleRappresentante"`
	Asset               assetRequest        `json:"bene"`
}

type catNatResponseDTO struct {
	PolicyNumber   string          `json:"numeroPolizza,omitempty"`
	ProposalNumber string          `json:"numeroProposta,omitempty"`
	Result         string          `json:"esito,omitempty"`
	AnnualGross    float64         `json:"imp_Lordo_Annuo,omitempty"`
	AnnualNet      float64         `json:"imp_Netto_Annuo,omitempty"`
	AnnualTax      float64         `json:"imp_Tasse_Annuo,omitempty"`
	AssetDetail    []assetResponse `json:"dettaglioBeni,omitempty"`
	Errors         []detail        `json:"errori,omitempty"`
	Reports        []detail        `json:"segnalazioni,omitempty"`
}

type assetResponse struct {
	ProgressiveNumber string            `json:"progressivoBene,omitempty"`
	GrossAmount       float64           `json:"imp_Lordo_Bene,omitempty"`
	NetAmount         float64           `json:"imp_Netto_Bene,omitempty"`
	TaxAmount         float64           `json:"imp_Tasse_Bene,omitempty"`
	GuaranteeDetail   []guaranteeDetail `json:"dettaglioGaranzie,omitempty"`
}

type guaranteeDetail struct {
	GuaranteeCode  string  `json:"codiceGaranzia,omitempty"`
	GuaranteeGross float64 `json:"imp_Lordo_Garanzia,omitempty"`
	GuaranteeNet   float64 `json:"imp_Netto_Garanzia,omitempty"`
	GuaranteeTax   float64 `json:"imp_Tasse_Garanzia,omitempty"`
}

type detail struct {
	Code        string `json:"codice,omitempty"`
	Description string `json:"descrizione,omitempty"`
}

func CatNatFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		err       error
		reqPolicy *models.Policy
	)

	log.SetPrefix("[CatNatFx] ")
	defer func() {
		r.Body.Close()
		if err != nil {
			log.Printf("error: %s", err.Error())
		}
		log.Println("Handler end ---------------------------------------------")
		log.SetPrefix("")
	}()
	log.Println("Handler start -----------------------------------------------")

	_, err = lib.GetAuthTokenFromIdToken(r.Header.Get("Authorization"))
	if err != nil {
		log.Printf("error getting authToken")
		return "", nil, err
	}

	if err = json.NewDecoder(r.Body).Decode(&reqPolicy); err != nil {
		log.Println("error decoding request body")
		return "", nil, err
	}

	catNatRequest, err := buildNetInsuranceDTO(reqPolicy)
	if err != nil {
		log.Printf("error building NetInsurance DTO: %s", err.Error())
		return "", nil, err
	}

	scope := "emettiPolizza_441-029-007"
	tokenUrl := "https://apigatewaydigital.netinsurance.it/Identity/connect/token"
	client := lib.ClientCredentials(os.Getenv("NETINS_ID"), os.Getenv("NETINS_SECRET"), scope, tokenUrl)

	resp, err := netInsuranceQuotation(client, catNatRequest)
	if err != nil {
		log.Printf("error calling NetInsurance api: %s", err.Error())
		return "", nil, err
	}

	out, err := json.Marshal(resp)
	if err != nil {
		log.Println("error encoding response %w", err.Error())
		return "", nil, err
	}

	return string(out), out, err
}

func buildNetInsuranceDTO(policy *models.Policy) (catNatRequestDTO, error) {
	dto := catNatRequestDTO{
		ProductCode:       "007",
		Date:              policy.StartDate.Format("2006-01-02"),
		ExternalReference: policy.Uid,
		DistributorCode:   "0155",
		SecondLevelCode:   "0001",
		ThirdLevelCode:    "00180",
		Splitting:         "01",
		Emission:          "no",
		SalesChannel:      "3",
	}

	contr := contractor{
		PersonalDataType:          "2",
		CompanyName:               policy.Contractor.Name,
		VatNumber:                 policy.Contractor.VatCode,
		FiscalCode:                policy.Contractor.FiscalCode,
		AtecoCode:                 policy.Assets[0].Building.Ateco,
		PostalCode:                policy.Contractor.PostalCode,
		Address:                   policy.Contractor.Address,
		Locality:                  policy.Contractor.Locality,
		CityCode:                  policy.Contractor.CityCode,
		Phone:                     policy.Contractor.Phone,
		Email:                     policy.Contractor.Mail,
		PrivacyConsentDate:        policy.StartDate.Format("2006-01-02"),
		ProcessingConsent:         "no",
		GenericMarketingConsent:   "no",
		MarketingProfilingConsent: "no",
		MarketingActivityConsent:  "no",
		DocumentationFormat:       1,
	}

	dto.Contractor = contr

	var legalRep legalRepresentative
	if policy.Contractors != nil {
		for _, v := range *policy.Contractors {
			if v.IsSignatory {
				legalRep.Name = v.Name
				legalRep.Surname = v.Surname
				legalRep.FiscalCode = v.FiscalCode
				legalRep.PostalCode = v.PostalCode
				legalRep.Address = v.Address
				legalRep.Locality = v.Locality
				legalRep.CityCode = v.CityCode
				legalRep.Phone = v.Phone
				legalRep.Email = v.Mail
				break
			}
		}
	}

	dto.LegalRepresentative = legalRep

	asset := assetRequest{
		ContractorAndTenant:  "si", // TODO
		EarthquakeCoverage:   "no", // TODO
		FloodCoverage:        "no", // TODO
		EarthquakePurchase:   "no",
		FloodPurchase:        "no",
		LandSlidePurchase:    "no",
		PostalCode:           "",
		Address:              "",
		Locality:             "",
		CityCode:             "",
		ConstructionMaterial: 0, // TODO
		ConstructionYear:     0, // TODO
		FloorNumber:          0, // TODO
		LowestFloor:          0, // TODO
		GuaranteeList:        make([]guaranteeList, 0),
	}

	for _, v := range policy.Assets {
		if v.Building != nil {
			asset.PostalCode = v.Building.PostalCode
			asset.Address = v.Building.Address
			asset.Locality = v.Building.Locality
			asset.CityCode = v.Building.CityCode
		}
		for _, g := range v.Guarantees {
			if g.Slug == "earthquake" { // TODO check slug
				asset.EarthquakePurchase = "si"
				if g.Value != nil {
					if g.Value.SumInsuredLimitOfIndemnity != 0 {
						var gL guaranteeList
						gL.GuaranteeCode = "211/00"
						gL.CapitalAmount = int(g.Value.SumInsuredLimitOfIndemnity)
						asset.GuaranteeList = append(asset.GuaranteeList, gL)
					}
					if g.Value.SumInsured != 0 {
						var gL guaranteeList
						gL.GuaranteeCode = "211/01"
						gL.CapitalAmount = int(g.Value.SumInsured)
						asset.GuaranteeList = append(asset.GuaranteeList, gL)
					}
					if g.Value.LimitOfIndemnity != 0 {
						var gL guaranteeList
						gL.GuaranteeCode = "211/02"
						gL.CapitalAmount = int(g.Value.LimitOfIndemnity)
						asset.GuaranteeList = append(asset.GuaranteeList, gL)
					}
				}
			}
			if g.Slug == "flood" { // TODO check slug
				asset.FloodPurchase = "si"
				if g.Value != nil {
					if g.Value.SumInsuredLimitOfIndemnity != 0 {
						var gL guaranteeList
						gL.GuaranteeCode = "212/00"
						gL.CapitalAmount = int(g.Value.SumInsuredLimitOfIndemnity)
						asset.GuaranteeList = append(asset.GuaranteeList, gL)
					}
					if g.Value.SumInsured != 0 {
						var gL guaranteeList
						gL.GuaranteeCode = "212/01"
						gL.CapitalAmount = int(g.Value.SumInsured)
						asset.GuaranteeList = append(asset.GuaranteeList, gL)
					}
					if g.Value.LimitOfIndemnity != 0 {
						var gL guaranteeList
						gL.GuaranteeCode = "212/02"
						gL.CapitalAmount = int(g.Value.LimitOfIndemnity)
						asset.GuaranteeList = append(asset.GuaranteeList, gL)
					}
				}
			}
			if g.Slug == "landslide" { // TODO check slug
				asset.LandSlidePurchase = "si"
				if g.Value != nil {
					if g.Value.SumInsuredLimitOfIndemnity != 0 {
						var gL guaranteeList
						gL.GuaranteeCode = "209/00"
						gL.CapitalAmount = int(g.Value.SumInsuredLimitOfIndemnity)
						asset.GuaranteeList = append(asset.GuaranteeList, gL)
					}
					if g.Value.SumInsured != 0 {
						var gL guaranteeList
						gL.GuaranteeCode = "209/01"
						gL.CapitalAmount = int(g.Value.SumInsured)
						asset.GuaranteeList = append(asset.GuaranteeList, gL)
					}
					if g.Value.LimitOfIndemnity != 0 {
						var gL guaranteeList
						gL.GuaranteeCode = "209/02"
						gL.CapitalAmount = int(g.Value.LimitOfIndemnity)
						asset.GuaranteeList = append(asset.GuaranteeList, gL)
					}
				}
			}
		}
	}

	return dto, nil
}

func netInsuranceQuotation(cl *http.Client, dto catNatRequestDTO) (catNatResponseDTO, error) {
	url := "https://apigatewaydigital.netinsurance.it/PolizzeGateway24/emettiPolizza/441-029-007"
	reqBodyBytes := new(bytes.Buffer)
	err := json.NewEncoder(reqBodyBytes).Encode(dto)
	if err != nil {
		return catNatResponseDTO{}, err
	}
	r := reqBodyBytes.Bytes()
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(r))
	req.Header.Set("Content-Type", "application/json")
	resp, err := cl.Do(req)
	if err != nil {
		return catNatResponseDTO{}, err
	}
	if resp != nil {
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return catNatResponseDTO{}, fmt.Errorf(string(body))
		}
		//body, err := io.ReadAll(res.Body)
		//lib.CheckError(err)
		cndto := catNatResponseDTO{}
		if err = json.NewDecoder(resp.Body).Decode(&cndto); err != nil {
			log.Println("error decoding catnat response")
			return catNatResponseDTO{}, err
		}

		resp.Body.Close()
		return cndto, nil

	}
	return catNatResponseDTO{}, fmt.Errorf("empty catnat response")
}
