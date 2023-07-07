package quote

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-gota/gota/dataframe"
	"github.com/pkg/errors"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/sellable"
)

func GapFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	log.Println("[GapFx] Handler start")

	req := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	var policy models.Policy
	err := json.Unmarshal(req, &policy)
	lib.CheckError(err)

	authToken, err := models.GetAuthTokenFromIdToken(r.Header.Get("Authorization"))
	lib.CheckError(err)

	Gap(authToken.Role, &policy)
	policyJson, err := json.Marshal(policy)
	return string(policyJson), nil, err
}

func getDataFrameFromCsv(path string) dataframe.DataFrame {
	csvFile := lib.GetFilesByEnv("quote/gap_matrix_base.csv")
	return lib.CsvToDataframe(csvFile)
}

func Gap(role string, p *models.Policy) {
	product, err := sellable.Gap(role, p)
	lib.CheckError(err)

	guarantees := make([]models.Guarante, 0, 10)
	for _, g := range product.Companies[0].GuaranteesMap {
		guarantees = append(guarantees, *g)
	}
	p.Assets[0].Guarantees = guarantees

	gapMatrices := map[string]dataframe.DataFrame{
		"base":     getDataFrameFromCsv("quote/gap_matrix_base.csv"),
		"complete": getDataFrameFromCsv("quote/gap_matrix_complete.csv"),
	}

	provincesMatrix := getDataFrameFromCsv("enrich/provinces.csv")

	residenceCode := p.Assets[0].Person.Residence.Locality
	vehiclePrice := p.Assets[0].Vehicle.PriceValue
	residenceArea := getResidentArea(provincesMatrix, residenceCode)
	duration := lib.ElapsedYears(p.StartDate, p.EndDate)

	// Getting the first tax, and assuming every others are the same
	tax := -1.0
	for _, g := range product.Companies[0].GuaranteesMap {
		tax = g.Tax
		break
	}
	// Just in case
	if tax == -1.0 {
		lib.CheckError(errors.New("Error: no tax found"))
	}

	for offer := range product.Offers {
		gapMatrix := gapMatrices[offer]
		gapMultiplier := mustGetGapMultipliers(gapMatrix, duration, residenceArea)

		initOfferPrices(p, offer, float64(vehiclePrice), gapMultiplier, tax)
	}

	roundOffersPrices(p.OffersPrices)
}

func roundOffersPrices(offersPrices map[string]map[string]*models.Price) {
	for offer, payments := range offersPrices {
		for paymentType, price := range payments {
			offersPrices[offer][paymentType].Net = lib.RoundFloat(price.Net, 2)
			offersPrices[offer][paymentType].Tax = lib.RoundFloat(price.Tax, 2)
			offersPrices[offer][paymentType].Gross = lib.RoundFloat(price.Gross, 2)
		}
	}
}

func initOfferPrices(
	p *models.Policy,
	offer string,
	vehiclePrice float64,
	gapMultiplier float64,
	tax float64,
) {
	netPrice := gapMultiplier * vehiclePrice
	TaxOnPrice := netPrice * (tax / 100)
	TaxOnPrice = lib.RoundFloat(TaxOnPrice, 2)
	grossPrice := netPrice + TaxOnPrice

	p.OffersPrices[offer] = map[string]*models.Price{
		"singleInstallment": {
			Net:      netPrice,
			Tax:      TaxOnPrice,
			Gross:    grossPrice,
			Delta:    0.0,
			Discount: 0.0,
		},
	}
}

// Returns the area (N,C,S) of the residence in input.
// If the return value is "" then no match is found.
func getResidentArea(provincesMatrix dataframe.DataFrame, residenceCode string) string {
	for _, row := range provincesMatrix.Records() {
		if row[1] == residenceCode {
			return row[2]
		}
	}
	return ""
}

func mustGetGapMultipliers(residences dataframe.DataFrame, duration int, area string) float64 {
	var matrixAreaRow []string

	// We assume that the area is unique, hence the first match is the only one
	for _, row := range residences.Records() {
		if row[0] == area {
			matrixAreaRow = row
			break
		}
	}

	taxString := strings.Replace(matrixAreaRow[duration], "%", "", -1)
	tax, err := strconv.ParseFloat(taxString, 64)
	lib.CheckError(err)

	return tax / 100
}
