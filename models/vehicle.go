package models

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
	Plate                   string `firestore:"plate" json:"plate,omitempty" bigquery:"plate"`
	Models                  string `firestore:"models" json:"models,omitempty" bigquery:"models"`
	Manufacturer            string `firestore:"manufacturer" json:"manufacturer,omitempty" bigquery:"manufacturer"`
	Year                    string `firestore:"year" json:"year,omitempty" bigquery:"year"`
	RegistrationDate        string `firestore:"registrationDate" json:"registrationDate,omitempty" bigquery:"registrationDate"`
	Vin                     string `firestore:"vin" json:"vin,omitempty" bigquery:"vin"`
	BodyType                string `firestore:"bodytype" json:"bodytype,omitempty" bigquery:"bodytype"`
	Setup                   string `firestore:"setup" json:"setup,omitempty" bigquery:"setup"`
	VehicleOwnerType        string `firestore:"vehicleOwnerType" json:"vehicleOwnerType,omitempty" bigquery:"vehicleOwnerType"`
	TowHosokPresence        bool   `firestore:"towHosokPresence" json:"towHosokPresence,omitempty" bigquery:"towHosokPresence"`
	KmPerYear               uint   `firestore:"kmPerYear" json:"kmPerYear,omitempty" bigquery:"kmPerYear"`
	OvernightVehicleShelter string `firestore:"overnightVehicleShelter" json:"overnightVehicleShelter,omitempty" bigquery:"overnightVehicleShelter"`
	MainUse                 string `firestore:"mainUse" json:"mainUse,omitempty" bigquery:"mainUse"`
	CarAlarmTypeInstalled   string `firestore:"carAlarmTypeInstalled" json:"carAlarmTypeInstalled,omitempty" bigquery:"carAlarmTypeInstalled"`
	AlreadyOwner            string `firestore:"alreadyOwner" json:"alreadyOwner,omitempty" bigquery:"alreadyOwner"`
	OwnerNumber             uint   `firestore:"ownerNumber" json:"ownerNumber,omitempty" bigquery:"ownerNumber"`
	VehicleValue            uint   `firestore:"vehicleValue" json:"vehicleValue,omitempty" bigquery:"vehicleValue"`
}
