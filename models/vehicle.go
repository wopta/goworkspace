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
	Plate                   string    `firestore:"plate" json:"plate,omitempty" bigquery:"-"`
	Models                  string    `firestore:"models" json:"models,omitempty" bigquery:"-"`
	Manufacturer            string    `firestore:"manufacturer" json:"manufacturer,omitempty" bigquery:"-"`
	Year                    string    `firestore:"year" json:"year,omitempty" bigquery:"-"`
	RegistrationDate        time.Time `firestore:"registrationDate" json:"registrationDate,omitempty" bigquery:"-"`
	Vin                     string    `firestore:"vin" json:"vin,omitempty" bigquery:"-"`
	BodyType                string    `firestore:"bodyType" json:"bodytype,omitempty" bigquery:"-"`
	Setup                   string    `firestore:"setup" json:"setup,omitempty" bigquery:"-"`
	VehicleOwnerType        string    `firestore:"vehicleOwnerType" json:"vehicleOwnerType,omitempty" bigquery:"-"`
	HasTowHook              bool      `firestore:"hasTowHook" json:"hasTwoHook,omitempty" bigquery:"-"`
	KmPerYear               int64     `firestore:"kmPerYear" json:"kmPerYear,omitempty" bigquery:"-"`
	OvernightVehicleShelter string    `firestore:"overnightVehicleShelter" json:"overnightVehicleShelter,omitempty" bigquery:"-"`
	MainUse                 string    `firestore:"mainUse" json:"mainUse,omitempty" bigquery:"-"`
	CarAlarmTypeInstalled   string    `firestore:"carAlarmTypeInstalled" json:"carAlarmTypeInstalled,omitempty" bigquery:"-"`
	OwnershipStatus         string    `firestore:"ownershipStatus" json:"ownershipStatus,omitempty" bigquery:"-"`
	NumberOfOwners          int64     `firestore:"numberOfOwners" json:"numberOfOwners,omitempty" bigquery:"-"`
	VehicleValue            int64     `firestore:"vehicleValue" json:"vehicleValue,omitempty" bigquery:"-"`
}
