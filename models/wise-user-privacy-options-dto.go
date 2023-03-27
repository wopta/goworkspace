package models

type WiseUserPrivacyOptionsDto struct {
	ThirdPartyCommercialConsent bool `json:"bConsensoCommercialeVersoTerzi"`
	DataConsent                 bool `json:"bConsensoDati"`
	SensibleDataConsent         bool `json:"bConsensoDatiSensibili"`
	CommercialPurposeConsent    bool `json:"bConsensoFiniCommerciali"`
	MarketingProfilingConsent   bool `json:"bConsensoProfilazioneMarketing"`
}
