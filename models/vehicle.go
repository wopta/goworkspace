package models

import "time"

/*
GIA’ PROPRIETARIO (si/in attesa/la comprerò in futuro)
DATA IMMATRICOLAZIONE
ANTIFURTO INSTALLATO (nessuno/meccanico/elettronico/satellitare)
UTILIZZO PRINCIPALE (tempo libero/tempo libero e casa-lavoro/lavoro)
RICOVERO NOTTURNO (box privato/garage pubblico/area recintata privata/in strada)
KM ANNUI
PRESENZA GANCIO TRAINO
TIPOLOGIA PROPRIETARIO (uomo/donna/società-p.iva)
TARGA
MODELLO
PRODUTTORE
ANNO
VIN (vehicle Identification Number)
BODY TYPE (sedan, hatchback, suv, coupe, convertible, wagon, minivan, pickup truck, van, crossover)
ALLESTIMENTO (cilindrata, kw/cv)
VALORE VEICOLO
*/

type Vehicle struct {
	Plate                   string    `firestore:"plate,omitempty" json:"plate,omitempty" bigquery:"-"`
	Model                   string    `firestore:"model,omitempty" json:"model,omitempty" bigquery:"-"`
	Manufacturer            string    `firestore:"manufacturer,omitempty" json:"manufacturer,omitempty" bigquery:"-"`
	Year                    string    `firestore:"year,omitempty" json:"year,omitempty" bigquery:"-"`
	RegistrationDate        time.Time `firestore:"registrationDate,omitempty" json:"registrationDate,omitempty" bigquery:"-"`
	Vin                     string    `firestore:"vin,omitempty" json:"vin,omitempty" bigquery:"-"`
	BodyType                string    `firestore:"bodyType,omitempty" json:"bodyType,omitempty" bigquery:"-"`
	Setup                   string    `firestore:"setup,omitempty" json:"setup,omitempty" bigquery:"-"`
	VehicleOwnerType        string    `firestore:"vehicleOwnerType,omitempty" json:"vehicleOwnerType,omitempty" bigquery:"-"`
	HasTowHook              bool      `firestore:"hasTowHook,omitempty" json:"hasTowHook,omitempty" bigquery:"-"`
	KmPerYear               int64     `firestore:"kmPerYear,omitempty" json:"kmPerYear,omitempty" bigquery:"-"`
	OvernightVehicleShelter string    `firestore:"overnightVehicleShelter,omitempty" json:"overnightVehicleShelter,omitempty" bigquery:"-"`
	MainUse                 string    `firestore:"mainUse,omitempty" json:"mainUse,omitempty" bigquery:"-"`
	AlarmTypeInstalled      string    `firestore:"alarmTypeInstalled,omitempty" json:"alarmTypeInstalled,omitempty" bigquery:"-"`
	OwnershipStatus         string    `firestore:"ownershipStatus,omitempty" json:"ownershipStatus,omitempty" bigquery:"-"`
	NumberOfOwners          int64     `firestore:"numberOfOwners,omitempty" json:"numberOfOwners,omitempty" bigquery:"-"`
	VehicleValue            int64     `firestore:"vehicleValue,omitempty" json:"vehicleValue,omitempty" bigquery:"-"`
}
