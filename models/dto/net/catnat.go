package net

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

type CatNatRequestDTO struct {
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

type CatNatResponseDTO struct {
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
