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
	FiscalCode              string                      `json:"txCodiceFiscale,omitempty"`
	Surname                 string                      `json:"txCognome,omitempty"`
	PlaceOfBirth            string                      `json:"txComuneNascita,omitempty"`
	FiscalResidencyNation   string                      `json:"txNazioneIdentFiscale,omitempty"`
	Name                    string                      `json:"txNome,omitempty"`
	VatNumber               string                      `json:"txPartitaIva,omitempty"`
	BusinessName            string                      `json:"txRagioneSociale,omitempty"`
	Gender                  string                      `json:"txSesso,omitempty"`
	Contacts                []WiseContactInfo           `json:"listRecapiti"`
}

type WiseContactInfo struct {
	TypeCode string `json:"cdTipoRecapito"`
	Type     string `json:"txTipoRecapito"`
	Contact  string `json:"txRecapito"`
}

const WISE_PHONE_CONTACT_TYPE_CODE = "4"
const WISE_EMAIL_CONTACT_TYPE_CODE = "5"

func (registry *WiseUserRegistryDto) ToDomain() *User {
	var person User

	person.Name = registry.Name
	person.Surname = registry.Surname
	person.FiscalCode = registry.FiscalCode
	if len(registry.Address.TxToponimo) > 0 {
		person.Address += registry.Address.TxToponimo
	}
	if len(registry.Address.AddressDescription) > 0 {
		person.Address += registry.Address.AddressDescription
	}
	person.StreetNumber = registry.Address.HouseNumber
	person.PostalCode = registry.Address.PostalCode
	person.City = registry.Address.Municipality
	person.BirthDate = registry.BirthDate.Format(time.RFC3339)

	for _, contact := range registry.Contacts {
		switch contact.TypeCode {
		case WISE_PHONE_CONTACT_TYPE_CODE:
			person.Phone = contact.Contact
		case WISE_EMAIL_CONTACT_TYPE_CODE:
			person.Mail = contact.Contact
		}
	}

	return &person
}
