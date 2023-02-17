package models

import (
	"time"
)

type WiseUserRegistryDto struct {
	RegistryType            string                      `json:"txTipoAnagrafica,omitempty"`
	RegistryTypeCode        string                      `json:"cdTipoAnagrafica,omitempty"`
	AtecoCode               string                      `json:"codiceAteco,omitempty"`
	EstabilishmentDate      time.Time                   `json:"dataCostituzione,omitempty"`
	PhysicalPersonExtraData interface{}                 `json:"datiAggPFisica,omitempty"`
	LegalPersonExtraData    interface{}                 `json:"datiAggPGiuridica,omitempty"`
	UserPrivacyOptions      *WiseUserPrivacyOptionsDto  `json:"datiPrivacy,omitempty"`
	AtecoDescription        string                      `json:"descrizioneAteco,omitempty"`
	BirthDate               time.Time                   `json:"dtNascita,omitempty"`
	LastUpdatedAt           time.Time                   `json:"dtUltimaVariazione,omitempty"`
	Id                      int                         `json:"id,omitempty"`
	Address                 *WiseUserAddressRegistryDto `json:"indirizzo,omitempty"`
	ContactDetails          []interface{}               `json:"listRecapiti,omitempty"`
	FiscalCode              string                      `json:"txCodiceFiscale,omitempty"`
	Surname                 string                      `json:"txCognome,omitempty"`
	PlaceOfBirth            string                      `json:"txComuneNascita,omitempty"`
	FiscalResidencyNation   string                      `json:"txNazioneIdentFiscale,omitempty"`
	Name                    string                      `json:"txNome,omitempty"`
	VatNumber               string                      `json:"txPartitaIva,omitempty"`
	BusinessName            string                      `json:"txRagioneSociale,omitempty"`
	Gender                  string                      `json:"txSesso,omitempty"`
}
