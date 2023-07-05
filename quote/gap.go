package quote

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-gota/gota/dataframe"
	lib "github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
)

func GapFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	req := lib.ErrorByte(io.ReadAll(r.Body))
	var data models.Policy
	defer r.Body.Close()
	e := json.Unmarshal(req, &data)

	authToken, err := models.GetAuthTokenFromIdToken(r.Header.Get("Authorization"))
	lib.CheckError(err)
	res, e := Gap(authToken.Role, data)
	s, e := json.Marshal(res)
	return string(s), nil, e

}

func Gap(role string, data models.Policy) (models.Policy, error) {
	var err error

	//sellable rules need to be called here

	baseMatrix := lib.GetFilesByEnv("quote/gap_matrix_base.csv")
	completeMatrix := lib.GetFilesByEnv("quote/gap_matrix_complete.csv")
	provincesMatrix := lib.GetFilesByEnv("enrich/provinces.csv")

	baseMatrixDF := lib.CsvToDataframe(baseMatrix)
	completeMatrixDF := lib.CsvToDataframe(completeMatrix)
	provincesMatrixDF := lib.CsvToDataframe(provincesMatrix)

	//get the:
	// - vehicle owner residence
	// - vechicle value
	var residenceCode string = data.Assets[len(data.Assets)-1].Person.Residence.Locality
	var vehicleValue int64 = data.Assets[len(data.Assets)-1].Vehicle.PriceValue

	//get the residence area
	residenceArea := getResidentArea(provincesMatrixDF, residenceCode)

	//get the duration
	duration := getDuration(data.StartDate, data.EndDate)

	fmt.Printf("valore base e' %d\n", duration)

	//get the base and complete multipliers
	baseGapMultiplierFloat, completeGapMultiplierFloat := getGapMultipliers(baseMatrixDF, completeMatrixDF, duration, residenceArea)

	//get the matrix area row
	fmt.Printf("valore base e' %f\n", baseGapMultiplierFloat)
	fmt.Printf("valore completo e' %f\n", completeGapMultiplierFloat)
	fmt.Printf("valore base veicolo*tax e' %f\n", baseGapMultiplierFloat*float64(vehicleValue))
	fmt.Printf("valore completo veicolo*tax e' %f\n", completeGapMultiplierFloat*float64(vehicleValue))

	//set the offer in the policy and round to 2 decimal number

	data = setOffersPrices(data, vehicleValue, baseGapMultiplierFloat, completeGapMultiplierFloat)

	roundGapOffersPrices(data.OffersPrices)

	return data, err
}

func roundGapOffersPrices(offersPrices map[string]map[string]*models.Price) {
	for offerKey, offerValue := range offersPrices {
		for paymentKey := range offerValue {
			offersPrices[offerKey][paymentKey].Net = lib.RoundFloat(offersPrices[offerKey][paymentKey].Net, 2)
		}
	}
}

func setOffersPrices(data models.Policy, vehicleValue int64, baseGapMultiplierFloat float64, completeGapMultiplierFloat float64) models.Policy {
	data.OffersPrices = make(map[string]map[string]*models.Price)

	data.OffersPrices["base"] = map[string]*models.Price{
		"singleInstallment": {
			Net:      baseGapMultiplierFloat * float64(vehicleValue),
			Tax:      0.0,
			Gross:    0.0,
			Delta:    0.0,
			Discount: 0.0,
		},
	}
	data.OffersPrices["complete"] = map[string]*models.Price{
		"singleInstallment": {
			Net:      completeGapMultiplierFloat * float64(vehicleValue),
			Tax:      0.0,
			Gross:    0.0,
			Delta:    0.0,
			Discount: 0.0,
		},
	}

	return data
}

func getResidentArea(provincesMatrixDF dataframe.DataFrame, residenceCode string) string {
	for _, row := range provincesMatrixDF.Records() {
		//fmt.Println(row)
		if row[1] == residenceCode {
			return row[2]
		}
	}
	return ""
}

func getDuration(startDate time.Time, endDate time.Time) int64 {

	return int64(endDate.Year() - startDate.Year())
}

func getGapMultipliers(baseMatrixDF dataframe.DataFrame, completeMatrixDF dataframe.DataFrame, duration int64, residenceArea string) (float64, float64) {

	var baseMatrixAreaRow []string
	var completeMatrixAreaRow []string

	for _, row := range baseMatrixDF.Records() {
		fmt.Println(row)
		if row[0] == residenceArea {
			baseMatrixAreaRow = row
		}
	}

	for _, row := range completeMatrixDF.Records() {
		fmt.Println(row)
		if row[0] == residenceArea {
			completeMatrixAreaRow = row
		}
	}

	baseTaxFloat, _ := strconv.ParseFloat(strings.Replace(baseMatrixAreaRow[duration], "%", "", -1), 64)
	completeTaxFloat, _ := strconv.ParseFloat(strings.Replace(completeMatrixAreaRow[duration], "%", "", -1), 64)

	return baseTaxFloat / 100, completeTaxFloat / 100
}
