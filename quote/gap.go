package quote

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-gota/gota/dataframe"
	"gitlab.dev.wopta.it/goworkspace/lib"
	"gitlab.dev.wopta.it/goworkspace/lib/log"
	"gitlab.dev.wopta.it/goworkspace/models"
	"gitlab.dev.wopta.it/goworkspace/network"
	"gitlab.dev.wopta.it/goworkspace/quote/internal"
	"gitlab.dev.wopta.it/goworkspace/sellable"
)

func GapFx(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	var (
		policy  *models.Policy
		warrant *models.Warrant
	)

	log.AddPrefix("GapFx")
	defer log.PopPrefix()

	log.Println("Handler start -----------------------------------------------")

	req := lib.ErrorByte(io.ReadAll(r.Body))
	defer r.Body.Close()

	err := json.Unmarshal(req, &policy)
	lib.CheckError(err)

	policy.Normalize()

	authToken, err := lib.GetAuthTokenFromIdToken(r.Header.Get("Authorization"))
	lib.CheckError(err)

	flow := authToken.GetChannelByRoleV2()

	log.Println("load network node")
	networkNode := network.GetNetworkNodeByUid(policy.ProducerUid)
	if networkNode != nil {
		warrant = networkNode.GetWarrant()
		if warrant != nil {
			flow = warrant.GetFlowName(policy.Name)
		}
	}

	Gap(policy, authToken.GetChannelByRoleV2(), networkNode, warrant, flow)
	policyJson, err := policy.Marshal()

	log.Println("Handler end -------------------------------------------------")

	return string(policyJson), policy, err
}

func Gap(policy *models.Policy, channel string, networkNode *models.NetworkNode, warrant *models.Warrant, flow string) {
	log.AddPrefix("Gap")
	defer log.PopPrefix()

	policy.StartDate = lib.SetDateToStartOfDay(policy.StartDate)

	product, err := sellable.Gap(policy, channel, networkNode, warrant)
	lib.CheckError(err)

	availableRates := internal.GetAvailableRates(product, flow)

	policy.Assets[0].Guarantees = getGuarantees(*product)

	calculateGapOfferPrices(policy, *product)

	log.Println("apply consultacy price")

	internal.AddConsultacyPrice(policy, product)

	internal.RemoveOfferRate(policy, availableRates)
}

func calculateGapOfferPrices(policy *models.Policy, product models.Product) {
	log.AddPrefix("CalculateGapOfferPrices")
	defer log.PopPrefix()
	duration := lib.ElapsedYears(policy.StartDate, policy.EndDate)
	residenceArea := getAreaByProvince(policy.Assets[0].Person.Residence.CityCode)
	policy.Assets[0].Person.Residence.Area = residenceArea
	if residenceArea == "" {
		log.Println("residence area not set")
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
	return lib.CsvToDataframe(lib.GetFilesByEnv(fmt.Sprintf("%s/%s/%s/taxes_%s.csv", models.ProductsFolder,
		productName, productVersion, offerName)))
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
