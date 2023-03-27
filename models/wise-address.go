package models

type WiseAddress struct {
	TypeCode            string  `json:"cdTipoIndirizzo"`
	Type                string  `json:"txTipoIndirizzo"`
	PostalCode          string  `json:"txCap"`
	Description         string  `json:"txDescIndirizzo"`
	Locality            string  `json:"txLocalita"`
	MunicipalityCode    string  `json:"cdComuneMinisteriale"`
	Municipality        string  `json:"txComune"`
	Province            string  `json:"txSiglaProvincia"`
	CountryUicCode      string  `json:"cdNazioneUic"`
	CountryIsoCode      string  `json:"cdNazioneIso"`
	CountryName         string  `json:"txNazioneUic"`
	IsMainAddress       bool    `json:"bPrincipale"`
	ToponymicPrefixCode string  `json:"cdPrefissoToponomastico"`
	ToponymicPrefix     string  `json:"txPrefissoToponomastico"`
	Toponym             string  `json:"txToponimo"`
	HouseNumber         string  `json:"txCivico"`
	Latitude            float64 `json:"nLatitudine"`
	Longitude           float64 `json:"nLongitudine"`
	IsGeolocalized      bool    `json:"bGeolocalizzato"`
	IsUnlisted          bool    `json:"bIndirizzoNonCensito"`
}