package quote

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-gota/gota/dataframe"
	"github.com/wopta/goworkspace/lib"
	"github.com/wopta/goworkspace/models"
	"github.com/wopta/goworkspace/network"
	"github.com/wopta/goworkspace/sellable"
)

func GapFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		policy  *models.Policy
		warrant *models.Warrant
	)

	log.Println("[GapFx] handler start --------------------------------------")

	req := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	err := json.Unmarshal(req, &policy)
	lib.CheckError(err)

	authToken, err := models.GetAuthTokenFromIdToken(r.Header.Get("Authorization"))
	lib.CheckError(err)

	log.Println("[GapFx] load network node")
	networkNode := network.GetNetworkNodeByUid(authToken.UserID)
	if networkNode != nil {
		warrant = networkNode.GetWarrant()
	}

	Gap(policy, authToken.GetChannelByRoleV2(), networkNode, warrant)
	policyJson, err := policy.Marshal()

	log.Printf("[GapFx] response: %s", string(policyJson))

	log.Println("[GapFx] handler end --------------------------------------")

	return string(policyJson), policy, err
}

func Gap(policy *models.Policy, channel string, networkNode *models.NetworkNode, warrant *models.Warrant) {
	product, err := sellable.Gap(policy, channel, networkNode, warrant)
	lib.CheckError(err)

	policy.Assets[0].Guarantees = getGuarantees(*product)

	calculateGapOfferPrices(policy, *product)
}

func calculateGapOfferPrices(policy *models.Policy, product models.Product) {
	duration := lib.ElapsedYears(policy.StartDate, policy.EndDate)
	residenceArea := getAreaByProvince(policy.Assets[0].Person.Residence.CityCode)
	policy.Assets[0].Person.Residence.Area = residenceArea
	if residenceArea == "" {
		log.Println("[CalculateGapOfferPrices] residence area not set")
		lib.CheckError(errors.New("residence area not set"))
	}
	vehicleValue := policy.Assets[0].Vehicle.PriceValue
	taxValue := getTax(product)

	policy.OffersPrices = make(map[string]map[string]*models.Price)

	for offerName := range product.Offers {
		matrix := getGapMatrix(policy.Name, policy.ProductVersion, offerName)
		multiplier := getGapMultiplier(matrix, duration, residenceArea)

		netPrice := vehicleValue * multiplier
		taxOnPrice := netPrice / 100 * taxValue

		policy.OffersPrices[offerName] = map[string]*models.Price{
			string(models.PaySplitSingleInstallment): {
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

func getGapMatrix(productName, productVersion, offerName string) dataframe.DataFrame {
	return lib.CsvToDataframe(lib.GetFilesByEnv(fmt.Sprintf("products-v2/%s/%s/taxes_%s.csv", productName, productVersion, offerName)))
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
