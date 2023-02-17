package models

type WiseUserAddressRegistryDto struct {
	IsMainAddress           bool    `json:"bPrincipale,omitempty"`
	IsoCountryCode          string  `json:"cdNazioneIso,omitempty"`
	UicCountryCode          string  `json:"cdNazioneUic,omitempty"`
	AddressTypeCode         string  `json:"cdTipoIndirizzo,omitempty"`
	Latitude                float64 `json:"nLatitudine,omitempty"`
	Longitude               float64 `json:"nLongitudine,omitempty"`
	PostalCode              string  `json:"txCap,omitempty"`
	HouseNumber             string  `json:"txCivico,omitempty"`
	Municipality            string  `json:"txComune,omitempty"`
	AddressDescription      string  `json:"txDescIndirizzo,omitempty"`
	Location                string  `json:"txLocalita,omitempty"`
	CountryName             string  `json:"txNazioneUic,omitempty"`
	AddressId               string  `json:"txRifIdIndirizzo,omitempty"`
	Province                string  `json:"txSiglaProvincia,omitempty"`
	AddressType             string  `json:"txTipoIndirizzo,omitempty"`
	CdComuneMinisteriale    string  `json:"cdComuneMinisteriale,omitempty"`
	CdPrefissoToponomastico string  `json:"cdPrefissoToponomastico,omitempty"`
	TxPrefissoToponomastico string  `json:"txPrefissoToponomastico,omitempty"`
	TxToponimo              string  `json:"txToponimo,omitempty"`
	BGeolocalizzato         bool    `json:"bGeolocalizzato,omitempty"`
	BIndirizzoNonCensito    bool    `json:"bIndirizzoNonCensito,omitempty"`
}
