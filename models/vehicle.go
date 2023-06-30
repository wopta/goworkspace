package models

import "time"

type Vehicle struct {
	Plate            string    `firestore:"plate,omitempty"            json:"plate,omitempty"            bigquery:"-"` // TARGA
	Model            string    `firestore:"model,omitempty"            json:"model,omitempty"            bigquery:"-"` // MODELLO
	Manufacturer     string    `firestore:"manufacturer,omitempty"     json:"manufacturer,omitempty"     bigquery:"-"` // PRODUTTORE
	RegistrationDate time.Time `firestore:"registrationDate,omitempty" json:"registrationDate,omitempty" bigquery:"-"` // DATA IMMATRICOLAZIONE
	OwnershipStatus  string    `firestore:"ownershipStatus,omitempty"  json:"ownershipStatus,omitempty"  bigquery:"-"` // GIA’ PROPRIETARIO (si/in attesa/la comprerò in futuro)
	NumberOfOwners   int64     `firestore:"numberOfOwners,omitempty"   json:"numberOfOwners,omitempty"   bigquery:"-"` // NUMERO DI PROPRIETARI
	PriceValue       int64     `firestore:"priceValue,omitempty"       json:"priceValue,omitempty"       bigquery:"-"` // VALORE VEICOLO
	VehicleType      string    `firestore:"vehicleType,omitempty"      json:"vehicleType,omitempty"      bigquery:"-"` // TIPO VEICOLO
	Weight           float64   `firestore:"weight,omitempty"           json:"weight,omitempty"           bigquery:"-"` // PESO VEICOLO IN QUINTALI
	PowerSupply      string    `firestore:"powerSupply,omitempty"      json:"powerSupply,omitempty"      bigquery:"-"` // ALIMENTAZIONE
	HasSatellite     bool      `firestore:"hasSatellite,omitempty"     json:"hasSatellite,omitempty"     bigquery:"-"`

	// NOTE: Unused attributes
	// Vin                string    `firestore:"vin,omitempty"                json:"vin,omitempty"                bigquery:"-"` // VIN (vehicle Identification Number)
	// AlarmTypeInstalled string    `firestore:"alarmTypeInstalled,omitempty" json:"alarmTypeInstalled,omitempty" bigquery:"-"` // ANTIFURTO INSTALLATO (nessuno/meccanico/elettronico/satellitare)
	// Year               string    `firestore:"year,omitempty"               json:"year,omitempty"               bigquery:"-"` // ANNO
	// BodyType           string    `firestore:"bodyType,omitempty"           json:"bodyType,omitempty"           bigquery:"-"` // BODY TYPE (sedan, hatchback, suv, coupe, convertible, wagon, minivan, pickup truck, van, crossover)
	// Setup              string    `firestore:"setup,omitempty"              json:"setup,omitempty"              bigquery:"-"` // ALLESTIMENTO (cilindrata, kw/cv)
	// OwnerType          string    `firestore:"ownerType,omitempty"          json:"ownerType,omitempty"          bigquery:"-"` // TIPOLOGIA PROPRIETARIO (uomo/donna/società-p.iva)
	// HasTowHook         bool      `firestore:"hasTowHook,omitempty"         json:"hasTowHook,omitempty"         bigquery:"-"` // PRESENZA GANCIO TRAINO
	// KmPerYear          int64     `firestore:"kmPerYear,omitempty"          json:"kmPerYear,omitempty"          bigquery:"-"` // KM ANNUI
	// OvernightShelter   string    `firestore:"overnightShelter,omitempty"   json:"overnightShelter,omitempty"   bigquery:"-"` // RICOVERO NOTTURNO (box privato/garage pubblico/area recintata privata/in strada)
	// MainUse            string    `firestore:"mainUse,omitempty"            json:"mainUse,omitempty"            bigquery:"-"` // UTILIZZO PRINCIPALE (tempo libero/tempo libero e casa-lavoro/lavoro)
}
