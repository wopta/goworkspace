package models

import "time"

type WiseContract struct {
	PolicyContractReferenceID           string    `json:"txRifIdPolContratto,omitempty"`
	InstallmentPaymentDeadlineDate      time.Time `json:"dtScadenzaRata,omitempty"`
	InstallmentPaymentEffectiveDate     time.Time `json:"dtEffettoRata,omitempty"`
	DurationCode                        string    `json:"cdDurata,omitempty"`
	DurationDescription                 string    `json:"txDurata,omitempty"`
	InsuranceCollectionMeansCode        string    `json:"cdMezzoIncassoPolizza,omitempty"`
	InsuranceCollectionMeansDescription string    `json:"txMezzoIncassoPolizza,omitempty"`
	IsElectronicInsuranceCollection     bool      `json:"bMezzoIncassoPolizzaElettronico,omitempty"`
	InstalmentTypeCode                  string    `json:"cdFrazionamento,omitempty"`
	InstalmentTypeDescription           string    `json:"txFrazionamento,omitempty"`
	CurrencyCode                        string    `json:"cdValuta,omitempty"`
	CurrencyDescription                 string    `json:"txValuta,omitempty"`
	GrossAmount                         float64   `json:"nImpLordo,omitempty"`
	NetAmount                           float64   `json:"nImpNetto,omitempty"`
	DiscountsIncreasesAmount            float64   `json:"nImpScontiAumento,omitempty"`
	TaxesAmount                         float64   `json:"nImpTasse,omitempty"`
	AccessoriesAmount                   float64   `json:"nImpAccessori,omitempty"`
	IsRenewalTac                        bool      `json:"bTacRinnovo,omitempty"`
	IsExemptFromTaxes                   bool      `json:"bEsentasse,omitempty"`
	IsDerogation                        bool      `json:"bDeroga,omitempty"`
	AdjustmentPeriodCode                string    `json:"cdPeriodoRegolazione,omitempty"`
	IsPremiumRegulation                 bool      `json:"bRegPremio,omitempty"`
	RegulationMonths                    int       `json:"nMesiRegolazione,omitempty"`
	NextRegulationDate                  time.Time `json:"dtProssimaRegolazione,omitempty"`
	IndexTypeCode                       string    `json:"cdTipoIndiciz,omitempty"`
	IndexTypeDescription                string    `json:"txTipoIndiciz,omitempty"`
	IndexationDate                      time.Time `json:"dtIndicizzazione,omitempty"`
	ExtensionMonths                     int       `json:"nMmDurataProroga,omitempty"`
	AccessoriesPercentage               float64   `json:"nPctAccessori,omitempty"`
	SalesNetworkAgencyRefID             string    `json:"txIdRifReteVenditaAgenzia,omitempty"`
	SalesNetworkAgencyCode              string    `json:"cdReteVenditaAgenzia,omitempty"`
	SalesNetworkAgencyDescription       string    `json:"txReteVenditaAgenzia,omitempty"`
	SalesNetworkProducerRefID           string    `json:"txIdRifReteVenditaProduttore,omitempty"`
	SalesNetworkProducerCode            string    `json:"cdReteVenditaProduttore,omitempty"`
	SalesNetworkProducerDescription     string    `json:"txReteVenditaProduttore,omitempty"`
	SalesNetworkProducerLevel           int       `json:"nLivelloReteVenditaProduttore,omitempty"`
	PolicyExpirationDate                time.Time `json:"dtScadenzaPolizza,omitempty"`
	ProvinceCode                        string    `json:"cdModProv,omitempty"`
	ProvinceDescription                 string    `json:"txModProv,omitempty"`
	PolicyAppendixDate                  time.Time `json:"dtAppendice,omitempty"`
	PolicyAppendixNumber                int       `json:"nPrgAppendice,omitempty"`
	DocumentNumber                      any       `json:"txNDocumento,omitempty"`
	AccessoriesPercentage1              float64   `json:"nPctAccessori1,omitempty"`
	IsContractTransferred               bool      `json:"bContrattoCeduto,omitempty"`
	CurrentMoraExpirationDate           time.Time `json:"dtScadenzaMoraAttuale,omitempty"`
	FutureMoraExpirationDate            time.Time `json:"dtScadenzaMoraFutura,omitempty"`
	InclusionLimit                      int       `json:"nLimiteInclusioni,omitempty"`
	InclusionLimitNumber                int       `json:"nPrgLimiteInclusioni,omitempty"`
	CancellationReasonCode              string    `json:"cdMotivoAnnullamento,omitempty"`
	CancellationReasonDescription       string    `json:"txMotivoAnnullamento,omitempty"`
	CancellationDate                    time.Time `json:"dtAnnullamento,omitempty"`
	AnnualityStartDate                  time.Time `json:"dtInizioAnnualita,omitempty"`
	AnnualityEndDate                    time.Time `json:"dtFineAnnualita,omitempty"`
	IsExternal                          bool      `json:"bIddEsterno,omitempty"`
}
