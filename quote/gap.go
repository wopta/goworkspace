package quote

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-gota/gota/dataframe"
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
	return string(policyJson), policy, err
}

func Gap(role string, policy *models.Policy) {
	product, err := sellable.Gap(role, policy)
	lib.CheckError(err)

	policy.Assets[0].Guarantees = getGuarantees(product)

	gapPaths := map[string]string{
		"base":     "quote/gap_matrix_base.csv",
		"complete": "quote/gap_matrix_complete.csv",
	}

	provincesMatrix := getDataFrameFromCsv("enrich/provinces.csv")
	for offer := range product.Offers {
		gapMatrix := getDataFrameFromCsv(gapPaths[offer])
		calculateGapOfferPrices(policy, product, offer, gapMatrix, provincesMatrix)
	}
}

func calculateGapOfferPrices(
	policy *models.Policy,
	product models.Product,
	offer string,
	gapMatrix dataframe.DataFrame,
	provincesMatrix dataframe.DataFrame,
) {
	var (
		duration      = lib.ElapsedYears(policy.StartDate, policy.EndDate)
		residenceCode = policy.Assets[0].Person.Residence.CityCode
		residenceArea = getResidentArea(provincesMatrix, residenceCode)
		multiplier    = getGapMultiplier(gapMatrix, duration, residenceArea)

		tax          = getTax(product)
		vehiclePrice = float64(policy.Assets[0].Vehicle.PriceValue)
		netPrice     = multiplier * vehiclePrice
		taxOnPrice   = lib.RoundFloat((netPrice/100)*tax, 2)
		grossPrice   = netPrice + taxOnPrice
	)

	// Check if OffersPrices is not initialized
	if policy.OffersPrices == nil {
		policy.OffersPrices = make(map[string]map[string]*models.Price)
	}

	policy.OffersPrices[offer] = map[string]*models.Price{
		"singleInstallment": {
			Net:      lib.RoundFloat(netPrice, 2),
			Tax:      lib.RoundFloat(taxOnPrice, 2),
			Gross:    lib.RoundFloat(grossPrice, 2),
			Delta:    0.0,
			Discount: 0.0,
		},
	}
}

// Getting the first tax, and assuming every others are the same
func getTax(product models.Product) float64 {
	for _, guarantee := range product.Companies[0].GuaranteesMap {
		return guarantee.Tax
	}
	panic("no tax found")
}

func getDataFrameFromCsv(path string) dataframe.DataFrame {
	csvFile := lib.GetFilesByEnv(path)
	return lib.CsvToDataframe(csvFile)
}

func getGuarantees(product models.Product) []models.Guarante {
	guarantees := make([]models.Guarante, 0)
	for _, guarantee := range product.Companies[0].GuaranteesMap {
		guarantees = append(guarantees, *guarantee)
	}
	return guarantees
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

func getGapMultiplier(residences dataframe.DataFrame, duration int, area string) float64 {
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
