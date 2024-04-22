package form

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/bigquery"
	lib "github.com/wopta/goworkspace/lib"
	//"google.golang.org/api/firebaseappcheck/v1"
)

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

var tway = FleetAssistenceInclusiveMovements{
	PolicyNumber:      "191222",
	Lob:               "A",
	PolicyType:        "C",
	CodeSetting:       "1",
	AssetType:         "2",
	Address:           "Piazza Walther Von Der Vogelweide, 22",
	Name:              "T-WAY SPA",
	FleetName:         "T-WAY SPA",
	VatCodeFiscalcode: "3682240043",
	Company:           "AXA",
	Cap:               "39100",
	City:              "BZ",
	Locality:          "Bolzano",
	TypeVehicle:       "3",
	WeightVehicle:     "4",
}

func FleetAssistenceInclusiveMovement(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {

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
		setRequestData(fleetAssistenceInclusiveMovement, &data)
		if fleetAssistenceInclusiveMovement.MovementType == "Inserimento" {

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
func setRequestData(req *FleetAssistenceInclusiveMovements, data *FleetAssistenceInclusiveMovements) *FleetAssistenceInclusiveMovements {
	formatdate := "2006-01-02"
	startdate, e := time.Parse(formatdate, req.CoverageStartDateString)
	log.Println(e)
	enddate, e := time.Parse(formatdate, req.CoverageEndDateString)
	log.Println(e)
	data.PlateVehicle = req.PlateVehicle
	data.ModelVehicle = req.ModelVehicle
	data.CreationDate = lib.GetBigQueryNullDateTime(time.Now())
	data.CoverageStartDate = lib.GetBigQueryNullDateTime(startdate)
	data.CoverageEndDate = lib.GetBigQueryNullDateTime(enddate)
	data.UpdatedDate = lib.GetBigQueryNullDateTime(enddate)
	data.MovementType = req.MovementType
	return data
}
