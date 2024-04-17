package form

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"cloud.google.com/go/bigquery"
	lib "github.com/wopta/goworkspace/lib"
	//"google.golang.org/api/firebaseappcheck/v1"
)

// {"responses":{"TIPO MOVIMENTO":"Inserimento","Targa Inserimento":"test","MODELLO VEICOLO":"test mod","DATA IMMATRICOLAZIONE":"1212-12-02","DATA INIZIO VALIDITA' COPERTURA":"1212-12-12"}}
// {"responses":{"TIPO MOVIMENTO":"Annullo","Targa Annullo":"targa","DATA FINE VALIDITA' COPERTURA":"0009-09-09"},"mail":"test@gmail.com"}

type FleetAssistenceInclusiveMovements struct {
	PolicyNumber            string                `json:"-" firestore:"-" bigquery:"policyNumber"`
	Lob                     string                `json:"-" firestore:"-" bigquery:"lob"`
	PolicyType              string                `json:"-" firestore:"-" bigquery:"policyType"`
	CodeSetting             string                `json:"-" firestore:"-" bigquery:"codeSetting"`
	VatCodeFiscalcode       string                `json:"-" firestore:"-" bigquery:"vatCodeFiscalcode"`
	IsActive                bool                  `json:"-" firestore:"-" bigquery:"isActive "`
	Id                      string                `json:"-" firestore:"-" bigquery:"id"`
	AssetType               string                `json:"-" firestore:"-" bigquery:"assetType"`
	Name                    string                `json:"-" firestore:"-" bigquery:"name"`
	FleetName               string                `json:"-" firestore:"-" bigquery:"fleetName"`
	Company                 string                `json:"-" firestore:"-" bigquery:"company"`
	Address                 string                `json:"-" firestore:"-" bigquery:"address"`
	Cap                     string                `json:"-" firestore:"-" bigquery:"cap"`
	City                    string                `json:"-" firestore:"-" bigquery:"city"`
	Locality                string                `json:"-" firestore:"-" bigquery:"locality"`
	PlateVehicle            string                `json:"TARGA" firestore:"-" bigquery:"plateVehicle"`
	ModelVehicle            string                `json:"MODELLO VEICOLO" firestore:"-" bigquery:"modelVehicle"`
	TypeVehicle             string                `json:"-" firestore:"-" bigquery:"typeVehicle"`
	BrandVehicle            string                `json:"-" firestore:"-" bigquery:"brandVehicle"`
	FrameVehicle            string                `json:"-" firestore:"-" bigquery:"frameVehicle"`
	WeightVehicle           string                `json:"-" firestore:"-" bigquery:"weightVehicle"`
	RegistrationVehicle     string                `json:"DATA IMMATRICOLAZIONE" firestore:"-" bigquery:"registrationVehicle"`
	MovementType            string                `json:"TIPO MOVIMENTO" firestore:"-" bigquery:"movementType"`
	Mail                    string                `json:"mail" firestore:"-" bigquery:"mail"`
	CoverageStartDateString string                `json:"DATA INIZIO VALIDITA COPERTURA" firestore:"-" bigquery:"-"`
	CoverageEndDateString   string                `json:"DATA FINE VALIDITA COPERTURA" firestore:"-" bigquery:"-"`
	CoverageStartDate       bigquery.NullDateTime `json:"-" firestore:"-" bigquery:"coverageStartDate"`
	CoverageEndDate         bigquery.NullDateTime `json:"-" firestore:"-" bigquery:"coverageEndDate"`
	CreationDate            bigquery.NullDateTime `json:"-" firestore:"-" bigquery:"creationDate"`
	UpdatedDate             bigquery.NullDateTime `json:"-" firestore:"-" bigquery:"UpdatedDate "`
}

func FleetAssistenceInclusiveMovement(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {

	var tway = FleetAssistenceInclusiveMovements{
		PolicyNumber:      "191222",
		Lob:               "A",
		PolicyType:        "C",
		CodeSetting:       "1",
		AssetType:         "2",
		Address:           "Piazza Walther Von Der Vogelweide, 22",
		Name:              "T-WAY SPA",
		FleetName: "T-WAY SPA",
		VatCodeFiscalcode: "3682240043",
		Company:           "AXA",
		Cap: "39100",
		City: "BZ",
		Locality:"Bolzano" ,
		TypeVehicle: "3",
		WeightVehicle: "4",
	}
	var users = map[string]FleetAssistenceInclusiveMovements{
		"elisabetta.lainatiassicura@gmail.com": tway,
	}
	var (
		fleetAssistenceInclusiveMovement *FleetAssistenceInclusiveMovements
	)
	body, e := io.ReadAll(r.Body)
	log.Println(e)
	log.Println(string(body))
	json.Unmarshal(body, &fleetAssistenceInclusiveMovement)
	data, ok := users[fleetAssistenceInclusiveMovement.Mail]

	if ok {
		if fleetAssistenceInclusiveMovement.MovementType == "Inserimento" {
			setRequestData(fleetAssistenceInclusiveMovement, &data)
			lib.InsertRowsBigQuery("wopta", "fleetAssistenceInclusiveMovements", data)
		} else {
			checkPlate, e := lib.QueryRowsBigQuery[FleetAssistenceInclusiveMovements]("")
			log.Println(e)
			if len(checkPlate) > 0 {
				lib.InsertRowsBigQuery("wopta", "fleetAssistenceInclusiveMovements", data)
			} else {

			}
		}

	} else {

	}

	return "", nil, nil
}
func setRequestData(req *FleetAssistenceInclusiveMovements, data *FleetAssistenceInclusiveMovements) {
	data.PlateVehicle = req.PlateVehicle
	data.ModelVehicle = req.ModelVehicle
	

	data.PlateVehicle = req.PlateVehicle
	data.PlateVehicle = req.PlateVehicle
}
