package quote

import (
	"encoding/json"
	"errors"
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

	policy.Assets[0].Guarantees = getGuarantees(*product)

	calculateGapOfferPrices(policy, *product)
}

func calculateGapOfferPrices(policy *models.Policy, product models.Product) {
	duration := lib.ElapsedYears(policy.StartDate, policy.EndDate)
	residenceArea := getAreaByProvince(policy.Assets[0].Person.Residence.CityCode)
	if residenceArea == "" {
		log.Println("[CalculateGapOfferPrices] residence area not set")
		lib.CheckError(errors.New("residence area not set"))
	}
	vehicleValue := float64(policy.Assets[0].Vehicle.PriceValue)
	taxValue := getTax(product)

	policy.OffersPrices = make(map[string]map[string]*models.Price)

	for offerName, _ := range product.Offers {
		matrix := getGapMatrix(offerName)
		multiplier := getGapMultiplier(matrix, duration, residenceArea)

		netPrice := vehicleValue * multiplier
		taxOnPrice := netPrice / 100 * taxValue

		policy.OffersPrices[offerName] = map[string]*models.Price{
			string(models.PaySingleInstallment): {
				Net:      lib.RoundFloat(netPrice, 2),
				Tax:      lib.RoundFloat(taxOnPrice, 2),
				Gross:    lib.RoundFloat(netPrice+taxOnPrice, 2),
				Delta:    0.0,
				Discount: 0.0,
			},
		}
	}
}

func getAreaByProvince(province string) string {
	provincesMatrix := lib.CsvToDataframe(lib.GetFilesByEnv("enrich/provinces.csv"))

	for _, row := range provincesMatrix.Records() {
		if row[1] == province {
			return row[2]
		}
	}
	return ""
}

func getGapMatrix(offerName string) dataframe.DataFrame {
	gapPaths := map[string]string{
		"base":     "quote/gap_matrix_base.csv",
		"complete": "quote/gap_matrix_complete.csv",
	}
	return lib.CsvToDataframe(lib.GetFilesByEnv(gapPaths[offerName]))
}

// Getting the first tax, and assuming every others are the same
func getTax(product models.Product) float64 {
	for _, guarantee := range product.Companies[0].GuaranteesMap {
		return guarantee.Tax
	}
	panic("no tax found")
}

func getGuarantees(product models.Product) []models.Guarante {
	guarantees := make([]models.Guarante, 0)
	for _, guarantee := range product.Companies[0].GuaranteesMap {
		guarantees = append(guarantees, *guarantee)
	}
	return guarantees
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
